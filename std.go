package readline

import (
	"io"
	"os"
)

var (
	Stdout io.WriteCloser = os.Stdout
	Stderr io.WriteCloser = os.Stderr
)
