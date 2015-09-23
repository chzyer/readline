package readline

import (
	"io"
	"os"
)

type Operation struct {
	cfg     *Config
	t       *Terminal
	buf     *RuneBuffer
	outchan chan []rune
	*opHistory
	*opSearch
}

type wrapWriter struct {
	r      *Operation
	target io.Writer
}

func (w *wrapWriter) Write(b []byte) (int, error) {
	buf := w.r.buf
	buf.Clean()
	n, err := w.target.Write(b)
	w.r.buf.Refresh()
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
	op.opSearch = newOpSearch(op.buf.w, op.buf, op.opHistory)
	go op.ioloop()
	return op
}

func (o *Operation) ioloop() {
	for {
		keepInSearchMode := false
		r := o.t.ReadRune()
		switch r {
		case CharCannel:
			if o.IsSearchMode() {
				o.ExitSearchMode(true)
				o.buf.Refresh()
			}
		case CharBckSearch:
			o.SearchMode(S_DIR_BCK)
			keepInSearchMode = true
		case CharFwdSearch:
			o.SearchMode(S_DIR_FWD)
			keepInSearchMode = true
		case CharKill:
			o.buf.Kill()
		case MetaNext:
			o.buf.MoveToNextWord()
		case CharTranspose:
			o.buf.Transpose()
		case MetaPrev:
			o.buf.MoveToPrevWord()
		case MetaDelete:
			o.buf.DeleteWord()
		case CharLineStart:
			o.buf.MoveToLineStart()
		case CharLineEnd:
			o.buf.MoveToLineEnd()
		case CharDelete:
			o.buf.Delete()
		case CharBackspace, CharCtrlH:
			if o.IsSearchMode() {
				o.SearchBackspace()
				keepInSearchMode = true
			} else {
				o.buf.Backspace()
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
			}
		case CharNext:
			buf, ok := o.NextHistory()
			if ok {
				o.buf.Set(buf)
			}
		case CharInterrupt:
			if o.IsSearchMode() {
				o.ExitSearchMode(false)
			}
			o.buf.MoveToLineEnd()
			o.buf.Refresh()
			o.buf.WriteString("^C\n")
			o.outchan <- nil
		default:
			if o.IsSearchMode() {
				o.SearchChar(r)
				keepInSearchMode = true
			} else {
				o.buf.WriteRune(r)
			}
		}
		if !keepInSearchMode && o.IsSearchMode() {
			o.ExitSearchMode(false)
			o.buf.Refresh()
		}
		if !o.IsSearchMode() {
			o.UpdateHistory(o.buf.Runes(), false)
		}
	}
}

func (o *Operation) Stderr() io.Writer {
	return &wrapWriter{target: os.Stderr, r: o}
}

func (o *Operation) String() (string, error) {
	r, err := o.Runes()
	if err != nil {
		return "", err
	}
	return string(r), nil
}

func (o *Operation) Runes() ([]rune, error) {
	o.buf.Refresh() // print prompt
	r := <-o.outchan
	if r == nil {
		return nil, io.EOF
	}
	return r, nil
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
