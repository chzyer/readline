package readline

func (l *Readline) PrevHistory() []rune {
	if l.current == nil {
		return nil
	}
	current := l.current.Prev()
	if current == nil {
		return nil
	}
	l.current = current
	return current.Value.([]rune)
}

func (l *Readline) NextHistory() []rune {
	if l.current == nil {
		return nil
	}
	current := l.current.Next()
	if current == nil {
		return nil
	}
	l.current = current
	return current.Value.([]rune)
}

func (l *Readline) NewHistory(current []rune) {
	l.UpdateHistory(current)
	if l.current != l.history.Back() {
		// move history item to current command
		l.history.Remove(l.current)
		use := l.current.Value.([]rune)
		l.current = l.history.Back()
		l.UpdateHistory(use)
	}

	// push a new one to commit current command
	l.PushHistory(nil)
}

func (l *Readline) UpdateHistory(s []rune) {
	if l.current == nil {
		l.PushHistory(s)
		return
	}
	r := l.current.Value.([]rune)
	l.current.Value = append(r[:0], s...)
}

func (l *Readline) PushHistory(s []rune) {
	// copy
	newCopy := make([]rune, len(s))
	copy(newCopy, s)
	elem := l.history.PushBack(newCopy)
	l.current = elem
}
