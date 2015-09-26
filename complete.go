package readline

import (
	"bytes"
	"fmt"
	"io"
)

type AutoCompleter interface {
	Do(line []rune, pos int) (newLine [][]rune, offset int)
}

type opCompleter struct {
	w  io.Writer
	op *Operation
	ac AutoCompleter

	inCompleteMode  bool
	inSelectMode    bool
	candicate       [][]rune
	candicateSource []rune
	candicateOff    int
	candicateChoise int
	candicateColNum int
}

func newOpCompleter(w io.Writer, op *Operation) *opCompleter {
	return &opCompleter{
		w:  w,
		op: op,
		ac: op.cfg.AutoComplete,
	}
}

func (o *opCompleter) doSelect() {
	if len(o.candicate) == 1 {
		o.op.buf.WriteRunes(o.candicate[0])
		o.ExitCompleteMode(false)
		return
	}
	o.nextCandicate(1)
	o.CompleteRefresh()
}

func (o *opCompleter) nextCandicate(i int) {
	o.candicateChoise += i
	o.candicateChoise = o.candicateChoise % len(o.candicate)
	if o.candicateChoise < 0 {
		o.candicateChoise = len(o.candicate) + o.candicateChoise
	}
}

func (o *opCompleter) OnComplete() {
	if o.IsInCompleteSelectMode() {
		o.doSelect()
		return
	}

	buf := o.op.buf
	rs := buf.Runes()

	if o.IsInCompleteMode() && EqualRunes(rs, o.candicateSource) {
		o.EnterCompleteSelectMode()
		o.doSelect()
		return
	}

	o.ExitCompleteSelectMode()
	o.candicateSource = rs
	newLines, offset := o.ac.Do(rs, buf.idx)
	if len(newLines) == 0 {
		o.ExitCompleteMode(false)
		return
	}

	// only Aggregate candicates in non-complete mode
	if !o.IsInCompleteMode() {
		if len(newLines) == 1 {
			buf.WriteRunes(newLines[0])
			o.ExitCompleteMode(false)
			return
		}

		same, size := AggRunes(newLines)
		if size > 0 {
			buf.WriteRunes(same)
			o.ExitCompleteMode(false)
			return
		}
	}

	o.EnterCompleteMode(offset, newLines)
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
		o.op.buf.WriteRunes(o.op.candicate[o.op.candicateChoise])
		o.ExitCompleteMode(false)
	case CharLineStart:
		num := o.candicateChoise % o.candicateColNum
		o.nextCandicate(-num)
	case CharLineEnd:
		num := o.candicateColNum - o.candicateChoise%o.candicateColNum - 1
		o.candicateChoise += num
		if o.candicateChoise >= len(o.candicate) {
			o.candicateChoise = len(o.candicate) - 1
		}
	case CharBackspace:
		o.ExitCompleteSelectMode()
		next = false
	case CharTab, CharForward:
		o.doSelect()
	case CharCancel, CharInterrupt:
		o.ExitCompleteMode(true)
		next = false
	case CharNext:
		tmpChoise := o.candicateChoise + o.candicateColNum
		if tmpChoise >= o.getMatrixSize() {
			tmpChoise -= o.getMatrixSize()
		} else if tmpChoise >= len(o.candicate) {
			tmpChoise += o.candicateColNum
			tmpChoise -= o.getMatrixSize()
		}
		o.candicateChoise = tmpChoise
	case CharBackward:
		o.nextCandicate(-1)
	case CharPrev:
		tmpChoise := o.candicateChoise - o.candicateColNum
		if tmpChoise < 0 {
			tmpChoise += o.getMatrixSize()
			if tmpChoise >= len(o.candicate) {
				tmpChoise -= o.candicateColNum
			}
		}
		o.candicateChoise = tmpChoise
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

func (o *opCompleter) getMatrixSize() int {
	line := len(o.candicate) / o.candicateColNum
	if len(o.candicate)%o.candicateColNum != 0 {
		line++
	}
	return line * o.candicateColNum
}

func (o *opCompleter) CompleteRefresh() {
	if !o.inCompleteMode {
		return
	}
	lineCnt := o.op.buf.CursorLineCount()
	colWidth := 0
	for _, c := range o.candicate {
		w := RunesWidth(c)
		if w > colWidth {
			colWidth = w
		}
	}
	colNum := getWidth() / (colWidth + o.candicateOff + 2)
	o.candicateColNum = colNum
	buf := bytes.NewBuffer(nil)
	buf.Write(bytes.Repeat([]byte("\n"), lineCnt))
	same := o.op.buf.RuneSlice(-o.candicateOff)
	colIdx := 0
	lines := 1
	buf.WriteString("\033[J")
	for idx, c := range o.candicate {
		inSelect := idx == o.candicateChoise && o.IsInCompleteSelectMode()
		if inSelect {
			buf.WriteString("\033[30;47m")
		}
		buf.WriteString(string(same))
		buf.WriteString(string(c))
		buf.Write(bytes.Repeat([]byte(" "), colWidth-len(c)))
		if inSelect {
			buf.WriteString("\033[0m")
		}

		buf.WriteString("  ")
		colIdx++
		if colIdx == colNum {
			buf.WriteString("\n")
			lines++
			colIdx = 0
		}
	}

	// move back
	fmt.Fprintf(buf, "\033[%dA\r", lineCnt-1+lines)
	fmt.Fprintf(buf, "\033[%dC", o.op.buf.idx+o.op.buf.PromptLen())
	o.w.Write(buf.Bytes())
}

func (o *opCompleter) aggCandicate(candicate [][]rune) int {
	offset := 0
	for i := 0; i < len(candicate[0]); i++ {
		for j := 0; j < len(candicate)-1; j++ {
			if i > len(candicate[j]) {
				goto aggregate
			}
			if candicate[j][i] != candicate[j+1][i] {
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
	o.candicateChoise = -1
	o.CompleteRefresh()
}

func (o *opCompleter) EnterCompleteMode(offset int, candicate [][]rune) {
	o.inCompleteMode = true
	o.candicate = candicate
	o.candicateOff = offset
	o.CompleteRefresh()
}

func (o *opCompleter) ExitCompleteSelectMode() {
	o.inSelectMode = false
	o.candicate = nil
	o.candicateChoise = -1
	o.candicateOff = -1
	o.candicateSource = nil
}

func (o *opCompleter) ExitCompleteMode(revent bool) {
	o.inCompleteMode = false
	o.ExitCompleteSelectMode()
}
