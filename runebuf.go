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
}

func NewRuneBuffer(w io.Writer, prompt string) *RuneBuffer {
	rb := &RuneBuffer{
		prompt: []rune(prompt),
		w:      w,
	}
	return rb
}

func (r *RuneBuffer) PromptLen() int {
	return RunesWidth(r.prompt)
}

func (r *RuneBuffer) Runes() []rune {
	return r.buf
}

func (r *RuneBuffer) Pos() int {
	return r.idx
}

func (r *RuneBuffer) Len() int {
	return len(r.buf)
}

func (r *RuneBuffer) MoveToLineStart() {
	if r.idx == 0 {
		return
	}
	r.idx = 0
	r.Refresh()
}

func (r *RuneBuffer) MoveBackward() {
	if r.idx == 0 {
		return
	}
	r.idx--
	r.Refresh()
}

func (r *RuneBuffer) WriteString(s string) {
	r.WriteRunes([]rune(s))
}

func (r *RuneBuffer) WriteRune(s rune) {
	r.WriteRunes([]rune{s})
}

func (r *RuneBuffer) WriteRunes(s []rune) {
	tail := append(s, r.buf[r.idx:]...)
	r.buf = append(r.buf[:r.idx], tail...)
	r.idx++
	r.Refresh()
}

func (r *RuneBuffer) MoveForward() {
	if r.idx == len(r.buf) {
		return
	}
	r.idx++
	r.Refresh()
}

func (r *RuneBuffer) Delete() {
	if r.idx == len(r.buf) {
		return
	}
	r.buf = append(r.buf[:r.idx], r.buf[r.idx+1:]...)
	r.Refresh()
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
			r.buf = append(r.buf[:r.idx], r.buf[i-1:]...)
			r.Refresh()
			return
		}
	}
	r.Kill()
}

func (r *RuneBuffer) MoveToPrevWord() {
	if r.idx == 0 {
		return
	}
	for i := r.idx - 1; i > 0; i-- {
		if !IsWordBreak(r.buf[i]) && IsWordBreak(r.buf[i-1]) {
			r.idx = i
			r.Refresh()
			return
		}
	}
	r.idx = 0
	r.Refresh()
}

func (r *RuneBuffer) SetIdx(idx int) (change int) {
	i := r.idx
	r.idx = idx
	return r.idx - i
}

func (r *RuneBuffer) Kill() {
	r.buf = r.buf[:r.idx]
	r.Refresh()
}

func (r *RuneBuffer) Transpose() {
	if len(r.buf) < 2 {
		if len(r.buf) == 1 {
			r.idx++
			r.Refresh()
		}
		return
	}
	if r.idx == 0 {
		r.idx = 1
	} else if r.idx >= len(r.buf) {
		r.idx = len(r.buf) - 1
	}
	r.buf[r.idx], r.buf[r.idx-1] = r.buf[r.idx-1], r.buf[r.idx]
	r.idx++
	r.Refresh()
}

func (r *RuneBuffer) MoveToNextWord() {
	for i := r.idx + 1; i < len(r.buf); i++ {
		if !IsWordBreak(r.buf[i]) && IsWordBreak(r.buf[i-1]) {
			r.idx = i
			r.Refresh()
			return
		}
	}
	r.idx = len(r.buf)
	r.Refresh()
}

func (r *RuneBuffer) BackEscapeWord() {
	if r.idx == 0 {
		return
	}
	for i := r.idx - 1; i > 0; i-- {
		if !IsWordBreak(r.buf[i]) && IsWordBreak(r.buf[i-1]) {
			r.buf = append(r.buf[:i], r.buf[r.idx:]...)
			r.idx = i
			r.Refresh()
			return
		}
	}

	r.buf = r.buf[:0]
	r.idx = 0
	r.Refresh()
}

func (r *RuneBuffer) Backspace() {
	if r.idx == 0 {
		return
	}
	r.idx--
	r.buf = append(r.buf[:r.idx], r.buf[r.idx+1:]...)
	r.Refresh()
}

func (r *RuneBuffer) MoveToLineEnd() {
	if r.idx == len(r.buf) {
		return
	}
	r.idx = len(r.buf)
	r.Refresh()
}

func (r *RuneBuffer) LineCount() int {
	return LineCount(RunesWidth(r.buf) + r.PromptLen())
}

func (r *RuneBuffer) IdxLine() int {
	totalWidth := RunesWidth(r.buf[:r.idx]) + r.PromptLen()
	w := getWidth()
	line := 0
	for totalWidth >= w {
		totalWidth -= w
		line++
	}
	return line
}

func (r *RuneBuffer) CursorLineCount() int {
	return r.LineCount() - r.IdxLine()
}

func (r *RuneBuffer) Refresh() {
	r.w.Write(r.Output())
}

func (r *RuneBuffer) Output() []byte {
	buf := bytes.NewBuffer(nil)
	buf.Write(r.CleanOutput())
	buf.WriteString(string(r.prompt))
	buf.Write([]byte(string(r.buf)))
	if len(r.buf) > r.idx {
		buf.Write(bytes.Repeat([]byte{'\b'}, len(r.buf)-r.idx))
	}
	return buf.Bytes()
}

func (r *RuneBuffer) CleanOutput() []byte {
	buf := bytes.NewBuffer(nil)
	buf.Write([]byte("\033[J")) // just like ^k :)

	// TODO: calculate how many line before cursor.
	for i := 0; i <= 100; i++ {
		buf.WriteString("\033[2K\r\b")
	}
	return buf.Bytes()
}

func (r *RuneBuffer) Clean() {
	r.w.Write(r.CleanOutput())
}

func (r *RuneBuffer) Reset() []rune {
	ret := r.buf
	r.buf = r.buf[:0]
	r.idx = 0
	return ret
}

func (r *RuneBuffer) SetStyle(start, end int, style string) {
	idx := r.idx
	if end < start {
		panic("end < start")
	}

	// goto start
	move := start - idx
	if move > 0 {
		r.w.Write([]byte(string(r.buf[r.idx : r.idx+move])))
	} else {
		r.w.Write(bytes.Repeat([]byte("\b"), -move))
	}
	r.w.Write([]byte("\033[" + style))
	r.w.Write([]byte(string(r.buf[start:end])))
	r.w.Write([]byte("\033[0m"))
	if move > 0 {
		r.w.Write(bytes.Repeat([]byte("\b"), -move+(end-start)))
	} else if -move < end-start {
		r.w.Write(bytes.Repeat([]byte("\b"), -move))
	} else {
		r.w.Write([]byte(string(r.buf[end:r.idx])))
	}
}

func (r *RuneBuffer) SetWithIdx(idx int, buf []rune) {
	r.buf = buf
	r.idx = idx
	r.Refresh()
}

func (r *RuneBuffer) Set(buf []rune) {
	r.SetWithIdx(len(buf), buf)
}
