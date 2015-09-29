// +build windows

package readline

import (
	"reflect"
	"syscall"
	"unsafe"
)

var (
	kernel = NewKernel()
	stdout = uintptr(syscall.Stdout)
)

type Kernel struct {
	SetConsoleCursorPosition,
	SetConsoleTextAttribute,
	GetConsoleScreenBufferInfo,
	FillConsoleOutputCharacterW,
	GetStdHandle CallFunc
}

type short int16
type word uint16

type _COORD struct {
	x short
	y short
}

func (c *_COORD) ptr() uintptr {
	return uintptr(*(*int32)(unsafe.Pointer(c)))
}

type _CONSOLE_SCREEN_BUFFER_INFO struct {
	dwSize              _COORD
	dwCursorPosition    _COORD
	wAttributes         word
	srWindow            _SMALL_RECT
	dwMaximumWindowSize _COORD
}

type _SMALL_RECT struct {
	left   short
	top    short
	right  short
	bottom short
}

type CallFunc func(u ...uintptr) error

func NewKernel() *Kernel {
	k := &Kernel{}
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	v := reflect.ValueOf(k).Elem()
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		name := t.Field(i).Name
		f := kernel32.NewProc(name)
		v.Field(i).Set(reflect.ValueOf(k.Wrap(f)))
	}
	return k
}

func (k *Kernel) Wrap(p *syscall.LazyProc) CallFunc {
	return func(args ...uintptr) error {
		var r0 uintptr
		var e1 syscall.Errno
		size := uintptr(len(args))
		if len(args) <= 3 {
			buf := make([]uintptr, 3)
			copy(buf, args)
			r0, _, e1 = syscall.Syscall(p.Addr(), size,
				buf[0], buf[1], buf[2])
		} else {
			buf := make([]uintptr, 6)
			copy(buf, args)
			r0, _, e1 = syscall.Syscall6(p.Addr(), size,
				buf[0], buf[1], buf[2], buf[3], buf[4], buf[5],
			)
		}

		if int(r0) == 0 {
			if e1 != 0 {
				return error(e1)
			} else {
				return syscall.EINVAL
			}
		}
		return nil
	}

}

func GetConsoleScreenBufferInfo() (*_CONSOLE_SCREEN_BUFFER_INFO, error) {
	t := new(_CONSOLE_SCREEN_BUFFER_INFO)
	err := kernel.GetConsoleScreenBufferInfo(
		stdout,
		uintptr(unsafe.Pointer(t)),
	)
	return t, err
}
