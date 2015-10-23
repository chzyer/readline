package readline

import (
	"bufio"
	"fmt"
	"sync"
	"sync/atomic"

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
	isReading int64
}

func NewTerminal(cfg *Config) (*Terminal, error) {
	if err := cfg.Init(); err != nil {
		return nil, err
	}
	state, err := MakeRaw(StdinFd)
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
	return t.cfg.Stdout.Write(b)
}

func (t *Terminal) Print(s string) {
	fmt.Fprintf(t.cfg.Stdout, "%s", s)
}

func (t *Terminal) PrintRune(r rune) {
	fmt.Fprintf(t.cfg.Stdout, "%c", r)
}

func (t *Terminal) Readline() *Operation {
	return NewOperation(t, t.cfg)
}

func (t *Terminal) ReadRune() rune {
	return <-t.outchan
}

func (t *Terminal) IsReading() bool {
	return atomic.LoadInt64(&t.isReading) == 1
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

	buf := bufio.NewReader(Stdin)
	for {
		if !expectNextChar {
			atomic.StoreInt64(&t.isReading, 0)
			select {
			case <-t.kickChan:
				atomic.StoreInt64(&t.isReading, 1)
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
			// if hw delete button is pressed it is specified as set ot 4 runes [27,91,51,126]. we are now at 51
			if r == CharDelete {
				if d, _, err := buf.ReadRune(); err != nil || d != 126 {
					buf.UnreadRune()
				}
			}
		}

		expectNextChar = true
		switch r {
		case CharEsc:
			if t.cfg.VimMode {
				t.outchan <- r
				break
			}
			isEscape = true
		case CharInterrupt, CharEnter, CharCtrlJ:
			expectNextChar = false
			fallthrough
		default:
			t.outchan <- r
		}
	}
}

func (t *Terminal) Bell() {
	fmt.Fprintf(t, "%c", CharBell)
}

func (t *Terminal) Close() error {
	if atomic.SwapInt64(&t.closed, 1) != 0 {
		return nil
	}
	t.stopChan <- struct{}{}
	t.wg.Wait()
	return Restore(StdinFd, t.state)
}
