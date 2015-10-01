package readline

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

type Operation struct {
	cfg     *Config
	t       *Terminal
	buf     *RuneBuffer
	outchan chan []rune
	w       io.Writer

	*opHistory
	*opSearch
	*opCompleter
	*opVim
}

type wrapWriter struct {
	r      *Operation
	t      *Terminal
	target io.Writer
}

func (w *wrapWriter) Write(b []byte) (int, error) {
	if !w.t.IsReading() {
		return w.target.Write(b)
	}

	var (
		n   int
		err error
	)
	w.r.buf.Refresh(func() {
		n, err = w.target.Write(b)
	})

	if w.r.IsSearchMode() {
		w.r.SearchRefresh(-1)
	}
	if w.r.IsInCompleteMode() {
		w.r.CompleteRefresh()
	}
	return n, err
}

func NewOperation(t *Terminal, cfg *Config) *Operation {
	op := &Operation{
		cfg:       cfg,
		t:         t,
		buf:       NewRuneBuffer(t, cfg.Prompt),
		outchan:   make(chan []rune),
		opHistory: newOpHistory(cfg.HistoryFile),
	}
	op.opVim = newVimMode(op)
	op.w = op.buf.w
	op.opSearch = newOpSearch(op.buf.w, op.buf, op.opHistory)
	op.opCompleter = newOpCompleter(op.buf.w, op)
	go op.ioloop()
	return op
}

func (o *Operation) SetPrompt(s string) {
	o.buf.SetPrompt(s)
}

func (o *Operation) ioloop() {
	for {
		keepInSearchMode := false
		keepInCompleteMode := false
		r := o.t.ReadRune()

		if o.IsInCompleteSelectMode() {
			keepInCompleteMode = o.HandleCompleteSelect(r)
			if keepInCompleteMode {
				continue
			}

			o.buf.Refresh(nil)
			switch r {
			case CharEnter, CharCtrlJ:
				o.UpdateHistory(o.buf.Runes(), false)
				fallthrough
			case CharInterrupt:
				o.t.KickRead()
				fallthrough
			case CharBell:
				continue
			}
		}

		if o.IsEnableVimMode() {
			var ok bool
			r, ok = o.HandleVim(r, o.t.ReadRune)
			if ok {
				continue
			}
		}

		switch r {
		case CharBell:
			if o.IsSearchMode() {
				o.ExitSearchMode(true)
				o.buf.Refresh(nil)
			}
			if o.IsInCompleteMode() {
				o.ExitCompleteMode(true)
				o.buf.Refresh(nil)
			}
		case CharTab:
			if o.opCompleter == nil {
				break
			}
			o.OnComplete()
			keepInCompleteMode = true
		case CharBckSearch:
			o.SearchMode(S_DIR_BCK)
			keepInSearchMode = true
		case CharCtrlU:
			o.buf.KillFront()
		case CharFwdSearch:
			o.SearchMode(S_DIR_FWD)
			keepInSearchMode = true
		case CharKill:
			o.buf.Kill()
			keepInCompleteMode = true
		case MetaForward:
			o.buf.MoveToNextWord()
		case CharTranspose:
			o.buf.Transpose()
		case MetaBackward:
			o.buf.MoveToPrevWord()
		case MetaDelete:
			o.buf.DeleteWord()
		case CharLineStart:
			o.buf.MoveToLineStart()
		case CharLineEnd:
			o.buf.MoveToLineEnd()
		case CharDelete:
			if !o.buf.Delete() {
				o.t.Bell()
			}
		case CharBackspace, CharCtrlH:
			if o.IsSearchMode() {
				o.SearchBackspace()
				keepInSearchMode = true
				break
			}

			if o.buf.Len() == 0 {
				o.t.Bell()
				break
			}
			o.buf.Backspace()
			if o.IsInCompleteMode() {
				o.OnComplete()
			}
		case MetaBackspace, CharCtrlW:
			o.buf.BackEscapeWord()
		case CharEnter, CharCtrlJ:
			if o.IsSearchMode() {
				o.ExitSearchMode(false)
			}
			o.buf.MoveToLineEnd()
			o.buf.WriteRune('\n')
			data := o.buf.Reset()
			data = data[:len(data)-1] // trim \n
			o.outchan <- data
			o.NewHistory(data)
		case CharBackward:
			o.buf.MoveBackward()
		case CharForward:
			o.buf.MoveForward()
		case CharPrev:
			buf := o.PrevHistory()
			if buf != nil {
				o.buf.Set(buf)
			} else {
				o.t.Bell()
			}
		case CharNext:
			buf, ok := o.NextHistory()
			if ok {
				o.buf.Set(buf)
			} else {
				o.t.Bell()
			}
		case CharInterrupt:
			if o.IsSearchMode() {
				o.t.KickRead()
				o.ExitSearchMode(true)
				break
			}
			if o.IsInCompleteMode() {
				o.t.KickRead()
				o.ExitCompleteMode(true)
				o.buf.Refresh(nil)
				break
			}
			o.buf.MoveToLineEnd()
			o.buf.Refresh(nil)
			o.buf.WriteString("^C\n")
			o.outchan <- nil
		default:
			if o.IsSearchMode() {
				o.SearchChar(r)
				keepInSearchMode = true
				break
			}
			o.buf.WriteRune(r)
			if o.IsInCompleteMode() {
				o.OnComplete()
				keepInCompleteMode = true
			}
		}

		if !keepInSearchMode && o.IsSearchMode() {
			o.ExitSearchMode(false)
			o.buf.Refresh(nil)
		} else if o.IsInCompleteMode() {
			if !keepInCompleteMode {
				o.ExitCompleteMode(false)
				o.buf.Refresh(nil)
			} else {
				o.buf.Refresh(nil)
				o.CompleteRefresh()
			}
		}
		if !o.IsSearchMode() {
			o.UpdateHistory(o.buf.Runes(), false)
		}
	}
}

func (o *Operation) Stderr() io.Writer {
	return &wrapWriter{target: o.cfg.Stderr, r: o, t: o.t}
}

func (o *Operation) Stdout() io.Writer {
	return &wrapWriter{target: o.cfg.Stdout, r: o, t: o.t}
}

func (o *Operation) String() (string, error) {
	r, err := o.Runes()
	if err != nil {
		return "", err
	}
	return string(r), nil
}

func (o *Operation) Runes() ([]rune, error) {
	o.buf.Refresh(nil) // print prompt
	o.t.KickRead()
	r := <-o.outchan
	if r == nil {
		return nil, io.EOF
	}
	return r, nil
}

func (o *Operation) Password(prompt string) ([]byte, error) {
	w := o.Stdout()
	if prompt != "" {
		fmt.Fprintf(w, prompt)
	}
	b, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	fmt.Fprint(w, "\r\n")
	return b, err
}

func (o *Operation) SetTitle(t string) {
	o.w.Write([]byte("\033[2;" + t + "\007"))
}

func (o *Operation) Slice() ([]byte, error) {
	r, err := o.Runes()
	if err != nil {
		return nil, err
	}
	return []byte(string(r)), nil
}

func (o *Operation) Close() {
	o.opHistory.Close()
}
