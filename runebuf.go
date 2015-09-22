package readline

import (
	"bytes"
	"io"
)

type RuneBuffer struct {
	buf         []rune
	idx         int
	prompt      []byte
	w           io.Writer
	hasPrompt   bool
	lastWritten int
}

func NewRuneBuffer(w io.Writer, prompt string) *RuneBuffer {
	rb := &RuneBuffer{
		prompt: []byte(prompt),
		w:      w,
	}
	return rb
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
	r.Refresh(-1, r.SetIdx(0))
}

func (r *RuneBuffer) MoveBackward() {
	if r.idx == 0 {
		return
	}
	r.idx--
	r.Refresh(0, -1)
}

func (rb *RuneBuffer) WriteString(s string) {
	rb.WriteRunes([]rune(s))
}

func (rb *RuneBuffer) WriteRune(r rune) {
	rb.WriteRunes([]rune{r})
}

func (rb *RuneBuffer) WriteRunes(r []rune) {
	tail := append(r, rb.buf[rb.idx:]...)
	rb.buf = append(rb.buf[:rb.idx], tail...)
	rb.idx++
	rb.Refresh(1, 1)
}

func (r *RuneBuffer) MoveForward() {
	if r.idx == len(r.buf) {
		return
	}
	r.idx++
	r.Refresh(0, 1)
}

func (r *RuneBuffer) Delete() {
	if r.idx == len(r.buf) {
		return
	}
	r.buf = append(r.buf[:r.idx], r.buf[r.idx+1:]...)
	r.Refresh(-1, 0)
}

func (r *RuneBuffer) DeleteWord() {
	if r.idx == len(r.buf) {
		return
	}
	init := r.idx
	for init < len(r.buf) && r.buf[init] == ' ' {
		init++
	}
	for i := init + 1; i < len(r.buf); i++ {
		if r.buf[i] != ' ' && r.buf[i-1] == ' ' {
			r.buf = append(r.buf[:r.idx], r.buf[i-1:]...)
			r.Refresh(r.idx-i+1, 0)
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
		if r.buf[i] != ' ' && r.buf[i-1] == ' ' {
			r.Refresh(0, r.SetIdx(i))
			return
		}
	}
	r.Refresh(0, r.SetIdx(0))
}

func (r *RuneBuffer) SetIdx(idx int) (change int) {
	i := r.idx
	r.idx = idx
	return r.idx - i
}

func (r *RuneBuffer) Kill() {
	length := len(r.buf)
	r.buf = r.buf[:r.idx]
	r.Refresh(r.idx-length, 0)
}

func (r *RuneBuffer) MoveToNextWord() {
	for i := r.idx + 1; i < len(r.buf); i++ {
		if r.buf[i] != ' ' && r.buf[i-1] == ' ' {
			r.Refresh(0, r.SetIdx(i))
			return
		}
	}
	r.Refresh(0, r.SetIdx(len(r.buf)))
}

func (r *RuneBuffer) BackEscapeWord() {
	if r.idx == 0 {
		return
	}
	for i := r.idx - 1; i > 0; i-- {
		if r.buf[i] != ' ' && r.buf[i-1] == ' ' {
			change := i - r.idx
			r.buf = append(r.buf[:i], r.buf[r.idx:]...)
			r.Refresh(change, r.SetIdx(i))
			return
		}
	}

	length := len(r.buf)
	r.buf = r.buf[:0]
	r.Refresh(-length, r.SetIdx(0))
}

func (r *RuneBuffer) Backspace() {
	if r.idx == 0 {
		return
	}
	r.idx--
	r.buf = append(r.buf[:r.idx], r.buf[r.idx+1:]...)
	r.Refresh(-1, -1)
}

func (r *RuneBuffer) MoveToLineEnd() {
	if r.idx == len(r.buf) {
		return
	}
	r.Refresh(0, r.SetIdx(len(r.buf)))
}

func (r *RuneBuffer) Refresh(chlen, chidx int) {
	s := r.Output(len(r.buf)-chlen, r.idx-chidx, true)
	r.w.Write(s)
}

func (r *RuneBuffer) RefreshSet(originLength, originIdx int) {
	r.w.Write(r.Output(originLength, originIdx, true))
}

func (r *RuneBuffer) Output(originLength, originIdx int, prompt bool) []byte {
	buf := bytes.NewBuffer(nil)
	buf.Write(r.CleanOutput())
	buf.Write(r.prompt)
	r.hasPrompt = prompt
	buf.Write([]byte(string(r.buf)))
	buf.Write(bytes.Repeat([]byte{'\b'}, len(r.buf)-r.idx))
	return buf.Bytes()
}

func (r *RuneBuffer) CleanOutput() []byte {
	buf := bytes.NewBuffer(nil)
	buf.Write([]byte("\033[J"))
	for i := 0; i <= 100; i++ {
		buf.WriteString("\033[2K\r\b")
	}
	return buf.Bytes()
}

func (r *RuneBuffer) Clean() {
	r.w.Write(r.CleanOutput())
}

func (r *RuneBuffer) ResetScreen() {
	r.w.Write(r.Output(0, 0, false))
}

func (r *RuneBuffer) Reset() []rune {
	ret := r.buf
	r.buf = r.buf[:0]
	r.idx = 0
	r.hasPrompt = false
	return ret
}

func (r *RuneBuffer) Set(buf []rune) {
	length, idx := len(r.buf), r.idx
	r.buf = buf
	r.idx = len(r.buf)
	r.RefreshSet(length, idx)
}
