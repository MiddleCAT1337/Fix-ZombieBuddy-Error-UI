package lua

import (
	"fmt"
	"os"
	"strings"
)

const (
	applySettingsMarker = "local function applySettings()"
	guardLine           = "    if not ZombieBuddy then return end"
	guardSubstring      = "if not ZombieBuddy then return end"
)

// NeedsPatch reports whether the file content is missing the nil guard.
func NeedsPatch(content string) bool {
	norm := strings.ReplaceAll(content, "\r\n", "\n")
	if !strings.Contains(norm, applySettingsMarker) {
		return false
	}
	return !strings.Contains(norm, guardSubstring)
}

// PatchContent inserts the guard after applySettings() when needed.
func PatchContent(content string) (string, bool, error) {
	if !strings.Contains(content, applySettingsMarker) {
		return content, false, fmt.Errorf("applySettings() not found")
	}
	if strings.Contains(content, guardSubstring) {
		return content, false, nil
	}
	idx := strings.Index(content, applySettingsMarker)
	if idx < 0 {
		return content, false, nil
	}
	insertAt := idx + len(applySettingsMarker)
	if insertAt < len(content) && content[insertAt] == '\r' {
		insertAt++
	}
	if insertAt < len(content) && content[insertAt] == '\n' {
		insertAt++
	}
	var b strings.Builder
	b.WriteString(content[:insertAt])
	b.WriteString("\n")
	b.WriteString(guardLine)
	b.WriteString(content[insertAt:])
	return b.String(), true, nil
}

// NeedsPatchFile reads path and checks if patching is required.
func NeedsPatchFile(path string) (bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return false, err
	}
	return NeedsPatch(string(data)), nil
}

// PatchFile applies the guard idempotently.
func PatchFile(path string) (changed bool, err error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return false, err
	}
	newContent, changed, err := PatchContent(string(data))
	if err != nil {
		return false, err
	}
	if !changed {
		return false, nil
	}
	return true, os.WriteFile(path, []byte(newContent), 0644)
}
