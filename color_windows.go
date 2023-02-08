//go:build windows

/*
Copyright (c) Jason Walton <dev@lucid.thedremaing.org> (https://www.thedreaming.org)
Copyright (c) Sindre Sorhus <sindresorhus@gmail.com> (https://sindresorhus.com)

Released under the MIT License:

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package readline

import (
	"sync"

	"golang.org/x/sys/windows"
)

var (
	ansiSuccess bool
	ansiOnce    sync.Once
)

func enableANSI() bool {
	ansiOnce.Do(func() {
		ansiSuccess = realEnableANSI()
	})
	return ansiSuccess
}

// this is `enableColor` from https://github.com/jwalton/go-supportscolor
func realEnableANSI() bool {
	handle, err := windows.GetStdHandle(windows.STD_OUTPUT_HANDLE)
	if err != nil {
		return false
	}

	// Get the existing console mode.
	var mode uint32
	err = windows.GetConsoleMode(handle, &mode)
	if err != nil {
		return false
	}

	// If ENABLE_VIRTUAL_TERMINAL_PROCESSING is not set, then set it.  This will
	// enable native ANSI color support from Windows.
	if mode&windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING != windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING {
		// Enable color.
		// See https://docs.microsoft.com/en-us/windows/console/console-virtual-terminal-sequences.
		mode = mode | windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING
		err = windows.SetConsoleMode(handle, mode)
		if err != nil {
			return false
		}
	}

	return true
}
