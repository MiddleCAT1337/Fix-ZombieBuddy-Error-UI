//go:build !windows

package steamcheck

func Running() bool {
	return false
}
