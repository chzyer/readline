// +build windows
// visual test
package readline

import (
	"fmt"
	"testing"
)

func TestAnsiWriter(t *testing.T) {
	rl, err := New("")
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		fmt.Fprintf(rl, "\033[3%[1]dmcolor\033[01;3%[1]dmcolor\033[4%[1]dmcolor\033[mt\033[0mreset\n", i)
	}
}
