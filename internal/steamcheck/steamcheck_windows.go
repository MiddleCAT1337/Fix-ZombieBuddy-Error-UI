//go:build windows

package steamcheck

import (
	"os/exec"
	"strings"
)

// Running reports whether steam.exe appears in the process list.
func Running() bool {
	cmd := exec.Command("tasklist", "/FI", "IMAGENAME eq steam.exe", "/NH")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.Contains(strings.ToLower(string(output)), "steam.exe")
}
