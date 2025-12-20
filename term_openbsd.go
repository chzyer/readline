// +build openbsd

package readline

import (
	"golang.org/x/sys/unix"
)

type Termios unix.Termios

func getTermios(fd int) (*Termios, error) {
	termios, err := unix.IoctlGetTermios(fd, unix.TIOCGETA)
	if err != nil {
		return nil, err
	}
	return (*Termios)(termios), nil
}

func setTermios(fd int, termios *Termios) error {
	return unix.IoctlSetTermios(fd, unix.TIOCSETA, (*unix.Termios)(termios))
}

func GetSize(fd int) (int, int, error) {
	uws, err := unix.IoctlGetWinsize(fd, unix.TIOCGWINSZ)
	if err != nil {
		return 0, 0, err
	}
	return int(uws.Col), int(uws.Row), nil
}

