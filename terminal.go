package readline

import (
	"bufio"
	"fmt"
	"os"
	"sync/atomic"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

const (
	MetaPrev = -iota - 1
	MetaNext
	MetaDelete
	MetaBackspace
)

const (
	KeyPrevChar  = 0x2
	KeyInterrupt = 0x3
	KeyNextChar  = 0x6
	KeyDelete    = 0x4
	KeyEsc       = 0x1b
)

type Terminal struct {
	state   *terminal.State
	outchan chan rune
	closed  int64
}

func NewTerminal() (*Terminal, error) {
	state, err := MakeRaw(syscall.Stdin)
	if err != nil {
		return nil, err
	}
	t := &Terminal{
		state:   state,
		outchan: make(chan rune),
	}

	go t.ioloop()
	return t, nil
}

func (t *Terminal) Write(b []byte) (int, error) {
	return os.Stdout.Write(b)
}

func (t *Terminal) Print(s string) {
	fmt.Fprintf(os.Stdout, "%s", s)
}

func (t *Terminal) PrintRune(r rune) {
	fmt.Fprintf(os.Stdout, "%c", r)
}

func (t *Terminal) Readline(prompt string) *Operation {
	return NewOperation(os.Stdin, t, prompt)
}

func (t *Terminal) ReadRune() rune {
	return <-t.outchan
}

func (t *Terminal) ioloop() {
	buf := bufio.NewReader(os.Stdin)
	prefix := false
	for {
		r, _, err := buf.ReadRune()
		if err != nil {
			break
		}

		if prefix {
			prefix = false
			r = prefixKey(r)
		}

		if IsPrintable(r) || r < 0 {
			t.outchan <- r
			continue
		}
		switch r {
		case KeyInterrupt:
			t.outchan <- r
			goto exit
		case KeyEsc:
			prefix = true
		case CharEnter, CharEnter2, KeyPrevChar, KeyNextChar, KeyDelete:
			fallthrough
		case CharLineEnd, CharLineStart, CharNext, CharPrev:
			t.outchan <- r
		default:
			println("np:", r)
		}
	}
exit:
}

func (t *Terminal) Close() error {
	if atomic.SwapInt64(&t.closed, 1) != 0 {
		return nil
	}
	return Restore(syscall.Stdin, t.state)
}
