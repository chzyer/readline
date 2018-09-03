package readline

import "os"

// State contains the state of a terminal.
type State struct {
	consctl *os.File
}

// MakeRaw put the terminal connected to the given file descriptor into raw
// mode and returns the previous state of the terminal so that it can be
// restored.
func MakeRaw(fd int) (*State, error) {
	f, err := os.OpenFile("/dev/consctl", os.O_WRONLY, 0)
	if err != nil {
		return nil, err
	}
	if _, err := f.WriteString("rawon"); err != nil {
		return nil, err
	}
	return &State{consctl: f}, nil
}

// Restore restores the terminal connected to the given file descriptor to a
// previous state.
func restoreTerm(fd int, state *State) error {
	return state.consctl.Close()
}
