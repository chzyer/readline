package readline

import (
	"container/list"
	"fmt"
	"os"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/crypto/ssh/terminal"
)

// IsTerminal returns true if the given file descriptor is a terminal.
func IsTerminal(fd int) bool {
	return terminal.IsTerminal(fd)
}

func MakeRaw(fd int) (*terminal.State, error) {
	return terminal.MakeRaw(fd)
}

func Restore(fd int, state *terminal.State) error {
	return terminal.Restore(fd, state)
}

func IsPrintable(key rune) bool {
	isInSurrogateArea := key >= 0xd800 && key <= 0xdbff
	return key >= 32 && !isInSurrogateArea
}

func escapeExKey(r rune) rune {
	switch r {
	case 'D':
		r = CharBackward
	case 'C':
		r = CharForward
	case 'A':
		r = CharPrev
	case 'B':
		r = CharNext
	}
	return r
}

func escapeKey(r rune) rune {
	switch r {
	case 'b':
		r = MetaPrev
	case 'f':
		r = MetaNext
	case 'd':
		r = MetaDelete
	case CharBackspace:
		r = MetaBackspace
	case KeyEsc:

	}
	return r
}

func Debug(o ...interface{}) {
	f, _ := os.OpenFile("debug.tmp", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	fmt.Fprintln(f, o...)
	f.Close()
}

type winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

func getWidth() int {
	ws := &winsize{}
	retCode, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)))

	if int(retCode) == -1 {
		panic(errno)
	}
	return int(ws.Col)
}

func debugList(l *list.List) {
	idx := 0
	for e := l.Front(); e != nil; e = e.Next() {
		Debug(idx, fmt.Sprintf("%+v", e.Value))
		idx++
	}
}

func equalRunes(a, b []rune) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func sleep(n int) {
	Debug(n)
	time.Sleep(2000 * time.Millisecond)
}
