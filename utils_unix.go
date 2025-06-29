// +build aix darwin dragonfly freebsd linux,!appengine netbsd openbsd os400 solaris

package readline

import (
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

// SuspendMe use to send suspend signal to myself, when we in the raw mode.
// For OSX it need to send to parent's pid
// For Linux it need to send to myself
func SuspendMe() {
	p, _ := os.FindProcess(os.Getppid())
	p.Signal(syscall.SIGTSTP)
	p, _ = os.FindProcess(os.Getpid())
	p.Signal(syscall.SIGTSTP)
}

// get width of the terminal
func getWidth(stdoutFd int) int {
	cols, _, err := GetSize(stdoutFd)
	if err != nil {
		return -1
	}
	return cols
}

func GetScreenWidth() int {
	w := getWidth(syscall.Stdout)
	if w < 0 {
		w = getWidth(syscall.Stderr)
	}
	return w
}

// getWidthHeight of the terminal using given file descriptor
func getWidthHeight(stdoutFd int) (width int, height int) {
	width, height, err := GetSize(stdoutFd)
	if err != nil {
		return -1, -1
	}
	return
}

// GetScreenSize returns the width/height of the terminal or -1,-1 or error
func GetScreenSize() (width int, height int) {
	width, height = getWidthHeight(syscall.Stdout)
	if width < 0 {
		width, height = getWidthHeight(syscall.Stderr)
	}
	return
}

// Ask the terminal for the current cursor position. The terminal will then
// write the position back to us via termainal stdin asynchronously.
func SendCursorPosition(t *Terminal) {
	t.Write([]byte("\033[6n"))
}

// ClearScreen clears the console screen
func ClearScreen(w io.Writer) (int, error) {
	return w.Write([]byte("\033[H"))
}

func DefaultIsTerminal() bool {
	return IsTerminal(syscall.Stdin) && (IsTerminal(syscall.Stdout) || IsTerminal(syscall.Stderr))
}

func GetStdin() int {
	return syscall.Stdin
}

// -----------------------------------------------------------------------------

var (
	sizeChange         sync.Once
	sizeChangeCallback func()
)

func DefaultOnWidthChanged(f func()) {
	DefaultOnSizeChanged(f)
}

func DefaultOnSizeChanged(f func()) {
	sizeChangeCallback = f
	sizeChange.Do(func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGWINCH)

		go func() {
			for {
				_, ok := <-ch
				if !ok {
					break
				}
				sizeChangeCallback()
			}
		}()
	})
}
