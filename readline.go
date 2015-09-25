package readline

import "io"

type Instance struct {
	t *Terminal
	o *Operation
}

type AutoCompleter interface {
	Do(line []rune, pos int) (newLine []rune, newPos int, ok bool)
}

type Config struct {
	Prompt       string
	HistoryFile  string
	AutoComplete AutoCompleter
}

func NewEx(cfg *Config) (*Instance, error) {
	t, err := NewTerminal(cfg)
	if err != nil {
		return nil, err
	}
	rl := t.Readline()
	return &Instance{
		t: t,
		o: rl,
	}, nil
}

func New(prompt string) (*Instance, error) {
	return NewEx(&Config{Prompt: prompt})
}

func (i *Instance) Stdout() io.Writer {
	return i.o.Stdout()
}

func (i *Instance) Stderr() io.Writer {
	return i.o.Stderr()
}

func (i *Instance) Readline() (string, error) {
	return i.o.String()
}

func (i *Instance) ReadSlice() ([]byte, error) {
	return i.o.Slice()
}

func (i *Instance) Close() error {
	if err := i.t.Close(); err != nil {
		return err
	}
	i.o.Close()
	return nil
}

func (i *Instance) Write(b []byte) (int, error) {
	return i.o.Stderr().Write(b)
}
