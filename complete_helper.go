package readline

import (
	"bytes"
	"strings"

	"github.com/chzyer/readline/runes"
)

type DynChildrenFunc func() [][]rune

type PrefixCompleter struct {
	Name        []rune
	Children    []*PrefixCompleter
	DynChildren DynChildrenFunc
}

func empty() [][]rune {
	return [][]rune{}
}

func (p *PrefixCompleter) Tree(prefix string) string {
	buf := bytes.NewBuffer(nil)
	p.Print(prefix, 0, buf)
	return buf.String()
}

func (p *PrefixCompleter) Print(prefix string, level int, buf *bytes.Buffer) {
	if strings.TrimSpace(string(p.Name)) != "" {
		buf.WriteString(prefix)
		if level > 0 {
			buf.WriteString("├")
			buf.WriteString(strings.Repeat("─", (level*4)-2))
			buf.WriteString(" ")
		}
		buf.WriteString(string(p.Name) + "\n")
		level++
	}
	for _, ch := range p.Children {
		ch.Print(prefix, level, buf)
	}
}

func NewPrefixCompleter(pc ...*PrefixCompleter) *PrefixCompleter {
	return PcItem("", pc...)
}

func PcDyn(name string, term DynChildrenFunc) *PrefixCompleter {
	name += " "
	return &PrefixCompleter{
		Name:        []rune(name),
		DynChildren: term,
	}
}

func PcItem(name string, pc ...*PrefixCompleter) *PrefixCompleter {
	name += " "
	return &PrefixCompleter{
		Name:        []rune(name),
		Children:    pc,
		DynChildren: empty,
	}
}

func (p *PrefixCompleter) Do(line []rune, pos int) (newLine [][]rune, offset int) {
	end := PcItem("")
	line = line[:pos]
	goNext := false
	var lineCompleter *PrefixCompleter
	for _, child := range p.Children {
		if len(line) >= len(child.Name) {
			if runes.HasPrefix(line, child.Name) {
				if len(line) == len(child.Name) {
					newLine = append(newLine, []rune{' '})
				} else {
					newLine = append(newLine, child.Name)
				}
				offset = len(child.Name)
				lineCompleter = child
				goNext = true
			}
		} else {
			if runes.HasPrefix(child.Name, line) {
				newLine = append(newLine, child.Name[len(line):])
				offset = len(line)
				lineCompleter = child
			}
		}
	}
	for _, dynChild := range p.DynChildren() {
		if len(line) >= len(dynChild) {
			if runes.HasPrefix(line, dynChild) {
				if len(line) == len(dynChild) {
					newLine = append(newLine, []rune{' '})
				} else {
					newLine = append(newLine, dynChild)
				}
				offset = len(dynChild)
				lineCompleter = end
			}
		} else {
			if runes.HasPrefix(dynChild, line) {
				newLine = append(newLine, dynChild[len(line):])
				offset = len(line)
				lineCompleter = end
			}
		}
	}

	if len(newLine) != 1 {
		return
	}

	tmpLine := make([]rune, 0, len(line))
	for i := offset; i < len(line); i++ {
		if line[i] == ' ' {
			continue
		}

		tmpLine = append(tmpLine, line[i:]...)
		return lineCompleter.Do(tmpLine, len(tmpLine))
	}

	if goNext {
		return lineCompleter.Do(nil, 0)
	}
	return
}
