// +build windows

package readline

import (
	"io"

	"gopkg.in/bufio.v1"
)

func init() {
	Stdout = NewANSIWriter(Stdout)
	Stderr = NewANSIWriter(Stderr)
}

type ANSIWriter struct {
	io.Writer
}

func NewANSIWriter(w io.Writer) *ANSIWriter {
	a := &ANSIWriter{
		Writer: w,
	}
	return a
}

func (a *ANSIWriter) Write(b []byte) (int, error) {
	var (
		isEsc    bool
		isEscSeq bool
	)
	buf := bufio.NewBuffer(nil)
	for i := 0; i < len(b); i++ {
		if isEscSeq {
			isEsc = false
			isEscSeq = false
			continue
		}
		if isEsc {
			isEsc = false
			continue
		}

		switch b[i] {
		case CharEsc:
			isEsc = true
			continue
		case '[':
			if isEsc {
				isEscSeq = true
				continue
			}
			fallthrough
		default:
		}
	}
	n, err := buf.WriteTo(a.Writer)
	if err != nil {
		return n, err
	}
	return len(b), nil
}
