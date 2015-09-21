package readline

import "io"

type Instance struct {
	t *Terminal
	o *Operation
}

func New(prompt string) (*Instance, error) {
	t, err := NewTerminal()
	if err != nil {
		return nil, err
	}
	rl := t.Readline(prompt)
	return &Instance{
		t: t,
		o: rl,
	}, nil
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
	return i.t.Close()
}
