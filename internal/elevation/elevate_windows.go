//go:build windows

package elevation

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

const swShowNormal = 1

// EnsureAdmin prompts UAC and re-launches this executable elevated when needed.
func EnsureAdmin() error {
	if isElevated() {
		return nil
	}
	return relaunchElevated()
}

func isElevated() bool {
	var token windows.Token
	err := windows.OpenProcessToken(windows.CurrentProcess(), windows.TOKEN_QUERY, &token)
	if err != nil {
		return false
	}
	defer token.Close()
	return token.IsElevated()
}

func relaunchElevated() error {
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("resolve executable: %w", err)
	}

	verb, err := windows.UTF16PtrFromString("runas")
	if err != nil {
		return err
	}
	exePath, err := windows.UTF16PtrFromString(exe)
	if err != nil {
		return err
	}

	cwd, err := os.Getwd()
	if err != nil {
		cwd = ""
	}
	cwdPtr, err := windows.UTF16PtrFromString(cwd)
	if err != nil {
		return err
	}

	shell32 := windows.NewLazySystemDLL("shell32.dll")
	shellExecuteW := shell32.NewProc("ShellExecuteW")

	ret, _, errno := shellExecuteW.Call(
		0,
		uintptr(unsafe.Pointer(verb)),
		uintptr(unsafe.Pointer(exePath)),
		0,
		uintptr(unsafe.Pointer(cwdPtr)),
		swShowNormal,
	)
	if ret <= 32 {
		if errno == syscall.Errno(1223) { // ERROR_CANCELLED
			return fmt.Errorf("administrator permission was cancelled")
		}
		return fmt.Errorf("could not restart as administrator (code %d)", ret)
	}

	os.Exit(0)
	return nil
}
