// +build windows

package readline

import (
	"bytes"
	"testing"
)

func TestANSIWriterColorBounds(t *testing.T) {
	orig := kernel.SetConsoleTextAttribute
	defer func() {
		kernel.SetConsoleTextAttribute = orig
	}()
	kernel.SetConsoleTextAttribute = func(u ...uintptr) error {
		return nil
	}

	buf := bytes.NewBuffer(nil)
	w := NewANSIWriter(buf)

	tests := []string{
		"\x1b[30m", "\x1b[37m", "\x1b[38m", "\x1b[39m",
		"\x1b[40m", "\x1b[47m", "\x1b[48m", "\x1b[49m",
		"\x1b[38;5;196m", "\x1b[48;5;12m",
	}

	for _, tt := range tests {
		if _, err := w.Write([]byte(tt)); err != nil {
			t.Fatalf("unexpected error for %q: %v", tt, err)
		}
	}
}
