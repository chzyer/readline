// +build windows

package readline

import (
	"fmt"
	"io"
	"syscall"
)

func SuspendMe() {
}

func GetStdin() int {
	return int(syscall.Stdin)
}

func init() {
	isWindows = true
}

// get width of the terminal
func GetScreenWidth() int {
	info, _ := GetConsoleScreenBufferInfo()
	if info == nil {
		return -1
	}
	return int(info.dwSize.x)
}

// GetScreenSize returns the width, height of the terminal or -1,-1
func GetScreenSize() (width int, height int) {
	info, _ := GetConsoleScreenBufferInfo()
	if info == nil {
		return -1, -1
	}
	height = int(info.srWindow.bottom) - int(info.srWindow.top) + 1
	width = int(info.srWindow.right) - int(info.srWindow.left) + 1
	return
}

// Send the Current cursor position to t.sizeChan.
func SendCursorPosition(t *Terminal) {
	info, err := GetConsoleScreenBufferInfo()
	if err != nil || info == nil {
		t.sizeChan <- "-1;-1"
	} else {
		t.sizeChan <- fmt.Sprintf("%d;%d", info.dwCursorPosition.y, info.dwCursorPosition.x)
	}
}

// ClearScreen clears the console screen
func ClearScreen(_ io.Writer) error {
	return SetConsoleCursorPosition(&_COORD{0, 0})
}

func DefaultIsTerminal() bool {
	return true
}

func DefaultOnWidthChanged(f func()) {
	DefaultOnSizeChanged(f)
}

func DefaultOnSizeChanged(f func()) {

}
