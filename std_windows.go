// +build windows

package readline

func init() {
	Stdin = NewRawReader()
}
