package readline

import (
	"bytes"
	"io"
)

type RuneBuffer struct {
	buf    []rune
	idx    int
	prompt []rune
	w      io.Writer

	cleanInScreen bool
}

func NewRuneBuffer(w io.Writer, prompt string) *RuneBuffer {
	rb := &RuneBuffer{
		w: w,
	}
	rb.SetPrompt(prompt)
	return rb
}

func (r *RuneBuffer) CurrentWidth(x int) int {
	return RunesWidth(r.buf[:x])
}

func (r *RuneBuffer) PromptLen() int {
	return RunesWidth(RunesColorFilter(r.prompt))
}

func (r *RuneBuffer) RuneSlice(i int) []rune {
	if i > 0 {
		rs := make([]rune, i)
		copy(rs, r.buf[r.idx:r.idx+i])
		return rs
	}
	rs := make([]rune, -i)
	copy(rs, r.buf[r.idx+i:r.idx])
	return rs
}

func (r *RuneBuffer) Runes() []rune {
	newr := make([]rune, len(r.buf))
	copy(newr, r.buf)
	return newr
}

func (r *RuneBuffer) Pos() int {
	return r.idx
}

func (r *RuneBuffer) Len() int {
	return len(r.buf)
}

func (r *RuneBuffer) MoveToLineStart() {
	r.Refresh(func() {
		if r.idx == 0 {
			return
		}
		r.idx = 0
	})
}

func (r *RuneBuffer) MoveBackward() {
	r.Refresh(func() {
		if r.idx == 0 {
			return
		}
		r.idx--
	})
}

func (r *RuneBuffer) WriteString(s string) {
	r.WriteRunes([]rune(s))
}

func (r *RuneBuffer) WriteRune(s rune) {
	r.WriteRunes([]rune{s})
}

func (r *RuneBuffer) WriteRunes(s []rune) {
	r.Refresh(func() {
		tail := append(s, r.buf[r.idx:]...)
		r.buf = append(r.buf[:r.idx], tail...)
		r.idx += len(s)
	})
}

func (r *RuneBuffer) MoveForward() {
	r.Refresh(func() {
		if r.idx == len(r.buf) {
			return
		}
		r.idx++
	})
}

func (r *RuneBuffer) Delete() {
	r.Refresh(func() {
		if r.idx == len(r.buf) {
			return
		}
		r.buf = append(r.buf[:r.idx], r.buf[r.idx+1:]...)
	})
}

func (r *RuneBuffer) DeleteWord() {
	if r.idx == len(r.buf) {
		return
	}
	init := r.idx
	for init < len(r.buf) && IsWordBreak(r.buf[init]) {
		init++
	}
	for i := init + 1; i < len(r.buf); i++ {
		if !IsWordBreak(r.buf[i]) && IsWordBreak(r.buf[i-1]) {
			r.Refresh(func() {
				r.buf = append(r.buf[:r.idx], r.buf[i-1:]...)
			})
			return
		}
	}
	r.Kill()
}

func (r *RuneBuffer) MoveToPrevWord() {
	r.Refresh(func() {
		if r.idx == 0 {
			return
		}

		for i := r.idx - 1; i > 0; i-- {
			if !IsWordBreak(r.buf[i]) && IsWordBreak(r.buf[i-1]) {
				r.idx = i
				return
			}
		}
		r.idx = 0
	})
}

func (r *RuneBuffer) KillFront() {
	r.Refresh(func() {
		if r.idx == 0 {
			return
		}

		length := len(r.buf) - r.idx
		copy(r.buf[:length], r.buf[r.idx:])
		r.idx = 0
		r.buf = r.buf[:length]
	})
}

func (r *RuneBuffer) Kill() {
	r.Refresh(func() {
		r.buf = r.buf[:r.idx]
	})
}

func (r *RuneBuffer) Transpose() {
	r.Refresh(func() {
		if len(r.buf) == 1 {
			r.idx++
		}

		if len(r.buf) < 2 {
			return
		}

		if r.idx == 0 {
			r.idx = 1
		} else if r.idx >= len(r.buf) {
			r.idx = len(r.buf) - 1
		}
		r.buf[r.idx], r.buf[r.idx-1] = r.buf[r.idx-1], r.buf[r.idx]
		r.idx++
	})
}

func (r *RuneBuffer) MoveToNextWord() {
	r.Refresh(func() {
		for i := r.idx + 1; i < len(r.buf); i++ {
			if !IsWordBreak(r.buf[i]) && IsWordBreak(r.buf[i-1]) {
				r.idx = i
				return
			}
		}

		r.idx = len(r.buf)
	})
}

func (r *RuneBuffer) BackEscapeWord() {
	r.Refresh(func() {
		if r.idx == 0 {
			return
		}
		for i := r.idx - 1; i > 0; i-- {
			if !IsWordBreak(r.buf[i]) && IsWordBreak(r.buf[i-1]) {
				r.buf = append(r.buf[:i], r.buf[r.idx:]...)
				r.idx = i
				return
			}
		}

		r.buf = r.buf[:0]
		r.idx = 0
	})
}

func (r *RuneBuffer) Backspace() {
	r.Refresh(func() {
		if r.idx == 0 {
			return
		}

		r.idx--
		r.buf = append(r.buf[:r.idx], r.buf[r.idx+1:]...)
	})
}

func (r *RuneBuffer) MoveToLineEnd() {
	r.Refresh(func() {
		if r.idx == len(r.buf) {
			return
		}

		r.idx = len(r.buf)
	})
}

func (r *RuneBuffer) LineCount() int {
	return LineCount(RunesWidth(r.buf) + r.PromptLen())
}

func (r *RuneBuffer) IdxLine() int {
	totalWidth := RunesWidth(r.buf[:r.idx]) + r.PromptLen()
	w := getWidth()
	if w == 0 {
		return 0
	}
	line := totalWidth / w

	// if cursor is in last colmun and not any character behind it
	// the cursor will in the first line, otherwise will in the second line
	// this situation only occurs in golang's Stdout
	// TODO: figure out why
	if totalWidth%w == 0 && len(r.buf) == r.idx && !isWindows {
		line--
	}

	return line
}

func (r *RuneBuffer) CursorLineCount() int {
	return r.LineCount() - r.IdxLine()
}

func (r *RuneBuffer) Refresh(f func()) {
	r.Clean()
	if f != nil {
		f()
	}
	r.w.Write(r.output())
	r.cleanInScreen = false
}

func (r *RuneBuffer) output() []byte {
	buf := bytes.NewBuffer(nil)
	buf.WriteString(string(r.prompt))
	buf.Write([]byte(string(r.buf)))
	if len(r.buf) > r.idx {
		buf.Write(bytes.Repeat([]byte{'\b'}, RunesWidth(r.buf[r.idx:])))
	}
	return buf.Bytes()
}

func (r *RuneBuffer) Reset() []rune {
	ret := r.buf
	r.buf = r.buf[:0]
	r.idx = 0
	return ret
}

func (r *RuneBuffer) calWidth(m int) int {
	if m > 0 {
		return RunesWidth(r.buf[r.idx : r.idx+m])
	}
	return RunesWidth(r.buf[r.idx+m : r.idx])
}

func (r *RuneBuffer) SetStyle(start, end int, style string) {
	if end < start {
		panic("end < start")
	}

	// goto start
	move := start - r.idx
	if move > 0 {
		r.w.Write([]byte(string(r.buf[r.idx : r.idx+move])))
	} else {
		r.w.Write(bytes.Repeat([]byte("\b"), r.calWidth(move)))
	}
	r.w.Write([]byte("\033[" + style + "m"))
	r.w.Write([]byte(string(r.buf[start:end])))
	r.w.Write([]byte("\033[0m"))
	// TODO: move back
}

func (r *RuneBuffer) SetWithIdx(idx int, buf []rune) {
	r.Refresh(func() {
		r.buf = buf
		r.idx = idx
	})
}

func (r *RuneBuffer) Set(buf []rune) {
	r.SetWithIdx(len(buf), buf)
}

func (r *RuneBuffer) SetPrompt(prompt string) {
	r.prompt = []rune(prompt)
}

func (r *RuneBuffer) cleanOutput() []byte {
	buf := bytes.NewBuffer(nil)
	buf.Write([]byte("\033[J")) // just like ^k :)

	idxLine := r.IdxLine()

	if idxLine == 0 {
		buf.WriteString("\033[2K\r")
		return buf.Bytes()
	}
	for i := 0; i < idxLine; i++ {
		buf.WriteString("\033[2K\r\033[A")
	}
	buf.WriteString("\033[2K\r")
	return buf.Bytes()
}

func (r *RuneBuffer) Clean() {
	if r.cleanInScreen {
		return
	}
	r.cleanInScreen = true
	r.w.Write(r.cleanOutput())
}
