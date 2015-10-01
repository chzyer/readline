package readline

const (
	VIM_NORMAL = iota
	VIM_INSERT
	VIM_VISUAL
)

type opVim struct {
	cfg     *Config
	op      *Operation
	vimMode int
}

func newVimMode(op *Operation) *opVim {
	ov := &opVim{
		cfg: op.cfg,
		op:  op,
	}
	ov.SetVimMode(ov.cfg.VimMode)
	return ov
}

func (o *opVim) SetVimMode(on bool) {
	if o.cfg.VimMode && !on { // turn off
		o.ExitVimMode()
	}
	o.cfg.VimMode = on
	o.vimMode = VIM_INSERT
}

func (o *opVim) ExitVimMode() {
	o.vimMode = VIM_NORMAL
}

func (o *opVim) IsEnableVimMode() bool {
	return o.cfg.VimMode
}

func (o *opVim) HandleVimNormal(r rune, readNext func() rune) (t rune, handle bool) {
	switch r {
	case CharEnter, CharInterrupt:
		return r, false
	}
	rb := o.op.buf
	handled := true
	{
		switch r {
		case 'h':
			t = CharBackward
		case 'j':
			t = CharNext
		case 'k':
			t = CharPrev
		case 'l':
			t = CharForward
		default:
			handled = false
		}
		if handled {
			return t, false
		}
	}

	{ // to insert
		handled = true
		switch r {
		case 'i':
		case 'I':
			rb.MoveToLineStart()
		case 'a':
			rb.MoveForward()
		case 'A':
			rb.MoveToLineEnd()
		case 's':
			rb.Delete()
		default:
			handled = false
		}
		if handled {
			o.EnterVimInsertMode()
			return r, true
		}
	}

	{ // movement
		handled = true
		switch r {
		case '0', '^':
			rb.MoveToLineStart()
		case '$':
			rb.MoveToLineEnd()
		case 'b':
			rb.MoveToPrevWord()
		case 'w':
			rb.MoveToNextWord()
		case 'f', 'F', 't', 'T':
			next := readNext()
			prevChar := r == 't' || r == 'T'
			reverse := r == 'F' || r == 'T'
			switch next {
			case CharEsc:
			default:
				if rb.MoveTo(next, prevChar, reverse) {
					return r, true
				}
			}
		default:
			handled = false
		}
		if handled {
			return r, true
		}
	}

	// invalid operation
	o.op.t.Bell()
	return r, true
}

func (o *opVim) EnterVimInsertMode() {
	o.vimMode = VIM_INSERT
}

func (o *opVim) ExitVimInsertMode() {
	o.vimMode = VIM_NORMAL
}

func (o *opVim) HandleVim(r rune, readNext func() rune) (rune, bool) {
	if o.vimMode == VIM_NORMAL {
		return o.HandleVimNormal(r, readNext)
	}
	if r == CharEsc {
		o.ExitVimInsertMode()
		return r, true
	}

	switch o.vimMode {
	case VIM_INSERT:
		return r, false
	case VIM_VISUAL:
	}
	return r, false
}
