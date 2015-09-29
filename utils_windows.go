// +build windows

package readline

// get width of the terminal
func getWidth() int {
	info, _ := GetConsoleScreenBufferInfo()
	if info == nil {
		return 0
	}
	return int(info.dwSize.x)
}
