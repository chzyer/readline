package readline

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

type Terminal struct {
	cfg       *Config
	state     *terminal.State
	outchan   chan rune
	closed    int64
	stopChan  chan struct{}
	kickChan  chan struct{}
	wg        sync.WaitGroup
	isReading bool
}

func NewTerminal(cfg *Config) (*Terminal, error) {
	state, err := MakeRaw(syscall.Stdin)
	if err != nil {
		return nil, err
	}
	t := &Terminal{
		cfg:      cfg,
		state:    state,
		kickChan: make(chan struct{}, 1),
		outchan:  make(chan rune),
		stopChan: make(chan struct{}, 1),
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

func (t *Terminal) IsReading() bool {
	return t.isReading
}

func (t *Terminal) KickRead() {
	select {
	case t.kickChan <- struct{}{}:
	default:
	}
}

func (t *Terminal) ioloop() {
	t.wg.Add(1)
	defer t.wg.Done()
	var (
		isEscape       bool
		isEscapeEx     bool
		expectNextChar bool
	)

	buf := bufio.NewReader(os.Stdin)
	for {
		if !expectNextChar {
			t.isReading = false
			select {
			case <-t.kickChan:
				t.isReading = true
			case <-t.stopChan:
				return
			}
		}
		expectNextChar = false
		r, _, err := buf.ReadRune()
		if err != nil {
			break
		}

		if isEscape {
			isEscape = false
			if r == CharEscapeEx {
				expectNextChar = true
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
			expectNextChar = true
		case CharEnter, CharCtrlJ:
			t.outchan <- r
		default:
			expectNextChar = true
			t.outchan <- r
		}
	}
exit:
}

func (t *Terminal) Close() error {
	if atomic.SwapInt64(&t.closed, 1) != 0 {
		return nil
	}
	t.stopChan <- struct{}{}
	t.wg.Wait()
	return Restore(syscall.Stdin, t.state)
}
