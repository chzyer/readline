// +build windows

package readline

func init() {
	// don't use RawReader if stdin is a pipe (see https://github.com/chzyer/readline/issues/229)
	stdinType, _ := kernel.GetFileType(stdin)
	if stdinType == FILE_TYPE_CHAR {
		Stdin = NewRawReader()
	}
	Stdout = NewANSIWriter(Stdout)
	Stderr = NewANSIWriter(Stderr)
}
