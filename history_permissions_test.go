package readline

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestNewHistoryFilesArePrivate(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Windows does not enforce Unix file permission bits")
	}

	path := filepath.Join(t.TempDir(), "history")
	history := newOpHistory(&Config{HistoryFile: path, HistoryLimit: 10})
	history.historyUpdatePath(path)
	defer history.Close()

	checkMode := func(phase string) {
		t.Helper()
		info, err := os.Stat(path)
		if err != nil {
			t.Fatal(err)
		}
		if got := info.Mode().Perm(); got&0077 != 0 {
			t.Errorf("%s: history permissions %04o expose entries to group/others", phase, got)
		}
	}

	checkMode("initial create")
	history.Push([]rune("private command"))
	history.Rewrite()
	checkMode("rewrite")
}
