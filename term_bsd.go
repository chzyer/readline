// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin || dragonfly || freebsd || netbsd || openbsd
// +build darwin dragonfly freebsd netbsd openbsd

package readline

import (
	"golang.org/x/sys/unix"
)

func getTermios(fd int) (*Termios, error) {
	termios, err := unix.IoctlGetTermios(fd, unix.TIOCGETA)
	return (*Termios)(termios), err
}

func setTermios(fd int, termios *Termios) error {
	return unix.IoctlSetTermios(fd, unix.TIOCSETA, (*unix.Termios)(termios))
}
