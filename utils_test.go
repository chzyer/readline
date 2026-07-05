package readline

import "testing"

// TestSplitByLineWidePrompt verifies that SplitByLine accounts for a prompt
// wider than the terminal. len(result)-1 is the terminal row the last rune sits
// on (idxLine); clean() uses it to decide how many rows to move up on redraw.
// When the prompt wraps on its own the extra rows must be counted, otherwise a
// narrow-terminal prompt repeats on every keystroke.
func TestSplitByLineWidePrompt(t *testing.T) {
	tests := []struct {
		name        string
		promptWidth int
		screenWidth int
		text        string
		wantRow     int // expected idxLine == len(result)-1
	}{
		{"normal prompt empty buffer", 5, 20, "", 0},
		{"normal prompt short text", 5, 20, "abc", 0},
		{"prompt exactly screen width, empty", 20, 20, "", 1},
		{"prompt just over screen width, empty", 21, 20, "", 1},
		{"prompt three screens wide, empty", 60, 20, "", 3},
		{"prompt three screens wide, text stays on last row", 60, 20, "abc", 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SplitByLine(tt.promptWidth, tt.screenWidth, []rune(tt.text))
			row := len(got) - 1
			if row != tt.wantRow {
				t.Fatalf("SplitByLine(start=%d, width=%d, %q) row = %d, want %d (lines=%q)",
					tt.promptWidth, tt.screenWidth, tt.text, row, tt.wantRow, got)
			}
		})
	}
}

// TestSplitByLineConsistentWithLineCount ensures the row index reported by
// SplitByLine for an empty buffer never exceeds the total number of rows the
// prompt occupies per LineCount, across a range of prompt/terminal widths. If
// SplitByLine over-counted, clean() would move the cursor up past the prompt and
// erase unrelated output printed above a narrow menu. (At an exact-multiple
// boundary the two can be equal due to the pending-wrap `>=` convention, so the
// guard is row > total, not row > total-1.)
func TestSplitByLineConsistentWithLineCount(t *testing.T) {
	screenWidth := 20
	for promptWidth := 0; promptWidth <= 100; promptWidth++ {
		row := len(SplitByLine(promptWidth, screenWidth, nil)) - 1
		total := LineCount(screenWidth, promptWidth)
		if row > total {
			t.Fatalf("promptWidth=%d: SplitByLine row=%d exceeds LineCount total=%d", promptWidth, row, total)
		}
	}
}
