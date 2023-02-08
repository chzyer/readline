//go:build !windows

package readline

func enableANSI() bool {
	return true
}
