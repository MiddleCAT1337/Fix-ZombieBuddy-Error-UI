//go:build !windows

package elevation

func EnsureAdmin() error {
	return nil
}
