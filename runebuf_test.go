package readline

import (
	"bytes"
	"strings"
	"testing"
)

func newTestRuneBuffer(t *testing.T, prompt string, width int) *RuneBuffer {
	cfg := &Config{
		Painter: &defaultPainter{},
	}
	if err := cfg.Init(); err != nil {
		t.Fatal(err)
	}
	return NewRuneBuffer(&bytes.Buffer{}, prompt, cfg, width)
}

type placeholderPainter struct {
	placeholder []rune
}

func (p *placeholderPainter) Paint(line []rune, pos int) []rune {
	if len(line) == 0 {
		return p.placeholder
	}
	return line
}

func TestIsInLineEdgeEmpty(t *testing.T) {
	rb := newTestRuneBuffer(t, "> ", 80)

	if rb.isInLineEdge() {
		t.Errorf("isInLineEdge() should be false when buffer is empty")
	}
}

func TestOutputEmptyBuffer(t *testing.T) {
	rb := newTestRuneBuffer(t, "> ", 80)
	out := string(rb.output())
	if out != "> " {
		t.Errorf("output() for empty buffer expected %q, got %q", "> ", out)
	}
}

func TestOutputWithPlaceholderPainter(t *testing.T) {
	cfg := &Config{
		Painter: &placeholderPainter{placeholder: []rune("hint")},
	}
	if err := cfg.Init(); err != nil {
		t.Fatal(err)
	}
	rb := NewRuneBuffer(&bytes.Buffer{}, "> ", cfg, 80)
	out := string(rb.output())
	expected := "> hint"
	if out != expected {
		t.Errorf("output() with placeholder painter expected %q, got %q", expected, out)
	}
}

func TestIsInLineEdgeBoundaries(t *testing.T) {
	// Prompt length 2: "> "
	// Width 10:
	// Line 1 fills at 10 columns (prompt 2 + 8 characters = 10)
	rb := newTestRuneBuffer(t, "> ", 10)

	// 7 characters: 2 + 7 = 9 < 10 (not at edge)
	rb.buf = []rune("1234567")
	rb.idx = len(rb.buf)
	if rb.isInLineEdge() {
		t.Errorf("isInLineEdge() should be false for width 9/10")
	}

	// 8 characters: 2 + 8 = 10 == 10 (at edge)
	rb.buf = []rune("12345678")
	rb.idx = len(rb.buf)
	if isWindows {
		if rb.isInLineEdge() {
			t.Errorf("isInLineEdge() should be false on Windows")
		}
	} else {
		if !rb.isInLineEdge() {
			t.Errorf("isInLineEdge() should be true when line exactly fills width 10")
		}
		out := string(rb.output())
		if !strings.HasSuffix(out, " \b") {
			t.Errorf("output() should end with ' \\b' when in line edge, got %q", out)
		}
	}

	// 9 characters: 8 on line 1, 1 on line 2 (not at edge)
	rb.buf = []rune("123456789")
	rb.idx = len(rb.buf)
	if rb.isInLineEdge() {
		t.Errorf("isInLineEdge() should be false after wrapping to new line")
	}

	// 18 characters: 8 on line 1, 10 on line 2 (line 2 at edge)
	rb.buf = []rune("123456789012345678")
	rb.idx = len(rb.buf)
	if isWindows {
		if rb.isInLineEdge() {
			t.Errorf("isInLineEdge() should be false on Windows")
		}
	} else {
		if !rb.isInLineEdge() {
			t.Errorf("isInLineEdge() should be true when line 2 fills width 10")
		}
	}
}

func TestIsInLineEdgeZeroWidth(t *testing.T) {
	rb := newTestRuneBuffer(t, "> ", 0)
	rb.buf = []rune("hello")
	rb.idx = len(rb.buf)
	if rb.isInLineEdge() {
		t.Errorf("isInLineEdge() should be false when width is 0")
	}
}
