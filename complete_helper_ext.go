package readline

import (
	"bytes"
	"strings"

	"github.com/chzyer/readline/runes"
)

type PrefixCompleterInterface interface {
	Print(prefix string, level int, buf *bytes.Buffer)
	Do(line []rune, pos int) (newLine [][]rune, length int)
	GetName() []rune
	GetChildren() []PrefixCompleterInterface
}

type PrefixCompleterExt struct {
	Name     []rune
	Children []PrefixCompleterInterface
}

func Print(p PrefixCompleterInterface, prefix string, level int, buf *bytes.Buffer) {
	if strings.TrimSpace(string(p.GetName())) != "" {
		buf.WriteString(prefix)
		if level > 0 {
			buf.WriteString("├")
			buf.WriteString(strings.Repeat("─", (level*4)-2))
			buf.WriteString(" ")
		}
		buf.WriteString(string(p.GetName()) + "\n")
		level++
	}
	for _, ch := range p.GetChildren() {
		ch.Print(prefix, level, buf)
	}
}

func (p *PrefixCompleterExt) Print(prefix string, level int, buf *bytes.Buffer) {
	Print(p, prefix, level, buf)
}

func (p *PrefixCompleterExt) GetName() []rune {
	return p.Name
}

func (p *PrefixCompleterExt) GetChildren() []PrefixCompleterInterface {
	return p.Children
}

func NewPrefixCompleterExt(pc ...PrefixCompleterInterface) *PrefixCompleterExt {
	return PceItem("", pc...)
}

func PceItem(name string, pc ...PrefixCompleterInterface) *PrefixCompleterExt {
	name += " "
	return &PrefixCompleterExt{
		Name:     []rune(name),
		Children: pc,
	}
}

func (p *PrefixCompleterExt) Do(line []rune, pos int) (newLine [][]rune, offset int) {
	return Do(p, line, pos)
}

func Do(p PrefixCompleterInterface, line []rune, pos int) (newLine [][]rune, offset int) {
	line = line[:pos]
	goNext := false
	var lineCompleter PrefixCompleterInterface
	for _, child := range p.GetChildren() {
		childName := child.GetName()
		if len(line) >= len(childName) {
			if runes.HasPrefix(line, childName) {
				if len(line) == len(childName) {
					newLine = append(newLine, []rune{' '})
				} else {
					newLine = append(newLine, childName)
				}
				offset = len(childName)
				lineCompleter = child
				goNext = true
			}
		} else {
			if runes.HasPrefix(childName, line) {
				newLine = append(newLine, childName[len(line):])
				offset = len(line)
				lineCompleter = child
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
