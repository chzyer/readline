package readline

import (
	"reflect"
	"testing"
)

func TestSplitByMultiLine(t *testing.T) {
	rs := []rune("hello!bye!!!!")
	expected := []string{"hell", "o!", "bye!", "!!", "!"}
	ret := SplitByMultiLine(0, 6, 4, rs)
	if !reflect.DeepEqual(ret, expected) {
		t.Fatal(ret, expected)
	}
}
