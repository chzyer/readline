package readline

import (
	"reflect"
	"testing"
)

type twidth struct {
	r      []rune
	length int
}

func TestRuneWidth(t *testing.T) {
	rs := []twidth{
		{[]rune("☭"), 1},
		{[]rune("a"), 1},
		{[]rune("你"), 2},
		{runes.ColorFilter([]rune("☭\033[13;1m你")), 3},
		{[]rune("，。！？；："), 12},
		{[]rune("你好，世界！"), 12},
		{[]rune("【测试】"), 8},
		{[]rune("《书名》"), 8},
		{[]rune("「引号」"), 8},
		{[]rune("a，b"), 4},
		{runes.ColorFilter([]rune("\033[32m，\033[0m")), 2},
	}
	for _, r := range rs {
		if w := runes.WidthAll(r.r); w != r.length {
			t.Fatalf("expected width %d for %q, got %d", r.length, string(r.r), w)
		}
	}
}

func TestFullwidthPunctuationWidth(t *testing.T) {
	fullwidthRunes := []rune{
		'，', // U+FF0C
		'。', // U+3002
		'！', // U+FF01
		'？', // U+FF1F
		'；', // U+FF1B
		'：', // U+FF1A
		'、', // U+3001
		'【', // U+3010
		'】', // U+3011
		'《', // U+300A
		'》', // U+300B
		'「', // U+300C
		'」', // U+300D
		'『', // U+300E
		'』', // U+300F
		'（', // U+FF08
		'）', // U+FF09
		'［', // U+FF3B
		'］', // U+FF3D
		'｛', // U+FF5B
		'｝', // U+FF5D
		'～', // U+FF5E
	}

	for _, r := range fullwidthRunes {
		if w := runes.Width(r); w != 2 {
			t.Fatalf("expected width 2 for fullwidth rune %c (%U), got %d", r, r, w)
		}
	}

	asciiPunctuation := []rune{',', '.', '!', '?', ';', ':', '[', ']', '{', '}'}
	for _, r := range asciiPunctuation {
		if w := runes.Width(r); w != 1 {
			t.Fatalf("expected width 1 for ascii rune %c (%U), got %d", r, r, w)
		}
	}

	bs := runes.Backspace([]rune("，"))
	if len(bs) != 2 {
		t.Fatalf("expected backspace length 2 for fullwidth comma, got %d", len(bs))
	}
}

type tagg struct {
	r      [][]rune
	e      [][]rune
	length int
}

func TestAggRunes(t *testing.T) {
	rs := []tagg{
		{
			[][]rune{[]rune("ab"), []rune("a"), []rune("abc")},
			[][]rune{[]rune("b"), []rune(""), []rune("bc")},
			1,
		},
		{
			[][]rune{[]rune("addb"), []rune("ajkajsdf"), []rune("aasdfkc")},
			[][]rune{[]rune("ddb"), []rune("jkajsdf"), []rune("asdfkc")},
			1,
		},
		{
			[][]rune{[]rune("ddb"), []rune("ajksdf"), []rune("aasdfkc")},
			[][]rune{[]rune("ddb"), []rune("ajksdf"), []rune("aasdfkc")},
			0,
		},
		{
			[][]rune{[]rune("ddb"), []rune("ddajksdf"), []rune("ddaasdfkc")},
			[][]rune{[]rune("b"), []rune("ajksdf"), []rune("aasdfkc")},
			2,
		},
	}
	for _, r := range rs {
		same, off := runes.Aggregate(r.r)
		if off != r.length {
			t.Fatal("result not expect", off)
		}
		if len(same) != off {
			t.Fatal("result not expect", same)
		}
		if !reflect.DeepEqual(r.r, r.e) {
			t.Fatal("result not expect")
		}
	}
}
