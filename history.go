package readline

type HisItem struct {
	Source  []rune
	Version int64
	Tmp     []rune
}

func (h *HisItem) Clean() {
	h.Source = nil
	h.Tmp = nil
}

func (o *Operation) showItem(obj interface{}) []rune {
	item := obj.(*HisItem)
	if item.Version == o.historyVer {
		return item.Tmp
	}
	return item.Source
}

func (o *Operation) PrevHistory() []rune {
	if o.current == nil {
		return nil
	}
	current := o.current.Prev()
	if current == nil {
		return nil
	}
	o.current = current
	return o.showItem(current.Value)
}

func (o *Operation) NextHistory() ([]rune, bool) {
	if o.current == nil {
		return nil, false
	}
	current := o.current.Next()
	if current == nil {
		return nil, false
	}

	o.current = current
	return o.showItem(current.Value), true
}

func (o *Operation) NewHistory(current []rune) {
	// if just use last command without modify
	// just clean lastest history
	if o.current == o.history.Back().Prev() {
		use := o.current.Value.(*HisItem)
		if equalRunes(use.Tmp, use.Source) {
			o.current = o.history.Back()
			o.current.Value.(*HisItem).Clean()
			o.historyVer++
			return
		}
	}

	if o.current != o.history.Back() {
		// move history item to current command
		use := o.current.Value.(*HisItem)
		o.current = o.history.Back()
		current = use.Tmp
	}

	o.UpdateHistory(current, true)

	// push a new one to commit current command
	o.historyVer++
	o.PushHistory(nil)
}

func (o *Operation) UpdateHistory(s []rune, commit bool) {
	if o.current == nil {
		o.PushHistory(s)
		return
	}
	r := o.current.Value.(*HisItem)
	r.Version = o.historyVer
	if commit {
		r.Source = make([]rune, len(s))
		copy(r.Source, s)
	} else {
		r.Tmp = append(r.Tmp[:0], s...)
	}
	o.current.Value = r
}

func (o *Operation) PushHistory(s []rune) {
	// copy
	newCopy := make([]rune, len(s))
	copy(newCopy, s)
	elem := o.history.PushBack(&HisItem{Source: newCopy})
	o.current = elem
}
