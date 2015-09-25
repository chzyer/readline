package readline

import "testing"

type Twidth struct {
	r      []rune
	length int
}

func TestRuneWidth(t *testing.T) {
	runes := []Twidth{
		{[]rune("☭"), 1},
		{[]rune("a"), 1},
		{[]rune("你"), 2},
		{RunesColorFilter([]rune("☭\033[13;1m你")), 3},
	}
	for _, r := range runes {
		if w := RunesWidth(r.r); w != r.length {
			t.Fatal("result not expect", r.r, r.length, w)
		}
	}
}
