package readline

import (
	"strconv"
	"syscall"
	"unicode"

	"golang.org/x/crypto/ssh/terminal"
)

var (
	StdinFd   = int(uintptr(syscall.Stdin))
	isWindows = false
)

// IsTerminal returns true if the given file descriptor is a terminal.
func IsTerminal(fd int) bool {
	return terminal.IsTerminal(fd)
}

func MakeRaw(fd int) (*terminal.State, error) {
	return terminal.MakeRaw(fd)
}

func Restore(fd int, state *terminal.State) error {
	err := terminal.Restore(fd, state)
	if err != nil {
		// errno 0 means everything is ok :)
		if err.Error() == "errno 0" {
			err = nil
		}
	}
	return nil
}

func IsPrintable(key rune) bool {
	isInSurrogateArea := key >= 0xd800 && key <= 0xdbff
	return key >= 32 && !isInSurrogateArea
}

// translate Esc[X
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

// translate EscX to Meta+X
func escapeKey(r rune) rune {
	switch r {
	case 'b':
		r = MetaBackward
	case 'f':
		r = MetaForward
	case 'd':
		r = MetaDelete
	case CharTranspose:
		r = MetaTranspose
	case CharBackspace:
		r = MetaBackspace
	case CharEsc:

	}
	return r
}

func RunesEqual(a, b []rune) bool {
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

// calculate how many lines for N character
func LineCount(w int) int {
	screenWidth := getWidth()
	r := w / screenWidth
	if w%screenWidth != 0 {
		r++
	}
	return r
}

// Search in runes from end to front
func RunesIndexBck(r, sub []rune) int {
	for i := len(r) - len(sub); i >= 0; i-- {
		found := true
		for j := 0; j < len(sub); j++ {
			if r[i+j] != sub[j] {
				found = false
				break
			}
		}
		if found {
			return i
		}
	}
	return -1
}

// Search in runes from front to end
func RunesIndex(r, sub []rune) int {
	for i := 0; i < len(r); i++ {
		found := true
		if len(r[i:]) < len(sub) {
			return -1
		}
		for j := 0; j < len(sub); j++ {
			if r[i+j] != sub[j] {
				found = false
				break
			}
		}
		if found {
			return i
		}
	}
	return -1
}

func IsWordBreak(i rune) bool {
	if i >= 'a' && i <= 'z' {
		return false
	}
	if i >= 'A' && i <= 'Z' {
		return false
	}
	return true
}

var zeroWidth = []*unicode.RangeTable{
	unicode.Mn,
	unicode.Me,
	unicode.Cc,
	unicode.Cf,
}

var doubleWidth = []*unicode.RangeTable{
	unicode.Han,
	unicode.Hangul,
	unicode.Hiragana,
	unicode.Katakana,
}

func RuneIndex(r rune, rs []rune) int {
	for i := 0; i < len(rs); i++ {
		if rs[i] == r {
			return i
		}
	}
	return -1
}

func RunesColorFilter(r []rune) []rune {
	newr := make([]rune, 0, len(r))
	for pos := 0; pos < len(r); pos++ {
		if r[pos] == '\033' && r[pos+1] == '[' {
			idx := RuneIndex('m', r[pos+2:])
			if idx == -1 {
				continue
			}
			pos += idx + 2
			continue
		}
		newr = append(newr, r[pos])
	}
	return newr
}

func RuneWidth(r rune) int {
	if unicode.IsOneOf(zeroWidth, r) {
		return 0
	}
	if unicode.IsOneOf(doubleWidth, r) {
		return 2
	}
	return 1
}

func RunesWidth(r []rune) (length int) {
	for i := 0; i < len(r); i++ {
		length += RuneWidth(r[i])
	}
	return
}

func RunesAggregate(candicate [][]rune) (same []rune, size int) {
	for i := 0; i < len(candicate[0]); i++ {
		for j := 0; j < len(candicate)-1; j++ {
			if i >= len(candicate[j]) || i >= len(candicate[j+1]) {
				goto aggregate
			}
			if candicate[j][i] != candicate[j+1][i] {
				goto aggregate
			}
		}
		size = i + 1
	}
aggregate:
	if size > 0 {
		same = RunesCopy(candicate[0][:size])
		for i := 0; i < len(candicate); i++ {
			n := RunesCopy(candicate[i])
			copy(n, n[size:])
			candicate[i] = n[:len(n)-size]
		}
	}
	return
}

func RunesCopy(r []rune) []rune {
	n := make([]rune, len(r))
	copy(n, r)
	return n
}

func RunesHasPrefix(r, prefix []rune) bool {
	if len(r) < len(prefix) {
		return false
	}
	return RunesEqual(r[:len(prefix)], prefix)
}

func GetInt(s []string, def int) int {
	if len(s) == 0 {
		return def
	}
	c, err := strconv.Atoi(s[0])
	if err != nil {
		return def
	}
	return c
}
