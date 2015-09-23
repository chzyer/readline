package readline

import (
	"bufio"
	"container/list"
	"os"
	"strings"
)

type hisItem struct {
	Source  []rune
	Version int64
	Tmp     []rune
}

func (h *hisItem) Clean() {
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

func (o *opHistory) FindHistoryBck(isNewSearch bool, rs []rune, start int) (int, *list.Element) {
	for elem := o.current; elem != nil; elem = elem.Prev() {
		item := o.showItem(elem.Value)
		if isNewSearch {
			start += len(rs)
		}
		if elem == o.current {
			if len(item) >= start {
				item = item[:start]
			}
		}
		idx := RunesIndexBck(item, rs)
		if idx < 0 {
			continue
		}
		return idx, elem
	}
	return -1, nil
}

func (o *opHistory) FindHistoryFwd(isNewSearch bool, rs []rune, start int) (int, *list.Element) {
	for elem := o.current; elem != nil; elem = elem.Next() {
		item := o.showItem(elem.Value)
		if isNewSearch {
			start -= len(rs)
			if start < 0 {
				start = 0
			}
		}
		if elem == o.current {
			if len(item)-1 >= start {
				item = item[start:]
			}
		}
		idx := RunesIndex(item, rs)
		if idx < 0 {
			continue
		}
		if elem == o.current {
			idx += start
		}
		return idx, elem
	}
	return -1, nil
}

func (o *opHistory) showItem(obj interface{}) []rune {
	item := obj.(*hisItem)
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
			use := o.showItem(o.current.Value.(*hisItem))
			if equalRunes(use, prev.Value.(*hisItem).Source) {
				o.current = o.history.Back()
				o.current.Value.(*hisItem).Clean()
				o.historyVer++
				return
			}
		}
	}
	if len(current) == 0 {
		o.current = o.history.Back()
		if o.current != nil {
			o.current.Value.(*hisItem).Clean()
			o.historyVer++
			return
		}
	}

	if o.current != o.history.Back() {
		// move history item to current command
		use := o.current.Value.(*hisItem)
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
	r := o.current.Value.(*hisItem)
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
	elem := o.history.PushBack(&hisItem{Source: newCopy})
	o.current = elem
}
