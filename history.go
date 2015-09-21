package readline

func (o *Operation) PrevHistory() []rune {
	if o.current == nil {
		return nil
	}
	current := o.current.Prev()
	if current == nil {
		return nil
	}
	o.current = current
	return current.Value.([]rune)
}

func (o *Operation) NextHistory() []rune {
	if o.current == nil {
		return nil
	}
	current := o.current.Next()
	if current == nil {
		return nil
	}
	o.current = current
	return current.Value.([]rune)
}

func (o *Operation) NewHistory(current []rune) {
	o.UpdateHistory(current)
	if o.current != o.history.Back() {
		// move history item to current command
		o.history.Remove(o.current)
		use := o.current.Value.([]rune)
		o.current = o.history.Back()
		o.UpdateHistory(use)
	}

	// push a new one to commit current command
	o.PushHistory(nil)
}

func (o *Operation) UpdateHistory(s []rune) {
	if o.current == nil {
		o.PushHistory(s)
		return
	}
	r := o.current.Value.([]rune)
	o.current.Value = append(r[:0], s...)
}

func (o *Operation) PushHistory(s []rune) {
	// copy
	newCopy := make([]rune, len(s))
	copy(newCopy, s)
	elem := o.history.PushBack(newCopy)
	o.current = elem
}
