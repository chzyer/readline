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
	w     io.Writer
	op    *Operation
	width int

	inCompleteMode  bool
	inSelectMode    bool
	candidate       [][]rune
	candidateSource []rune
	candidateOff    int
	candidateChoise int
	candidateColNum int

	// State for pagination
	maxLine int // Maximum allowed columns on a single terminal
	pageIdx int // Current page
}

func newOpCompleter(w io.Writer, op *Operation, width int) *opCompleter {
	return &opCompleter{
		w:       w,
		op:      op,
		width:   width,
		maxLine: 5,
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
	// Number of elements in a full screen
	numCandidatePerPage := o.candidateColNum * o.maxLine
	// Number of elements in the current screen, which could be non-full
	matrixSize := o.numCandidatesCurPage()
	pageStart := o.pageIdx * numCandidatePerPage

	candidateChoiceModPage := o.candidateChoise - pageStart
	candidateChoiceModPage += i
	candidateChoiceModPage %= matrixSize
	if candidateChoiceModPage < 0 {
		candidateChoiceModPage += matrixSize
	}

	o.candidateChoise = candidateChoiceModPage + pageStart
}

func (o *opCompleter) nextLine(i int) {
	// Number of elements in a full screen
	candidatesCurPage := o.numCandidatesCurPage()
	numCandidatesFullPage := o.maxLine * o.candidateColNum
	pageStart := o.pageIdx * numCandidatesFullPage

	// Number of elements in the current screen, which could be non-full
	numLines := o.getLinesCurPage()
	rectangleSize := numLines * o.candidateColNum

	candidateChoiceModPage := o.candidateChoise - pageStart
	candidateChoiceModPage += i * o.candidateColNum

	if candidateChoiceModPage >= candidatesCurPage {
		if candidateChoiceModPage < rectangleSize {
			candidateChoiceModPage += o.candidateColNum
		}
		candidateChoiceModPage -= rectangleSize
	} else if candidateChoiceModPage < 0 {
		candidateChoiceModPage += rectangleSize
		if candidateChoiceModPage > candidatesCurPage {
			candidateChoiceModPage -= o.candidateColNum
		}
	}

	o.candidateChoise = candidateChoiceModPage + pageStart
}

func (o *opCompleter) OnComplete() bool {
	if o.width == 0 {
		return false
	}
	if o.IsInCompleteSelectMode() {
		o.doSelect()
		return true
	}

	buf := o.op.buf
	rs := buf.Runes()

	if o.IsInCompleteMode() && o.candidateSource != nil && runes.Equal(rs, o.candidateSource) {
		o.EnterCompleteSelectMode()
		o.doSelect()
		return true
	}

	o.ExitCompleteSelectMode()
	o.candidateSource = rs
	newLines, offset := o.op.cfg.AutoComplete.Do(rs, buf.idx)
	if len(newLines) == 0 {
		o.ExitCompleteMode(false)
		return true
	}

	// only Aggregate candidates in non-complete mode
	if !o.IsInCompleteMode() {
		if len(newLines) == 1 {
			buf.WriteRunes(newLines[0])
			o.ExitCompleteMode(false)
			return true
		}

		same, size := runes.Aggregate(newLines)
		if size > 0 {
			buf.WriteRunes(same)
			o.ExitCompleteMode(false)
			return true
		}
	}

	o.EnterCompleteMode(offset, newLines)
	return true
}

func (o *opCompleter) IsInCompleteSelectMode() bool {
	return o.inSelectMode
}

func (o *opCompleter) IsInCompleteMode() bool {
	return o.inCompleteMode
}

func (o *opCompleter) HandleCompleteSelect(r rune) bool {
	next := true
	switch r {
	case CharEnter, CharCtrlJ:
		next = false
		o.op.buf.WriteRunes(o.op.candidate[o.op.candidateChoise])
		o.ExitCompleteMode(false)
	case CharLineStart:
		num := o.candidateChoise % o.candidateColNum
		o.nextCandidate(-num)
	case CharLineEnd:
		num := o.candidateColNum - o.candidateChoise%o.candidateColNum - 1
		o.candidateChoise += num
		if o.candidateChoise >= len(o.candidate) {
			o.candidateChoise = len(o.candidate) - 1
		}
	case CharBackspace:
		o.ExitCompleteSelectMode()
		next = false
	case CharTab, CharForward:
		o.doSelect()
	case CharBell, CharInterrupt:
		o.ExitCompleteMode(true)
		next = false
	case CharNext:
		o.nextLine(1)
	case CharBackward:
		o.nextCandidate(-1)
	case CharPrev:
		o.nextLine(-1)
	case CharK:
		o.updatePage(1)
	case CharJ:
		o.updatePage(-1)
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

// Number of candidate completions we can show on the current page,
// which might be different from number of candidates on a full page if
// we are on the last page.
func (o *opCompleter) numCandidatesCurPage() int {
	numCandidatePerPage := o.candidateColNum * o.maxLine
	pageStart := o.pageIdx * numCandidatePerPage
	if len(o.candidate)-pageStart >= numCandidatePerPage {
		return numCandidatePerPage
	}
	return len(o.candidate) - pageStart
}

// Number of lines on the current page
func (o *opCompleter) getLinesCurPage() int {
	curPageSize := o.numCandidatesCurPage()
	numLines := curPageSize / o.candidateColNum
	if curPageSize%o.candidateColNum != 0 {
		numLines += 1
	}
	return numLines
}

func (o *opCompleter) OnWidthChange(newWidth int) {
	o.width = newWidth
}

// Move page
func (o *opCompleter) updatePage(offset int) {
	if !o.inCompleteMode {
		return
	}

	nextPageIdx := o.pageIdx + offset
	if nextPageIdx < 0 {
		return
	}
	nextPageStart := nextPageIdx * o.candidateColNum * o.maxLine
	if nextPageStart > len(o.candidate) {
		return
	}
	o.pageIdx = nextPageIdx
	o.candidateChoise = nextPageStart
}

func (o *opCompleter) CompleteRefresh() {
	if !o.inCompleteMode {
		return
	}
	lineCnt := o.op.buf.CursorLineCount()
	colWidth := 0
	for _, c := range o.candidate {
		w := runes.WidthAll(c)
		if w > colWidth {
			colWidth = w
		}
	}
	colWidth += o.candidateOff + 1
	same := o.op.buf.RuneSlice(-o.candidateOff)

	// -1 to avoid reach the end of line
	width := o.width - 1
	colNum := width / colWidth
	if colNum != 0 {
		colWidth += (width - (colWidth * colNum)) / colNum
	}

	o.candidateColNum = colNum
	buf := bufio.NewWriter(o.w)
	buf.Write(bytes.Repeat([]byte("\n"), lineCnt))

	colIdx := 0
	lines := 1
	buf.WriteString("\033[J")

	// Compute the candidates to show on the current page
	numCandidatePerPage := o.candidateColNum * o.maxLine
	startIdx := o.pageIdx * numCandidatePerPage
	endIdx := (o.pageIdx + 1) * numCandidatePerPage
	if endIdx > len(o.candidate) {
		endIdx = len(o.candidate)
	}

	for idx := startIdx; idx < endIdx; idx += 1 {
		c := o.candidate[idx]
		inSelect := idx == o.candidateChoise && o.IsInCompleteSelectMode()
		if inSelect {
			buf.WriteString("\033[30;47m")
		}
		buf.WriteString(string(same))
		buf.WriteString(string(c))
		buf.Write(bytes.Repeat([]byte(" "), colWidth-runes.WidthAll(c)-runes.WidthAll(same)))

		if inSelect {
			buf.WriteString("\033[0m")
		}

		colIdx++
		if colIdx == colNum {
			buf.WriteString("\n")
			lines++
			colIdx = 0
		}
	}

	// Add an extra line for navigation instructions
	if colIdx != 0 {
		buf.WriteString("\n")
		lines++
	}
	navigationMsg := "(j: prev page)    (k: next page)"
	buf.WriteString(navigationMsg)
	buf.Write(bytes.Repeat([]byte(" "), width-len(navigationMsg)))

	// move back
	fmt.Fprintf(buf, "\033[%dA\r", lineCnt-1+lines)
	fmt.Fprintf(buf, "\033[%dC", o.op.buf.idx+o.op.buf.PromptLen())
	buf.Flush()
}

func (o *opCompleter) aggCandidate(candidate [][]rune) int {
	offset := 0
	for i := 0; i < len(candidate[0]); i++ {
		for j := 0; j < len(candidate)-1; j++ {
			if i > len(candidate[j]) {
				goto aggregate
			}
			if candidate[j][i] != candidate[j+1][i] {
				goto aggregate
			}
		}
		offset = i
	}
aggregate:
	return offset
}

func (o *opCompleter) EnterCompleteSelectMode() {
	o.inSelectMode = true
	o.candidateChoise = -1
	o.CompleteRefresh()
}

func (o *opCompleter) EnterCompleteMode(offset int, candidate [][]rune) {
	o.inCompleteMode = true
	o.candidate = candidate
	o.candidateOff = offset

	// Initialize for complete mode
	colWidth := 0
	for _, c := range candidate {
		w := runes.WidthAll(c)
		if w > colWidth {
			colWidth = w
		}
	}
	colWidth += offset + 1
	width := o.width - 1
	colNum := width / colWidth
	if colNum != 0 {
		colWidth += (width - (colWidth * colNum)) / colNum
	}
	o.candidateColNum = colNum
	o.pageIdx = 0

	o.CompleteRefresh()
}

func (o *opCompleter) ExitCompleteSelectMode() {
	o.inSelectMode = false
	o.candidate = nil
	o.candidateChoise = -1
	o.candidateOff = -1
	o.candidateSource = nil
}

func (o *opCompleter) ExitCompleteMode(revent bool) {
	o.inCompleteMode = false
	o.ExitCompleteSelectMode()
}
