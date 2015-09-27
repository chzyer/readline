package readline

import "io"

type Operation struct {
	cfg     *Config
	t       *Terminal
	buf     *RuneBuffer
	outchan chan []rune

	*opHistory
	*opSearch
	*opCompleter
}

type wrapWriter struct {
	r      *Operation
	t      *Terminal
	target io.Writer
}

func (w *wrapWriter) Write(b []byte) (int, error) {
	buf := w.r.buf
	buf.Clean()
	n, err := w.target.Write(b)
	if w.t.IsReading() {
		w.r.buf.Refresh()
	}
	if w.r.IsSearchMode() {
		w.r.SearchRefresh(-1)
	}
	if w.r.IsInCompleteMode() {
		w.r.CompleteRefresh()
	}
	return n, err
}

func NewOperation(t *Terminal, cfg *Config) *Operation {
	op := &Operation{
		cfg:       cfg,
		t:         t,
		buf:       NewRuneBuffer(t, cfg.Prompt),
		outchan:   make(chan []rune),
		opHistory: newOpHistory(cfg.HistoryFile),
	}
	op.opSearch = newOpSearch(op.buf.w, op.buf, op.opHistory)
	op.opCompleter = newOpCompleter(op.buf.w, op)
	go op.ioloop()
	return op
}

func (o *Operation) ioloop() {
	for {
		keepInSearchMode := false
		keepInCompleteMode := false
		r := o.t.ReadRune()

		if o.IsInCompleteSelectMode() {
			keepInCompleteMode = o.HandleCompleteSelect(r)
			if keepInCompleteMode {
				continue
			}

			o.buf.Refresh()
			switch r {
			case CharEnter, CharCtrlJ:
				o.UpdateHistory(o.buf.Runes(), false)
				fallthrough
			case CharInterrupt:
				o.t.KickRead()
				fallthrough
			case CharCancel:
				continue
			}
		}

		switch r {
		case CharCancel:
			if o.IsSearchMode() {
				o.ExitSearchMode(true)
				o.buf.Refresh()
			}
			if o.IsInCompleteMode() {
				o.ExitCompleteMode(true)
				o.buf.Refresh()
			}
		case CharTab:
			if o.opCompleter == nil {
				break
			}
			o.OnComplete()
			keepInCompleteMode = true
		case CharBckSearch:
			o.SearchMode(S_DIR_BCK)
			keepInSearchMode = true
		case CharFwdSearch:
			o.SearchMode(S_DIR_FWD)
			keepInSearchMode = true
		case CharKill:
			o.buf.Kill()
			keepInCompleteMode = true
		case MetaNext:
			o.buf.MoveToNextWord()
		case CharTranspose:
			o.buf.Transpose()
		case MetaPrev:
			o.buf.MoveToPrevWord()
		case MetaDelete:
			o.buf.DeleteWord()
		case CharLineStart:
			o.buf.MoveToLineStart()
		case CharLineEnd:
			o.buf.MoveToLineEnd()
		case CharDelete:
			o.buf.Delete()
		case CharBackspace, CharCtrlH:
			if o.IsSearchMode() {
				o.SearchBackspace()
				keepInSearchMode = true
				break
			}

			if o.buf.Len() == 0 {
				break
			}
			o.buf.Backspace()
			if o.IsInCompleteMode() {
				o.OnComplete()
			}
		case MetaBackspace, CharCtrlW:
			o.buf.BackEscapeWord()
		case CharEnter, CharCtrlJ:
			if o.IsSearchMode() {
				o.ExitSearchMode(false)
			}
			o.buf.MoveToLineEnd()
			o.buf.WriteRune('\n')
			data := o.buf.Reset()
			data = data[:len(data)-1] // trim \n
			o.outchan <- data
			o.NewHistory(data)
		case CharBackward:
			o.buf.MoveBackward()
		case CharForward:
			o.buf.MoveForward()
		case CharPrev:
			buf := o.PrevHistory()
			if buf != nil {
				o.buf.Set(buf)
			}
		case CharNext:
			buf, ok := o.NextHistory()
			if ok {
				o.buf.Set(buf)
			}
		case CharInterrupt:
			if o.IsSearchMode() {
				o.t.KickRead()
				o.ExitSearchMode(true)
				break
			}
			if o.IsInCompleteMode() {
				o.t.KickRead()
				o.ExitCompleteMode(true)
				o.buf.Refresh()
				break
			}
			o.buf.MoveToLineEnd()
			o.buf.Refresh()
			o.buf.WriteString("^C\n")
			o.outchan <- nil
		default:
			if o.IsSearchMode() {
				o.SearchChar(r)
				keepInSearchMode = true
				break
			}
			o.buf.WriteRune(r)
			if o.IsInCompleteMode() {
				o.OnComplete()
				keepInCompleteMode = true
			}
		}
		if !keepInSearchMode && o.IsSearchMode() {
			o.ExitSearchMode(false)
			o.buf.Refresh()
		} else if o.IsInCompleteMode() {
			if !keepInCompleteMode {
				o.ExitCompleteMode(false)
				o.buf.Refresh()
			} else {
				o.buf.Refresh()
				o.CompleteRefresh()
			}
		}
		if !o.IsSearchMode() {
			o.UpdateHistory(o.buf.Runes(), false)
		}
	}
}

func (o *Operation) Stderr() io.Writer {
	return &wrapWriter{target: o.cfg.Stderr, r: o, t: o.t}
}

func (o *Operation) Stdout() io.Writer {
	return &wrapWriter{target: o.cfg.Stdout, r: o, t: o.t}
}

func (o *Operation) String() (string, error) {
	r, err := o.Runes()
	if err != nil {
		return "", err
	}
	return string(r), nil
}

func (o *Operation) Runes() ([]rune, error) {
	o.buf.Refresh() // print prompt
	o.t.KickRead()
	r := <-o.outchan
	if r == nil {
		return nil, io.EOF
	}
	return r, nil
}

func (o *Operation) Slice() ([]byte, error) {
	r, err := o.Runes()
	if err != nil {
		return nil, err
	}
	return []byte(string(r)), nil
}

func (o *Operation) Close() {
	o.opHistory.Close()
}
