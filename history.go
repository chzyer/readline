package readline

import (
	"bufio"
	"container/list"
	"os"
	"strings"
)

type HisItem struct {
	Source  []rune
	Version int64
	Tmp     []rune
}

func (h *HisItem) Clean() {
	h.Source = nil
	h.Tmp = nil
}

type opHistory struct {
	path       string
	history    *list.List
	historyVer int64
	current    *list.Element
	fd         *os.File
}

func newOpHistory(path string) (o *opHistory) {
	o = &opHistory{
		path:    path,
		history: list.New(),
	}
	if o.path == "" {
		return
	}
	f, err := os.OpenFile(o.path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return
	}
	o.fd = f
	r := bufio.NewReader(o.fd)
	for {
		line, err := r.ReadSlice('\n')
		if err != nil {
			break
		}
		o.PushHistory([]rune(strings.TrimSpace(string(line))))
	}
	o.historyVer++
	o.PushHistory(nil)
	return
}

func (o *opHistory) Close() {
	if o.fd != nil {
		o.fd.Close()
	}
}

func (o *opHistory) FindHistoryBck(rs []rune) (int, *list.Element) {
	for elem := o.current; elem != nil; elem = elem.Prev() {
		idx := RunesIndex(o.showItem(elem.Value), rs)
		if idx < 0 {
			continue
		}
		return idx, elem
	}
	return -1, nil
}

func (o *opHistory) FindHistoryFwd(rs []rune) (int, *list.Element) {
	for elem := o.current; elem != nil; elem = elem.Next() {
		idx := RunesIndex(o.showItem(elem.Value), rs)
		if idx < 0 {
			continue
		}
		return idx, elem
	}
	return -1, nil
}

func (o *opHistory) showItem(obj interface{}) []rune {
	item := obj.(*HisItem)
	if item.Version == o.historyVer {
		return item.Tmp
	}
	return item.Source
}

func (o *opHistory) PrevHistory() []rune {
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

func (o *opHistory) NextHistory() ([]rune, bool) {
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

func (o *opHistory) NewHistory(current []rune) {
	// if just use last command without modify
	// just clean lastest history
	if back := o.history.Back(); back != nil {
		prev := back.Prev()
		if prev != nil {
			use := o.current.Value.(*HisItem)
			if equalRunes(use.Tmp, prev.Value.(*HisItem).Source) {
				o.current = o.history.Back()
				o.current.Value.(*HisItem).Clean()
				o.historyVer++
				return
			}
		}
	}
	if len(current) == 0 {
		o.current = o.history.Back()
		if o.current != nil {
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

func (o *opHistory) UpdateHistory(s []rune, commit bool) {
	if o.current == nil {
		o.PushHistory(s)
		return
	}
	r := o.current.Value.(*HisItem)
	r.Version = o.historyVer
	if commit {
		r.Source = make([]rune, len(s))
		copy(r.Source, s)
		if o.fd != nil {
			o.fd.Write([]byte(string(r.Source) + "\n"))
		}
	} else {
		r.Tmp = append(r.Tmp[:0], s...)
	}
	o.current.Value = r
}

func (o *opHistory) PushHistory(s []rune) {
	// copy
	newCopy := make([]rune, len(s))
	copy(newCopy, s)
	elem := o.history.PushBack(&HisItem{Source: newCopy})
	o.current = elem
}
