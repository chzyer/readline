package readline

import "testing"

func TestRuneSliceClampsToBufferBounds(t *testing.T) {
	empty := &RuneBuffer{}
	if got := empty.RuneSlice(1); len(got) != 0 {
		t.Fatalf("forward slice on empty buffer length = %d, want 0", len(got))
	}
	if got := empty.RuneSlice(-1); len(got) != 0 {
		t.Fatalf("backward slice on empty buffer length = %d, want 0", len(got))
	}

	rb := &RuneBuffer{buf: []rune("abc"), idx: 1}
	if got := string(rb.RuneSlice(5)); got != "bc" {
		t.Fatalf("forward slice = %q, want %q", got, "bc")
	}
	if got := string(rb.RuneSlice(-5)); got != "a" {
		t.Fatalf("backward slice = %q, want %q", got, "a")
	}
}
