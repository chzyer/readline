package readline

import (
	"container/list"
	"io"
	"os"
)

type Operation struct {
	r       *os.File
	t       *Terminal
	buf     *RuneBuffer
	outchan chan []rune

	history    *list.List
	historyVer int64
	current    *list.Element
}

const (
	CharLineStart = 0x1
	CharLineEnd   = 0x5
	CharKill      = 11
	CharNext      = 0xe
	CharPrev      = 0x10
	CharBackward  = 0x2
	CharForward   = 0x6
	CharBackspace = 0x7f
	CharEnter     = 0xd
	CharEnter2    = 0xa
)

type wrapWriter struct {
	r      *Operation
	target io.Writer
}

func (w *wrapWriter) Write(b []byte) (int, error) {
	buf := w.r.buf
	buf.Clean()
	n, err := w.target.Write(b)
	w.r.buf.RefreshSet(0, 0)
	return n, err
}

func NewOperation(r *os.File, t *Terminal, prompt string) *Operation {
	op := &Operation{
		r:       r,
		t:       t,
		buf:     NewRuneBuffer(t, prompt),
		outchan: make(chan []rune),
		history: list.New(),
	}
	go op.ioloop()
	return op
}

func (l *Operation) ioloop() {
	for {
		r := l.t.ReadRune()
		switch r {
		case CharKill:
			l.buf.Kill()
		case MetaNext:
			l.buf.MoveToNextWord()
		case MetaPrev:
			l.buf.MoveToPrevWord()
		case MetaDelete:
			l.buf.DeleteWord()
		case CharLineStart:
			l.buf.MoveToLineStart()
		case CharLineEnd:
			l.buf.MoveToLineEnd()
		case KeyDelete:
			l.buf.Delete()
		case CharBackspace:
			l.buf.Backspace()
		case MetaBackspace:
			l.buf.BackEscapeWord()
		case CharEnter, CharEnter2:
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
		case KeyInterrupt:
			l.buf.WriteString("^C\n")
			l.outchan <- nil
		default:
			l.buf.WriteRune(r)
		}
		l.UpdateHistory(l.buf.Runes(), false)
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
	l.buf.Refresh(0, 0) // print prompt
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
