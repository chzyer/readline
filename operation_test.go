package readline

import (
	"io"
	"os"
	"testing"
)

func TestReadlineReturnsAfterListener(t *testing.T) {
	for _, uniqueEditLine := range []bool{false, true} {
		for _, tt := range []struct {
			name  string
			input string
			key   rune
			line  string
			err   error
		}{
			{"enter", "yes\r", CharEnter, "yes", nil},
			{"empty-enter", "\r", CharEnter, "", nil},
			{"ctrl-j", "yes\n", CharCtrlJ, "yes", nil},
			{"ctrl-c", "yes\x03", CharInterrupt, "yes", ErrInterrupt},
			{"ctrl-d", "\x04", CharDelete, "", io.EOF},
		} {
			name := tt.name
			if uniqueEditLine {
				name += "/unique-edit-line"
			}
			t.Run(name, func(t *testing.T) {
				input, writer := io.Pipe()
				defer input.Close()
				defer writer.Close()
				output, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
				if err != nil {
					t.Fatal(err)
				}
				defer output.Close()

				listenerDone := make(chan struct{})
				listenerCalled := false
				cfg := &Config{
					Stdin:              input,
					Stdout:             output,
					Stderr:             output,
					UniqueEditLine:     uniqueEditLine,
					FuncGetWidth:       func() int { return 80 },
					FuncIsTerminal:     func() bool { return false },
					FuncMakeRaw:        func() error { return nil },
					FuncExitRaw:        func() error { return nil },
					FuncOnWidthChanged: func(func()) {},
				}
				cfg.SetListener(func(line []rune, pos int, key rune) ([]rune, int, bool) {
					if key == tt.key {
						listenerCalled = true
						close(listenerDone)
					}
					return nil, 0, false
				})
				rl, err := NewEx(cfg)
				if err != nil {
					t.Fatal(err)
				}
				defer rl.Close()

				written := make(chan error, 1)
				go func() {
					_, err := io.WriteString(writer, tt.input)
					written <- err
				}()
				line, err := rl.Readline()
				// A caller must be able to read listener-owned state after Readline returns.
				calledAtReturn := listenerCalled
				// Join the callback after that read, including when testing the broken implementation.
				<-listenerDone
				if !calledAtReturn {
					t.Error("Readline returned before the listener completed")
				}
				if line != tt.line || err != tt.err {
					t.Errorf("Readline() = (%q, %v), want (%q, %v)", line, err, tt.line, tt.err)
				}
				if err := <-written; err != nil {
					t.Fatal(err)
				}
			})
		}
	}
}
