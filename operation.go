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

func (l *Operation) ioloop() {
	for {
		keepInSearchMode := false
		r := l.t.ReadRune()
		switch r {
		case CharCannel:
			if l.IsSearchMode() {
				l.ExitSearchMode(true)
				l.buf.Refresh()
			}
		case CharBckSearch:
			l.SearchMode(S_DIR_BCK)
			keepInSearchMode = true
		case CharFwdSearch:
			l.SearchMode(S_DIR_FWD)
			keepInSearchMode = true
		case CharKill:
			l.buf.Kill()
		case MetaNext:
			l.buf.MoveToNextWord()
		case CharTranspose:
			l.buf.Transpose()
		case MetaPrev:
			l.buf.MoveToPrevWord()
		case MetaDelete:
			l.buf.DeleteWord()
		case CharLineStart:
			l.buf.MoveToLineStart()
		case CharLineEnd:
			l.buf.MoveToLineEnd()
		case CharDelete:
			l.buf.Delete()
		case CharBackspace, CharCtrlH:
			if l.IsSearchMode() {
				l.SearchBackspace()
				keepInSearchMode = true
			} else {
				l.buf.Backspace()
			}
		case MetaBackspace, CharCtrlW:
			l.buf.BackEscapeWord()
		case CharEnter, CharCtrlJ:
			if l.IsSearchMode() {
				l.ExitSearchMode(false)
			}
			l.buf.MoveToLineEnd()
			l.buf.WriteRune('\n')
			data := l.buf.Reset()
			data = data[:len(data)-1] // trim \n
			l.outchan <- data
			l.NewHistory(data)
		case CharBackward:
			l.buf.MoveBackward()
		case CharForward:
			l.buf.MoveForward()
		case CharPrev:
			buf := l.PrevHistory()
			if buf != nil {
				l.buf.Set(buf)
			}
		case CharNext:
			buf, ok := l.NextHistory()
			if ok {
				l.buf.Set(buf)
			}
		case CharInterrupt:
			if l.IsSearchMode() {
				l.ExitSearchMode(false)
			}
			l.buf.MoveToLineEnd()
			l.buf.Refresh()
			l.buf.WriteString("^C\n")
			l.outchan <- nil
		default:
			if l.IsSearchMode() {
				l.SearchChar(r)
				keepInSearchMode = true
			} else {
				l.buf.WriteRune(r)
			}
		}
		if !keepInSearchMode && l.IsSearchMode() {
			l.ExitSearchMode(false)
			l.buf.Refresh()
		}
		if !l.IsSearchMode() {
			l.UpdateHistory(l.buf.Runes(), false)
		}
	}
}

func (l *Operation) Stderr() io.Writer {
	return &wrapWriter{target: os.Stderr, r: l}
}

func (l *Operation) String() (string, error) {
	r, err := l.Runes()
	if err != nil {
		return "", err
	}
	return string(r), nil
}

func (l *Operation) Runes() ([]rune, error) {
	l.buf.Refresh() // print prompt
	r := <-l.outchan
	if r == nil {
		return nil, io.EOF
	}
	return r, nil
}

func (l *Operation) Slice() ([]byte, error) {
	r, err := l.Runes()
	if err != nil {
		return nil, err
	}
	return []byte(string(r)), nil
}

func (l *Operation) Close() {
	l.opHistory.Close()
}
