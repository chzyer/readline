// +build windows

package readline

import (
	"syscall"
	"unsafe"
)

type (
	short int16
	word  uint16

	small_rect struct {
		left   short
		top    short
		right  short
		bottom short
	}

	coord struct {
		x short
		y short
	}

	console_screen_buffer_info struct {
		size                coord
		cursor_position     coord
		attributes          word
		window              small_rect
		maximum_window_size coord
	}
)

var (
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	tmp_info console_screen_buffer_info

	proc_get_console_screen_buffer_info = kernel32.NewProc("GetConsoleScreenBufferInfo")
)

func get_console_screen_buffer_info(h syscall.Handle, info *console_screen_buffer_info) (err error) {
	r0, _, e1 := syscall.Syscall(proc_get_console_screen_buffer_info.Addr(),
		2, uintptr(h), uintptr(unsafe.Pointer(info)), 0)
	if int(r0) == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func get_term_size(out syscall.Handle) coord {
	err := get_console_screen_buffer_info(out, &tmp_info)
	if err != nil {
		panic(err)
	}
	return tmp_info.size
}

// get width of the terminal
func getWidth() int {
	return int(get_term_size(syscall.Stdout).x)
}
