package readline

import (
	"io"
	"io/ioutil"
	"strconv"
	"syscall"
)

func SuspendMe() {
}

func GetScreenWidth() int {
	const w = 0

	// $COLS is set by vt(1). Read the environment variable directly
	// because Go might be caching the value.
	// See https://github.com/golang/go/issues/25234
	b, err := ioutil.ReadFile("/env/COLS")
	if err != nil {
		return w
	}
	cols, err := strconv.Atoi(string(b))
	if err != nil {
		return w
	}
	return cols
}

// ClearScreen clears the console screen
func ClearScreen(w io.Writer) (int, error) {
	return w.Write([]byte("\033[H"))
}

func DefaultIsTerminal() bool {
	return true
}

func GetStdin() int {
	return syscall.Stdin
}

func DefaultOnWidthChanged(f func()) {
}
