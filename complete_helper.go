package readline

import "github.com/chzyer/readline/runes"

type PrefixCompleter struct {
	Name     []rune
	Children []*PrefixCompleter
}

func NewPrefixCompleter(pc ...*PrefixCompleter) *PrefixCompleter {
	return PcItem("", pc...)
}

func PcItem(name string, pc ...*PrefixCompleter) *PrefixCompleter {
	name += " "
	return &PrefixCompleter{
		Name:     []rune(name),
		Children: pc,
	}
}

func (p *PrefixCompleter) Do(line []rune, pos int) (newLine [][]rune, offset int) {
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
