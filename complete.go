package readline

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

type AutoCompleter interface {
	// Readline will pass the whole line and current offset to it
	// Completer need to pass all the candidates, and how long they shared the same characters in line
	// Example:
	//   [go, git, git-shell, grep]
	//   Do("g", 1) => ["o", "it", "it-shell", "rep"], 1
	//   Do("gi", 2) => ["t", "t-shell"], 2
	//   Do("git", 3) => ["", "-shell"], 3
	Do(line []rune, pos int) (newLine [][]rune, length int)
}

type TabCompleter struct{}

func (t *TabCompleter) Do([]rune, int) ([][]rune, int) {
	return [][]rune{[]rune("\t")}, 0
}

type opCompleter struct {
	w               io.Writer
	op              *Operation
	width           int
	height          int

	inCompleteMode  bool
	inSelectMode    bool
	inPagerMode	bool

	candidate          [][]rune      // list of candidates
	candidateSource    []rune        // buffer string when tab was pressed
	candidateOff       int           // num runes in common from buf where candidate start
	candidateChoise    int           // candidate chosen (-1 for nothing yet) also used for paging
	candidateColNum    int           // num columns candidates take 0..wraps, 1 col, 2 cols etc.
	candidateColWidth  int           // width of candidate columns
}

func newOpCompleter(w io.Writer, op *Operation, width, height int) *opCompleter {
	return &opCompleter{
		w:      w,
		op:     op,
		width:  width,
		height: height,
	}
}

func (o *opCompleter) doSelect() {
	if len(o.candidate) == 1 {
		o.op.buf.WriteRunes(o.candidate[0])
		o.ExitCompleteMode(false)
		return
	}
	o.nextCandidate(1)
	o.CompleteRefresh()
}

func (o *opCompleter) nextCandidate(i int) {
	o.candidateChoise += i
	o.candidateChoise = o.candidateChoise % len(o.candidate)
	if o.candidateChoise < 0 {
		o.candidateChoise = len(o.candidate) + o.candidateChoise
	}
}

// OnComplete returns true if complete mode is available. Used to ring bell
// when tab pressed if cannot do complete for reason such as width unknown
// or no candidates available.
func (o *opCompleter) OnComplete() (ringBell bool) {
	if o.width == 0 || o.height < 3 {
		return false
	}
	if o.IsInCompleteSelectMode() {
		o.doSelect()
		return true
	}

	buf := o.op.buf
	rs := buf.Runes()

	// If in complete mode and nothing else typed then we must be entering select mode
	if o.IsInCompleteMode() && o.candidateSource != nil && runes.Equal(rs, o.candidateSource) {
		if len(o.candidate) > 1 {
			same, size := runes.Aggregate(o.candidate)
			if size > 0 {
				buf.WriteRunes(same)
				o.ExitCompleteMode(false)
				return false  // partial completion so ring the bell
			}
		}
		o.EnterCompleteSelectMode()
		o.doSelect()
		return true
	}

	newLines, offset := o.op.cfg.AutoComplete.Do(rs, buf.idx)
	if len(newLines) == 0 || (len(newLines) == 1 && len(newLines[0]) == 0) {
		o.ExitCompleteMode(false)
		return false // will ring bell on initial tab press
	} 
	if o.candidateOff > offset {
		// part of buffer we are completing has changed. Example might be that we were completing "ls" and
		// user typed space so we are no longer completing "ls" but now we are completing an argument of
		// the ls command. Instead of continuing in complete mode, we exit.
		o.ExitCompleteMode(false)
		return true
	}
	o.candidateSource = rs

	// only Aggregate candidates in non-complete mode
	if !o.IsInCompleteMode() {
		if len(newLines) == 1 {
			// not yet in complete mode but only 1 candidate so complete it
			buf.WriteRunes(newLines[0])
			o.ExitCompleteMode(false)
			return true
		}

		// check if all candidates have common prefix and return it and its size
		same, size := runes.Aggregate(newLines)
		if size > 0 {
			buf.WriteRunes(same)
			o.ExitCompleteMode(false)
			return false  // partial completion so ring the bell
		}
	}

	// otherwise, we just enter complete mode (which does a refresh)
	o.EnterCompleteMode(offset, newLines)
	return true
}

func (o *opCompleter) IsInCompleteSelectMode() bool {
	return o.inSelectMode
}

func (o *opCompleter) IsInCompleteMode() bool {
	return o.inCompleteMode
}

func (o *opCompleter) IsInPagerMode() bool {
	return o.inPagerMode
}

func (o *opCompleter) HandleCompleteSelect(r rune) (stayInMode bool) {
	next := true
	switch r {
	case CharEnter, CharCtrlJ:
		next = false
		o.op.buf.WriteRunes(o.op.candidate[o.op.candidateChoise])
		o.ExitCompleteMode(false)
	case CharLineStart:
		if o.candidateColNum > 1 {
			num := o.candidateChoise % o.candidateColNum
			o.nextCandidate(-num)
		}
	case CharLineEnd:
		if o.candidateColNum > 1 {
			num := o.candidateColNum - o.candidateChoise % o.candidateColNum - 1
			o.candidateChoise += num
			if o.candidateChoise >= len(o.candidate) {
				o.candidateChoise = len(o.candidate) - 1
			}
		}
	case CharBackspace:
		o.ExitCompleteSelectMode()
		next = false
	case CharTab, CharForward:
		o.nextCandidate(1)
	case CharBell, CharInterrupt:
		o.ExitCompleteMode(true)
		next = false
	case CharNext:
		colNum := 1
		if o.candidateColNum > 1 {
			colNum = o.candidateColNum
		}
		tmpChoise := o.candidateChoise + colNum
		if tmpChoise >= o.getMatrixSize() {
			tmpChoise -= o.getMatrixSize()
		} else if tmpChoise >= len(o.candidate) {
			tmpChoise += colNum
			tmpChoise -= o.getMatrixSize()
		}
		o.candidateChoise = tmpChoise
	case CharBackward:
		o.nextCandidate(-1)
	case CharPrev:
		colNum := 1
		if o.candidateColNum > 1 {
			colNum = o.candidateColNum 
		}
		tmpChoise := o.candidateChoise - colNum
		if tmpChoise < 0 {
			tmpChoise += o.getMatrixSize()
			if tmpChoise >= len(o.candidate) {
				tmpChoise -= colNum
			}
		}
		o.candidateChoise = tmpChoise
	default:
		next = false
		o.ExitCompleteSelectMode()
	}
	if next {
		o.CompleteRefresh()
		return true
	}
	return false
}

// HandlePagerMode handles user input when in pager mode.
// The user can only press certain keys either viewing another
// page or quitting the pager and going back to the prompt.
// returns true if we are still in pager mode or false if
// we exit pager mode.
func (o *opCompleter) HandlePagerMode(r rune) (stayInMode bool) {
	switch r {
	case ' ', 'Y', 'y':                // yes, show me more
		return o.pagerRefresh()    // last page exits
	case 'q','Q', 'N', 'n':            // no, quit giving me more
		o.scrollOutOfPagerMode()   // adjust for prompt
		o.ExitCompleteMode(true)   // completely exit complete mode
		return false
	default:                           // invalid choice
		o.op.t.Bell()              // ring bell
		return true                // stay in pager mode
	}
}

func (o *opCompleter) getMatrixSize() int {
	colNum := 1
	if o.candidateColNum > 1 {
		colNum = o.candidateColNum
	}
	line := len(o.candidate) / colNum
	if len(o.candidate) % colNum != 0 {
		line++
	}
	return line * colNum
}

func (o *opCompleter) OnWidthChange(newWidth int) {
	o.width = newWidth
}

func (o *opCompleter) OnSizeChange(newWidth, newHeight int) {
	o.width = newWidth
	o.height = newHeight
}

// setColumnInfo calculates column width and number of columns required
// to present the list of candidates on the terminal.
func (o *opCompleter) setColumnInfo() {
	same := o.op.buf.RuneSlice(-o.candidateOff)
	sameWidth := runes.WidthAll(same)

	colWidth := 0
	for _, c := range o.candidate {
		w := sameWidth + runes.WidthAll(c)
		if w > colWidth {
			colWidth = w
		}
	}
	colWidth++ // whitespace between cols

	// -1 to avoid end of line issues
	width := o.width - 1
	colNum := width / colWidth
	if colNum != 0 {
		colWidth += (width - (colWidth * colNum)) / colNum
	}

	o.candidateColNum = colNum
	o.candidateColWidth = colWidth
}

// needPagerMode returns true if number of candidates would go off the page
func (o *opCompleter) needPagerMode() bool {
	buflineCnt := o.op.buf.LineCount()           // lines taken by buffer content
	linesAvail := o.height - buflineCnt          // lines available without scrolling buffer off screen
	if o.candidateColNum > 0 {
		// Normal case where each candidate at least fits on a line
		maxOrPage := linesAvail * o.candidateColNum  // max candiates without needing to page
		return len(o.candidate) > maxOrPage
	}

	same := o.op.buf.RuneSlice(-o.candidateOff)
	sameWidth := runes.WidthAll(same)

	// 1 or more candidates take up multiple lines
	lines := 1
	for _, c := range o.candidate {
		cWidth := sameWidth + runes.WidthAll(c)
		cLines := 1
		if o.width > 0 {
			cLines = cWidth / o.width
			if cWidth % o.width > 0 {
				cLines++
			}
		}
		lines += cLines
		if lines > linesAvail {
			return true
		}
	}
	return false
}

// CompleteRefresh is used for completemode and selectmode
func (o *opCompleter) CompleteRefresh() {
	if !o.inCompleteMode {
		return
	}

	buf := bufio.NewWriter(o.w)
	// calculate num lines from cursor pos to where choices should be written
	lineCnt := o.op.buf.CursorLineCount()
	buf.Write(bytes.Repeat([]byte("\n"), lineCnt))  // move down from cursor to start of candidates
	buf.WriteString("\033[J")

	same := o.op.buf.RuneSlice(-o.candidateOff)

	colIdx := 0
	lines := 0
	sameWidth := runes.WidthAll(same)
	for idx, c := range o.candidate {
		inSelect := idx == o.candidateChoise && o.IsInCompleteSelectMode()
		cWidth := sameWidth + runes.WidthAll(c)
		cLines := 1
		if o.width > 0 {
			sWidth := 0
			if isWindows && inSelect {
				sWidth = 1 // adjust for hightlighting on Windows
			}
			cLines = (cWidth + sWidth) / o.width
			if (cWidth + sWidth) % o.width > 0 {
				cLines++
			}
		}
		if lines > 0 && colIdx == 0 {
			// After line 1, if we're printing to the first column
			// goto a new line. We do it here, instead of at the end
			// of the loop, to avoid the last \n taking up a blank
			// line at the end and stealing realestate.
			buf.WriteString("\n")
		}

		if inSelect {
			buf.WriteString("\033[30;47m")
		}

		buf.WriteString(string(same))
		buf.WriteString(string(c))
		if o.candidateColNum >= 1 {
			// only output spaces between columns if everything fits
			buf.Write(bytes.Repeat([]byte(" "), o.candidateColWidth - cWidth))
		}

		if inSelect {
			buf.WriteString("\033[0m")
		}

		colIdx++
		if colIdx >= o.candidateColNum {
			lines += cLines
			colIdx = 0
			if isWindows {
				// Windows EOL edge-case.
				buf.WriteString("\b")
			}
		}
	}
	if colIdx > 0 {
		lines++  // mid-line so count it.
	}

	// wrote out choices over "lines", move back to cursor (positioned at index)
	fmt.Fprintf(buf, "\033[%dA", lines)
	buf.Write(o.op.buf.getBackspaceSequence())
	buf.Flush()
}

// pagerRefresh writes a page full of candidates starting from candidateChoise to the screen
// followed by --More-- if there are more candidates or exiting complete mode entirely if done
// after which the prompt would be printed below the list (similar to bash). This is different
// to CompleteMode and CompleteSelectMode which leave the prompt on the same page and need
// to reset the cursor back up it after drawing the list.
func (o *opCompleter) pagerRefresh() (stayInMode bool) {
	buf := bufio.NewWriter(o.w)
	firstPage := o.candidateChoise == 0
	if firstPage {
		o.op.buf.SetOffset("1;1")     // paging, so reset any prompt offset
		// move down from cursor to where candidates should start
		lineCnt := o.op.buf.CursorLineCount()
		buf.Write(bytes.Repeat([]byte("\n"), lineCnt))
	} else {
		// after first page, redraw over --More--
		buf.WriteString("\r")
	}		
	buf.WriteString("\033[J")   // clear anything below

	same := o.op.buf.RuneSlice(-o.candidateOff)
	sameWidth := runes.WidthAll(same)

	colIdx := 0
	lines := 1
	for ; o.candidateChoise < len(o.candidate) ; o.candidateChoise++ {
		c := o.candidate[o.candidateChoise]
		cWidth := sameWidth + runes.WidthAll(c)
		cLines := 1
		if o.width > 0 {
			cLines = cWidth / o.width
			if cWidth % o.width > 0 {
				cLines++
			}
		}
		if lines > 1 && lines + cLines > o.height {
			break // won't fit on page, stop early.
		}
		buf.WriteString(string(same))
		buf.WriteString(string(c))
		if o.candidateColNum > 1 {
			// only output spaces between columns if more than 1
			buf.Write(bytes.Repeat([]byte(" "), o.candidateColWidth - cWidth))
		}
		colIdx++
		if colIdx >= o.candidateColNum {
			if isWindows {
				// Windows EOL edge-case.
				buf.WriteString("\b")
			}
			buf.WriteString("\n")
			lines += cLines
			colIdx = 0
		}
	}
	if colIdx != 0 {
		buf.WriteString("\n")
	}
	if firstPage || o.candidateChoise < len(o.candidate) {
		stayInMode = true
		buf.WriteString("--More--")
	} else {
		stayInMode = false
		o.scrollOutOfPagerMode()
		o.ExitCompleteMode(true)
	}
	buf.Flush()
	return
}

// scrollOutOfPagerMode adds enough new lines after the pager content such that when
// we rewrite the prompt it does not over write the page content. The code to rewrite
// the prompt assumes the cursor is at the index line, so we add enough blank lines.
func (o *opCompleter) scrollOutOfPagerMode() {
	lineCnt := o.op.buf.IdxLine(o.width)
	if lineCnt > 0 {
		buf := bufio.NewWriter(o.w)
		buf.Write(bytes.Repeat([]byte("\n"), lineCnt))
		buf.Flush()
	}
}

func (o *opCompleter) EnterCompleteSelectMode() {
	o.inSelectMode = true
	o.candidateChoise = -1
}

func (o *opCompleter) EnterCompleteMode(offset int, candidate [][]rune) {
	o.inCompleteMode = true
	o.candidate = candidate
	o.candidateOff = offset
	o.setColumnInfo()
	if o.needPagerMode() {
		o.EnterPagerMode()
	} else {
		o.CompleteRefresh()
	}
}

func (o *opCompleter) EnterPagerMode() {
	o.inPagerMode = true
	o.candidateChoise = 0   // next candidate to list on next page
	o.pagerRefresh()
}

func (o *opCompleter) ExitCompleteSelectMode() {
	o.inSelectMode = false
	o.candidateChoise = -1
}

func (o *opCompleter) ExitCompleteMode(revent bool) {
	o.inCompleteMode = false
	o.candidate = nil
	o.candidateOff = -1
	o.candidateSource = nil
	o.ExitCompleteSelectMode()
	o.ExitPagerMode()
}

func (o *opCompleter) ExitPagerMode() {
	o.inPagerMode = false
}
