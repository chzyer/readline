package readline

import (
	"reflect"
	"testing"
)

func TestRuneSliceClampsBounds(t *testing.T) {
	r := &RuneBuffer{buf: []rune("abc"), idx: 1}

	// Forward slice longer than remaining buffer.
	got := r.RuneSlice(10)
	want := []rune("bc")
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("forward clamp: got %q want %q", string(got), string(want))
	}

	// Backward slice longer than available prefix (issue #241).
	got = r.RuneSlice(-10)
	want = []rune("a")
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("backward clamp: got %q want %q", string(got), string(want))
	}

	// Exact sizes still work.
	got = r.RuneSlice(1)
	if !reflect.DeepEqual(got, []rune("b")) {
		t.Fatalf("forward exact: got %q", string(got))
	}
	got = r.RuneSlice(-1)
	if !reflect.DeepEqual(got, []rune("a")) {
		t.Fatalf("backward exact: got %q", string(got))
	}

	// Empty buffer / cursor at start.
	r = &RuneBuffer{buf: []rune{}, idx: 0}
	if got = r.RuneSlice(-1); len(got) != 0 {
		t.Fatalf("empty backward: got %q", string(got))
	}
	if got = r.RuneSlice(1); len(got) != 0 {
		t.Fatalf("empty forward: got %q", string(got))
	}
}
