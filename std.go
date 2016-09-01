package readline

import (
	"io"
	"os"
	"sync"
)

var (
	Stdin  io.ReadCloser  = NewCancelableStdin()
	Stdout io.WriteCloser = os.Stdout
	Stderr io.WriteCloser = os.Stderr
)

var (
	std     *Instance
	stdOnce sync.Once
)

// global instance will not submit history automatic
func getInstance() *Instance {
	stdOnce.Do(func() {
		std, _ = NewEx(&Config{
			DisableAutoSaveHistory: true,
		})
	})
	return std
}

// let readline load history from filepath
// and try to persist history into disk
// set fp to "" to prevent readline persisting history to disk
// so the `AddHistory` will return nil error forever.
func SetHistoryPath(fp string) {
	ins := getInstance()
	cfg := ins.Config.Clone()
	cfg.HistoryFile = fp
	ins.SetConfig(cfg)
}

// set auto completer to global instance
func SetAutoComplete(completer AutoCompleter) {
	ins := getInstance()
	cfg := ins.Config.Clone()
	cfg.AutoComplete = completer
	ins.SetConfig(cfg)
}

// add history to global instance manually
// raise error only if `SetHistoryPath` is set with a non-empty path
func AddHistory(content string) error {
	ins := getInstance()
	return ins.SaveHistory(content)
}

func Password(prompt string) ([]byte, error) {
	ins := getInstance()
	return ins.ReadPassword(prompt)
}

// readline with global configs
func Line(prompt string) (string, error) {
	ins := getInstance()
	ins.SetPrompt(prompt)
	return ins.Readline()
}

type CancelableStdin struct {
	mutex  sync.Mutex
	stop   chan struct{}
	notify chan struct{}
	data   []byte
	read   int
	err    error
}

func NewCancelableStdin() *CancelableStdin {
	c := &CancelableStdin{
		notify: make(chan struct{}),
		stop:   make(chan struct{}),
	}
	go c.ioloop()
	return c
}

func (c *CancelableStdin) ioloop() {
loop:
	for {
		select {
		case <-c.notify:
			c.read, c.err = os.Stdin.Read(c.data)
			<-c.notify
		case <-c.stop:
			break loop
		}
	}
}

func (c *CancelableStdin) Read(b []byte) (n int, err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data = b
	c.notify <- struct{}{}
	select {
	case <-c.notify:
		return c.read, c.err
	case <-c.stop:
		return 0, io.EOF
	}
}

func (c *CancelableStdin) Close() error {
	close(c.stop)
	return nil
}
