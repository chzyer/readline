//go:build !windows

package readline

func enableANSI() error {
	return nil
}
