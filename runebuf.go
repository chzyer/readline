package readline

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
	"sync"
)

type runeBufferBck struct {
	buf []rune
	idx int
}

type RuneBuffer struct {
	buf    []rune
	idx    int
	prompt []rune
	w      io.Writer

	interactive bool
	cfg         *Config

	width int

	bck *runeBufferBck

	offset string          // is offset useful? scrolling means row varies
	ppos   int             // prompt start position (0 == column 1)

	lastKill []rune

	sync.Mutex
}

func (r *RuneBuffer) pushKill(text []rune) {
	r.lastKill = append([]rune{}, text...)
}

func (r *RuneBuffer) OnWidthChange(newWidth int) {
	r.Lock()
	r.width = newWidth
	r.Unlock()
}

func (r *RuneBuffer) Backup() {
	r.Lock()
	r.bck = &runeBufferBck{r.buf, r.idx}
	r.Unlock()
}

func (r *RuneBuffer) Restore() {
	r.Refresh(func() {
		if r.bck == nil {
			return
		}
		r.buf = r.bck.buf
		r.idx = r.bck.idx
	})
}

func NewRuneBuffer(w io.Writer, prompt string, cfg *Config, width int) *RuneBuffer {
	rb := &RuneBuffer{
		w:           w,
		interactive: cfg.useInteractive(),
		cfg:         cfg,
		width:       width,
	}
	rb.SetPrompt(prompt)
	return rb
}

func (r *RuneBuffer) SetConfig(cfg *Config) {
	r.Lock()
	r.cfg = cfg
	r.interactive = cfg.useInteractive()
	r.Unlock()
}

func (r *RuneBuffer) SetMask(m rune) {
	r.Lock()
	r.cfg.MaskRune = m
	r.Unlock()
}

func (r *RuneBuffer) CurrentWidth(x int) int {
	r.Lock()
	defer r.Unlock()
	return runes.WidthAll(r.buf[:x])
}

func (r *RuneBuffer) PromptLen() int {
	r.Lock()
	width := r.promptLen()
	r.Unlock()
	return width
}

func (r *RuneBuffer) promptLen() int {
	return runes.WidthAll(runes.ColorFilter(r.prompt))
}

func (r *RuneBuffer) RuneSlice(i int) []rune {
	r.Lock()
	defer r.Unlock()

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
	r.Lock()
	newr := make([]rune, len(r.buf))
	copy(newr, r.buf)
	r.Unlock()
	return newr
}

func (r *RuneBuffer) Pos() int {
	r.Lock()
	defer r.Unlock()
	return r.idx
}

func (r *RuneBuffer) Len() int {
	r.Lock()
	defer r.Unlock()
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
	r.Lock()
	defer r.Unlock()

	if r.idx == len(r.buf) {
		// cursor is already at end of buf data so just call
		// append instead of refesh to save redrawing.
		r.buf = append(r.buf, s...)
		r.idx += len(s)
		if r.interactive {
			r.append(s)
		}
	} else {
		// writing into the data somewhere so do a refresh
		r.refresh(func() {
			tail := append(s, r.buf[r.idx:]...)
			r.buf = append(r.buf[:r.idx], tail...)
			r.idx += len(s)
		})
	}
}

func (r *RuneBuffer) MoveForward() {
	r.Refresh(func() {
		if r.idx == len(r.buf) {
			return
		}
		r.idx++
	})
}

func (r *RuneBuffer) IsCursorInEnd() bool {
	r.Lock()
	defer r.Unlock()
	return r.idx == len(r.buf)
}

func (r *RuneBuffer) Replace(ch rune) {
	r.Refresh(func() {
		r.buf[r.idx] = ch
	})
}

func (r *RuneBuffer) Erase() {
	r.Refresh(func() {
		r.idx = 0
		r.pushKill(r.buf[:])
		r.buf = r.buf[:0]
	})
}

func (r *RuneBuffer) Delete() (success bool) {
	r.Refresh(func() {
		if r.idx == len(r.buf) {
			return
		}
		r.pushKill(r.buf[r.idx : r.idx+1])
		r.buf = append(r.buf[:r.idx], r.buf[r.idx+1:]...)
		success = true
	})
	return
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
			r.pushKill(r.buf[r.idx : i-1])
			r.Refresh(func() {
				r.buf = append(r.buf[:r.idx], r.buf[i-1:]...)
			})
			return
		}
	}
	r.Kill()
}

func (r *RuneBuffer) MoveToPrevWord() (success bool) {
	r.Refresh(func() {
		if r.idx == 0 {
			return
		}

		for i := r.idx - 1; i > 0; i-- {
			if !IsWordBreak(r.buf[i]) && IsWordBreak(r.buf[i-1]) {
				r.idx = i
				success = true
				return
			}
		}
		r.idx = 0
		success = true
	})
	return
}

func (r *RuneBuffer) KillFront() {
	r.Refresh(func() {
		if r.idx == 0 {
			return
		}

		length := len(r.buf) - r.idx
		r.pushKill(r.buf[:r.idx])
		copy(r.buf[:length], r.buf[r.idx:])
		r.idx = 0
		r.buf = r.buf[:length]
	})
}

func (r *RuneBuffer) Kill() {
	r.Refresh(func() {
		r.pushKill(r.buf[r.idx:])
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

func (r *RuneBuffer) MoveToEndWord() {
	r.Refresh(func() {
		// already at the end, so do nothing
		if r.idx == len(r.buf) {
			return
		}
		// if we are at the end of a word already, go to next
		if !IsWordBreak(r.buf[r.idx]) && IsWordBreak(r.buf[r.idx+1]) {
			r.idx++
		}

		// keep going until at the end of a word
		for i := r.idx + 1; i < len(r.buf); i++ {
			if IsWordBreak(r.buf[i]) && !IsWordBreak(r.buf[i-1]) {
				r.idx = i - 1
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
				r.pushKill(r.buf[i:r.idx])
				r.buf = append(r.buf[:i], r.buf[r.idx:]...)
				r.idx = i
				return
			}
		}

		r.buf = r.buf[:0]
		r.idx = 0
	})
}

func (r *RuneBuffer) Yank() {
	if len(r.lastKill) == 0 {
		return
	}
	r.Refresh(func() {
		buf := make([]rune, 0, len(r.buf)+len(r.lastKill))
		buf = append(buf, r.buf[:r.idx]...)
		buf = append(buf, r.lastKill...)
		buf = append(buf, r.buf[r.idx:]...)
		r.buf = buf
		r.idx += len(r.lastKill)
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
	r.Lock()
	defer r.Unlock()
	if r.idx == len(r.buf) {
		return
	}
	r.refresh(func() {
		r.idx = len(r.buf)
	})
}

func (r *RuneBuffer) LineCount(width int) int {
	if width == -1 {
		width = r.width
	}
	return LineCount(width,
		runes.WidthAll(r.buf)+r.PromptLen())
}

func (r *RuneBuffer) MoveTo(ch rune, prevChar, reverse bool) (success bool) {
	r.Refresh(func() {
		if reverse {
			for i := r.idx - 1; i >= 0; i-- {
				if r.buf[i] == ch {
					r.idx = i
					if prevChar {
						r.idx++
					}
					success = true
					return
				}
			}
			return
		}
		for i := r.idx + 1; i < len(r.buf); i++ {
			if r.buf[i] == ch {
				r.idx = i
				if prevChar {
					r.idx--
				}
				success = true
				return
			}
		}
	})
	return
}

func (r *RuneBuffer) isInLineEdge() bool {
	if isWindows {
		return false
	}
	sp := r.getSplitByLine(r.buf, 1)
	return len(sp[len(sp)-1]) == 0  // last line is 0 len
}

func (r *RuneBuffer) getSplitByLine(rs []rune, nextWidth int) [][]rune {
	if r.cfg.EnableMask {
		w := runes.Width(r.cfg.MaskRune)
		masked := []rune(strings.Repeat(string(r.cfg.MaskRune), len(rs)))
		return SplitByLine(runes.ColorFilter(r.prompt), masked, r.ppos, r.width, w)
	} else {
		return SplitByLine(runes.ColorFilter(r.prompt), rs, r.ppos, r.width, nextWidth)
	}
}

func (r *RuneBuffer) IdxLine(width int) int {
	r.Lock()
	defer r.Unlock()
	return r.idxLine(width)
}

func (r *RuneBuffer) idxLine(width int) int {
	if width == 0 {
		return 0
	}
	nextWidth := 1
	if r.idx < len(r.buf) {
		nextWidth = runes.Width(r.buf[r.idx])
	}
	sp := r.getSplitByLine(r.buf[:r.idx], nextWidth)
	return len(sp) - 1
}

func (r *RuneBuffer) CursorLineCount() int {
	return r.LineCount(r.width) - r.IdxLine(r.width)
}

func (r *RuneBuffer) Refresh(f func()) {
	r.Lock()
	defer r.Unlock()
	r.refresh(f)
}

func (r *RuneBuffer) refresh(f func()) {
	if !r.interactive {
		if f != nil {
			f()
		}
		return
	}

	r.clean()
	if f != nil {
		f()
	}
	r.print()
}

// getAndSetOffset queries the terminal for the current cursor position by
// writing a control sequence to the terminal. This call is asynchronous
// and it returns before any offset has actually been set as the terminal
// will write the offset back to us via stdin and there may already be 
// other data in the stdin buffer ahead of it.
// This function is called at the start of readline each time.
func (r *RuneBuffer) getAndSetOffset(t *Terminal) {
	if !r.interactive {
		return
	}
	if !isWindows {
		// Handle lineedge cases where existing text before before
		// the prompt is printed would leave us at the right edge of
		// the screen but the next character would actually be printed
		// at the beginning of the next line.
		r.w.Write([]byte(" \b"))
	}
	t.GetOffset(r.setOffset)
}

func (r *RuneBuffer) SetOffset(offset string) {
	r.Lock()
	defer r.Unlock()
	r.setOffset(offset)
}

func (r *RuneBuffer) setOffset(offset string) {
	r.offset = offset
	if _, c, ok := (&escapeKeyPair{attr:offset}).Get2(); ok && c > 0 && c < r.width {
		r.ppos = c - 1  // c should be 1..width
	} else {
		r.ppos = 0
	}
}

// append s to the end of the current output. append is called in
// place of print() when clean() was avoided. As output is appended on
// the end, the cursor also needs no extra adjustment.
// NOTE: assumes len(s) >= 1 which should always be true for append.
func (r *RuneBuffer) append(s []rune) {
	buf := bytes.NewBuffer(nil)
	slen := len(s)
	if r.cfg.EnableMask {
		if slen > 1 && r.cfg.MaskRune != 0 {
			// write a mask character for all runes except the last rune
			buf.WriteString(strings.Repeat(string(r.cfg.MaskRune), slen-1))
		}
		// for the last rune, write \n or mask it otherwise.
		if s[slen-1] == '\n' {
			buf.WriteRune('\n')
		} else if r.cfg.MaskRune != 0 {
			buf.WriteRune(r.cfg.MaskRune)
		}
	} else {
		for _, e := range r.cfg.Painter.Paint(s, slen) {
			if e == '\t' {
				buf.WriteString(strings.Repeat(" ", TabWidth))
			} else {
				buf.WriteRune(e)
			}
		}
	}
	if r.isInLineEdge() {
		buf.WriteString(" \b")
	}
	r.w.Write(buf.Bytes())
}

// Print writes out the prompt and buffer contents at the current cursor position
func (r *RuneBuffer) Print() {
	r.Lock()
	defer r.Unlock()
	if !r.interactive {
		return
	}
	r.print()
}

func (r *RuneBuffer) print() {
	r.w.Write(r.output())
}

func (r *RuneBuffer) output() []byte {
	buf := bytes.NewBuffer(nil)
	buf.WriteString(string(r.prompt))
	if r.cfg.EnableMask && len(r.buf) > 0 {
		if r.cfg.MaskRune != 0 {
			buf.WriteString(strings.Repeat(string(r.cfg.MaskRune), len(r.buf)-1))
		}
		if r.buf[len(r.buf)-1] == '\n' {
			buf.WriteRune('\n')
		} else if r.cfg.MaskRune != 0 {
			buf.WriteRune(r.cfg.MaskRune)
		}
	} else {
		for _, e := range r.cfg.Painter.Paint(r.buf, r.idx) {
			if e == '\t' {
				buf.WriteString(strings.Repeat(" ", TabWidth))
			} else {
				buf.WriteRune(e)
			}
		}
	}
	if r.isInLineEdge() {
		buf.WriteString(" \b")
	}
	// cursor position
	if len(r.buf) > r.idx {
		buf.Write(r.getBackspaceSequence())
	}
	return buf.Bytes()
}

func (r *RuneBuffer) getBackspaceSequence() []byte {
	bcnt := len(r.buf) - r.idx  // backwards count to index
	sp := r.getSplitByLine(r.buf, 1)

	// Calculate how many lines up to the index line
	up := 0
	spi := len(sp) - 1
	for spi >= 0 {
		bcnt -= len(sp[spi])
		if bcnt <= 0 {
			break
		}
		up++
		spi--
	}

	// Calculate what column the index should be set to
	column := 1
	if spi == 0 {
		column += r.ppos
	}
	for _, rune := range sp[spi] {
		if bcnt >= 0 {
			break
		}
		column += runes.Width(rune)
		bcnt++
	}

	buf := bytes.NewBuffer(nil)
	if up > 0 {
		fmt.Fprintf(buf, "\033[%dA", up) // move cursor up to index line
	}
	fmt.Fprintf(buf, "\033[%dG", column) // move cursor to column

	return buf.Bytes()
}

func (r *RuneBuffer) Reset() []rune {
	ret := runes.Copy(r.buf)
	r.buf = r.buf[:0]
	r.idx = 0
	return ret
}

func (r *RuneBuffer) calWidth(m int) int {
	if m > 0 {
		return runes.WidthAll(r.buf[r.idx : r.idx+m])
	}
	return runes.WidthAll(r.buf[r.idx+m : r.idx])
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
	r.Lock()
	r.prompt = []rune(prompt)
	r.Unlock()
}

func (r *RuneBuffer) cleanOutput(w io.Writer, idxLine int) {
	buf := bufio.NewWriter(w)

	if r.width == 0 {
		buf.WriteString(strings.Repeat("\r\b", len(r.buf)+r.promptLen()))
		buf.Write([]byte("\033[J"))
	} else {
		if idxLine > 0 {
			fmt.Fprintf(buf, "\033[%dA", idxLine) // move cursor up by idxLine
		}
		fmt.Fprintf(buf, "\033[%dG", r.ppos + 1) // move cursor back to initial ppos position
		buf.Write([]byte("\033[J"))  // clear from cursor to end of screen
	}
	buf.Flush()
	return
}

func (r *RuneBuffer) Clean() {
	r.Lock()
	r.clean()
	r.Unlock()
}

func (r *RuneBuffer) clean() {
	r.cleanWithIdxLine(r.idxLine(r.width))
}

func (r *RuneBuffer) cleanWithIdxLine(idxLine int) {
	if !r.interactive {
		return
	}
	r.cleanOutput(r.w, idxLine)
}
