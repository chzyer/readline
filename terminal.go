package readline

import (
	"bufio"
	"fmt"
	"os"
	"sync/atomic"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

type Terminal struct {
	cfg     *Config
	state   *terminal.State
	outchan chan rune
	closed  int64
}

func NewTerminal(cfg *Config) (*Terminal, error) {
	state, err := MakeRaw(syscall.Stdin)
	if err != nil {
		return nil, err
	}
	t := &Terminal{
		cfg:     cfg,
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

func (t *Terminal) Readline() *Operation {
	return NewOperation(t, t.cfg)
}

func (t *Terminal) ReadRune() rune {
	return <-t.outchan
}

func (t *Terminal) ioloop() {
	buf := bufio.NewReader(os.Stdin)
	isEscape := false
	isEscapeEx := false
	for {
		r, _, err := buf.ReadRune()
		if err != nil {
			break
		}

		if isEscape {
			isEscape = false
			if r == CharEscapeEx {
				isEscapeEx = true
				continue
			}
			r = escapeKey(r)
		} else if isEscapeEx {
			isEscapeEx = false
			r = escapeExKey(r)
		}

		switch r {
		case CharInterrupt:
			t.outchan <- r
			goto exit
		case CharEsc:
			isEscape = true
		default:
			t.outchan <- r
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
