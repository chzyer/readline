package readline

type PrefixCompleter struct {
	Name     []rune
	Children []AutoCompleter
}

func (p PrefixCompleter) GetName() []rune {
	return p.Name
}

func (p PrefixCompleter) GetChildren() []AutoCompleter {
	return p.Children
}

func NewPrefixCompleter(pc ...AutoCompleter) AutoCompleter {
	return PcItem("", pc...)
}

func PcItem(name string, pc ...AutoCompleter) AutoCompleter {
	name += " "
	result := AutoCompleter(PrefixCompleter{
		Name:     []rune(name),
		Children: pc,
	})
	return result
}
