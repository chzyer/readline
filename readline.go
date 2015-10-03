package readline

import "io"

type Instance struct {
	Terminal  *Terminal
	Operation *Operation
}

type Config struct {
	Prompt       string
	HistoryFile  string
	AutoComplete AutoCompleter
	VimMode      bool
	Stdout       io.Writer
	Stderr       io.Writer

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
	return nil
}

func NewEx(cfg *Config) (*Instance, error) {
	t, err := NewTerminal(cfg)
	if err != nil {
		return nil, err
	}
	rl := t.Readline()
	return &Instance{
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

func (i *Instance) Stdout() io.Writer {
	return i.Operation.Stdout()
}

func (i *Instance) Stderr() io.Writer {
	return i.Operation.Stderr()
}

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

func (i *Instance) ReadSlice() ([]byte, error) {
	return i.Operation.Slice()
}

func (i *Instance) Close() error {
	if err := i.Terminal.Close(); err != nil {
		return err
	}
	i.Operation.Close()
	return nil
}
