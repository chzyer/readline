// +build windows

package readline

import (
	"io"
	"strconv"
	"strings"
	"sync"
	"unicode/utf8"
	"unsafe"

	"gopkg.in/bufio.v1"
)

func init() {
	Stdout = NewANSIWriter(Stdout)
	Stderr = NewANSIWriter(Stderr)
}

type ANSIWriter struct {
	target io.Writer
	ch     chan rune
	wg     sync.WaitGroup
	sync.Mutex
}

func NewANSIWriter(w io.Writer) *ANSIWriter {
	a := &ANSIWriter{
		target: w,
		ch:     make(chan rune, 1024),
	}
	go a.ioloop()
	return a
}

func (a *ANSIWriter) Close() error {
	close(a.ch)
	a.wg.Wait()
	return nil
}

func (a *ANSIWriter) ioloop() {
	a.wg.Add(1)
	defer a.wg.Done()

	var (
		ok       bool
		isEsc    bool
		isEscSeq bool

		char rune
		arg  []string

		target = bufio.NewWriter(a.target)
	)

	peek := func() rune {
		select {
		case ch := <-a.ch:
			return ch
		default:
			return 0
		}
	}

read:
	r := char
	if char == 0 {
		r, ok = <-a.ch
		if !ok {
			target.Flush()
			return
		}
	} else {
		char = 0
	}

	if isEscSeq {
		isEscSeq = a.ioloopEscSeq(target, r, &arg)
		goto read
	}

	switch r {
	case CharEsc:
		isEsc = true
	case '[':
		if isEsc {
			arg = nil
			isEscSeq = true
			isEsc = false
			break
		}
		fallthrough
	default:
		target.WriteRune(r)
		char = peek()
		if char == 0 || char == CharEsc {
			target.Flush()
		}
	}
	goto read
}

func (a *ANSIWriter) ioloopEscSeq(w *bufio.Writer, r rune, argptr *[]string) bool {
	arg := *argptr
	var err error
	switch r {
	case 'J':
		eraseLine()
	case 'K':
		eraseLine()
	case 'm':
		color := word(0)
		for _, item := range arg {
			var c int
			c, err = strconv.Atoi(item)
			if err != nil {
				w.WriteString("[" + strings.Join(arg, ";") + "m")
				break
			}
			if c >= 30 && c < 40 {
				color |= ColorTableFg[c-30]
			} else if c == 0 {
				color = ColorTableFg[7]
			}
		}
		if err != nil {
			break
		}
		kernel.SetConsoleTextAttribute(stdout, uintptr(color))
	case 'A':
	case 'B':
	case 'C':
	case 'D':
	case '\007':
	case ';':
		if len(arg) == 0 || arg[len(arg)-1] != "" {
			arg = append(arg, "")
		}
		fallthrough
	default:
		if len(arg) == 0 {
			arg = append(arg, "")
		}
		arg[len(arg)-1] += string(r)
		*argptr = arg
		return true
	}
	*argptr = nil
	return false
}

func (a *ANSIWriter) Write(b []byte) (int, error) {
	a.Lock()
	defer a.Unlock()

	off := 0
	for len(b) > off {
		r, size := utf8.DecodeRune(b[off:])
		if size == 0 {
			return off, io.ErrShortWrite
		}
		off += size
		a.ch <- r
	}
	return off, nil
}

func eraseLine() error {
	sbi, err := GetConsoleScreenBufferInfo()
	if err != nil {
		return err
	}

	var written int
	return kernel.FillConsoleOutputCharacterW(stdout, uintptr(' '),
		uintptr(sbi.dwSize.x-sbi.dwCursorPosition.x),
		sbi.dwCursorPosition.ptr(),
		uintptr(unsafe.Pointer(&written)),
	)
}

const (
	_                = uint16(0)
	COLOR_FBLUE      = 0x0001
	COLOR_FGREEN     = 0x0002
	COLOR_FRED       = 0x0004
	COLOR_FINTENSITY = 0x0008

	COLOR_BBLUE      = 0x0010
	COLOR_BGREEN     = 0x0020
	COLOR_BRED       = 0x0040
	COLOR_BINTENSITY = 0x0080
)

var ColorTableFg = []word{
	0,                                       // 30: Black
	COLOR_FRED,                              // 31: Red
	COLOR_FGREEN,                            // 32: Green
	COLOR_FRED | COLOR_FGREEN,               // 33: Yellow
	COLOR_FBLUE,                             // 34: Blue
	COLOR_FRED | COLOR_FBLUE,                // 35: Magenta
	COLOR_FGREEN | COLOR_FBLUE,              // 36: Cyan
	COLOR_FRED | COLOR_FBLUE | COLOR_FGREEN, // 37: White
}

var ColorTableBg = []word{
	0,                                       // 40: Black
	COLOR_BRED,                              // 41: Red
	COLOR_BGREEN,                            // 42: Green
	COLOR_BRED | COLOR_BGREEN,               // 43: Yellow
	COLOR_BBLUE,                             // 44: Blue
	COLOR_BRED | COLOR_BBLUE,                // 45: Magenta
	COLOR_BGREEN | COLOR_BBLUE,              // 46: Cyan
	COLOR_BRED | COLOR_BBLUE | COLOR_BGREEN, // 47: White
}
