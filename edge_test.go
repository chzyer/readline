package readline

import "testing"

func TestAggregateEmpty(t *testing.T) {
	same, size := runes.Aggregate(nil)
	if same != nil || size != 0 {
		t.Fatalf("nil: same=%v size=%d", same, size)
	}
	same, size = runes.Aggregate([][]rune{})
	if same != nil || size != 0 {
		t.Fatalf("empty: same=%v size=%d", same, size)
	}
	same, size = runes.Aggregate([][]rune{[]rune("abc"), []rune("abd")})
	if string(same) != "ab" || size != 2 {
		t.Fatalf("common: same=%q size=%d", string(same), size)
	}
}

func TestLineCountNonPositiveWidth(t *testing.T) {
	if got := LineCount(0, 10); got != 1 {
		t.Fatalf("LineCount(0,10)=%d want 1", got)
	}
	if got := LineCount(-1, 10); got != 1 {
		t.Fatalf("LineCount(-1,10)=%d want 1", got)
	}
	if got := LineCount(0, 0); got != 0 {
		t.Fatalf("LineCount(0,0)=%d want 0", got)
	}
	if got := LineCount(80, 100); got != 2 {
		t.Fatalf("LineCount(80,100)=%d want 2", got)
	}
}
