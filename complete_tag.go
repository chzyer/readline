package readline

type TabCompleter struct{}

func (t *TabCompleter) Do([]rune, int) ([][]rune, int) {
	return [][]rune{[]rune("\t")}, 0
}
