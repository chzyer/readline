package readline

import "io"

type Instance struct {
	Config    *Config
	Terminal  *Terminal
	Operation *Operation
}

type Config struct {
	// prompt supports ANSI escape sequence, so we can color some characters even in windows
	Prompt string

	// readline will persist historys to file where HistoryFile specified
	HistoryFile string
	// specify the max length of historys, it's 500 by default
	HistoryLimit int

	// AutoCompleter will called once user press TAB
	AutoComplete AutoCompleter

	// If VimMode is true, readline will in vim.insert mode by default
	VimMode bool

	Stdout io.Writer
	Stderr io.Writer

	inited bool
}

func (c *Config) Init() error {
	if c.inited {
		return nil
	}
	c.inited = true
	if c.Stdout == nil {
		c.Stdout = Stdout
	}
	if c.Stderr == nil {
		c.Stderr = Stderr
	}
	if c.HistoryLimit <= 0 {
		c.HistoryLimit = 500
	}
	return nil
}

func NewEx(cfg *Config) (*Instance, error) {
	t, err := NewTerminal(cfg)
	if err != nil {
		return nil, err
	}
	rl := t.Readline()
	return &Instance{
		Config:    cfg,
		Terminal:  t,
		Operation: rl,
	}, nil
}

func New(prompt string) (*Instance, error) {
	return NewEx(&Config{Prompt: prompt})
}

func (i *Instance) SetPrompt(s string) {
	i.Operation.SetPrompt(s)
}

// change hisotry persistence in runtime
func (i *Instance) SetHistoryPath(p string) {
	i.Operation.SetHistoryPath(p)
}

// readline will refresh automatic when write through Stdout()
func (i *Instance) Stdout() io.Writer {
	return i.Operation.Stdout()
}

// readline will refresh automatic when write through Stdout()
func (i *Instance) Stderr() io.Writer {
	return i.Operation.Stderr()
}

// switch VimMode in runtime
func (i *Instance) SetVimMode(on bool) {
	i.Operation.SetVimMode(on)
}

func (i *Instance) IsVimMode() bool {
	return i.Operation.IsEnableVimMode()
}

func (i *Instance) ReadPassword(prompt string) ([]byte, error) {
	return i.Operation.Password(prompt)
}

func (i *Instance) Readline() (string, error) {
	return i.Operation.String()
}

// same as readline
func (i *Instance) ReadSlice() ([]byte, error) {
	return i.Operation.Slice()
}

// we must make sure that call Close() before process exit.
func (i *Instance) Close() error {
	if err := i.Terminal.Close(); err != nil {
		return err
	}
	i.Operation.Close()
	return nil
}
