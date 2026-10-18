//go:build !windows

package discover

import "fmt"

func detectSteamPath() (string, error) {
	return "", fmt.Errorf("Steam registry lookup is only supported on Windows")
}
