// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin || dragonfly || freebsd || (linux && !appengine) || netbsd || openbsd
// +build darwin dragonfly freebsd linux,!appengine netbsd openbsd

package readline

import (
	"syscall"

	"golang.org/x/sys/unix"
)

type Termios syscall.Termios

// GetSize returns the dimensions of the given terminal.
func GetSize(fd int) (int, int, error) {
	winsize, err := unix.IoctlGetWinsize(fd, unix.TIOCGWINSZ)
	if err != nil {
		return 0, 0, err
	}
	return int(winsize.Col), int(winsize.Row), nil
}
