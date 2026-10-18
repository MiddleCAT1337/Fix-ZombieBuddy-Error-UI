package apply

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/whitehorizon/zb-patch/internal/launcher"
	"github.com/whitehorizon/zb-patch/internal/lua"
	"github.com/whitehorizon/zb-patch/internal/plan"
)

const backupSuffix = ".zb-patch.bak"

// Run applies all planned changes.
func Run(cs plan.ChangeSet) error {
	for _, c := range cs.Copies {
		if err := copyFile(c.Src, c.Dst); err != nil {
			return fmt.Errorf("copy %s: %w", c.Src, err)
		}
	}

	if cs.JSONPatch {
		if err := backupOnce(cs.JSONPath); err != nil {
			return err
		}
		if _, err := launcher.Patch(cs.JSONPath); err != nil {
			return fmt.Errorf("patch launcher JSON: %w", err)
		}
	}

	for _, p := range cs.LuaToPatch {
		if err := backupOnce(p); err != nil {
			return err
		}
		if _, err := lua.PatchFile(p); err != nil {
			return fmt.Errorf("patch lua %s: %w", p, err)
		}
	}

	return nil
}

func backupOnce(path string) error {
	bak := path + backupSuffix
	if _, err := os.Stat(bak); err == nil {
		return nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("backup read %s: %w", path, err)
	}
	if err := os.WriteFile(bak, data, 0644); err != nil {
		return fmt.Errorf("backup write %s: %w", bak, err)
	}
	return nil
}

func copyFile(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}
