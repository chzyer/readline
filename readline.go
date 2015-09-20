package readline

import (
	"io"
	"os"
)

type Readline struct {
	r       *os.File
	t       *Terminal
	buf     *RuneBuffer
	outchan chan []rune
}

const (
	CharLineStart = 0x1
	CharLineEnd   = 0x5
	CharPrev      = 0x2
	CharNext      = 0x6
	CharEscape    = 0x7f
	CharEnter     = 0xd
)

type wrapWriter struct {
	r      *Readline
	target io.Writer
}

func (w *wrapWriter) Write(b []byte) (int, error) {
	buf := w.r.buf
	buf.Clean()
	n, err := w.target.Write(b)
	w.r.buf.RefreshSet(0, 0)
	return n, err
}

func newReadline(r *os.File, t *Terminal, prompt string) *Readline {
	rl := &Readline{
		r:       r,
		t:       t,
		buf:     NewRuneBuffer(t, prompt),
		outchan: make(chan []rune),
	}
	go rl.ioloop()
	return rl
}

func (l *Readline) ioloop() {
	for {
		r := l.t.ReadRune()
		switch r {
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
		case CharEscape:
			l.buf.BackEscape()
		case CharEnter:
			l.buf.WriteRune('\n')
			data := l.buf.Reset()
			l.outchan <- data[:len(data)-1]
		case CharPrev:
			l.buf.MovePrev()
		case CharNext:
			l.buf.MoveNext()
		case KeyInterrupt:
			l.buf.WriteString("^C\n")
			l.outchan <- nil
			break
		default:
			l.buf.WriteRune(r)
		}
	}
}

func (l *Readline) Stderr() io.Writer {
	return &wrapWriter{target: os.Stderr, r: l}
}

func (l *Readline) Readline() (string, error) {
	r, err := l.ReadlineSlice()
	if err != nil {
		return "", err
	}
	return string(r), nil
}

func (l *Readline) ReadlineSlice() ([]byte, error) {
	l.buf.Refresh(0, 0)
	r := <-l.outchan
	if r == nil {
		return nil, io.EOF
	}
	return []byte(string(r)), nil
}
