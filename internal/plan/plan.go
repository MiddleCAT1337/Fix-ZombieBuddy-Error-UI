package plan

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/whitehorizon/zb-patch/internal/discover"
	"github.com/whitehorizon/zb-patch/internal/launcher"
	"github.com/whitehorizon/zb-patch/internal/lua"
)

// CopyAction is one jar/dll copy from workshop to game dir.
type CopyAction struct {
	Src string
	Dst string
}

// ChangeSet describes planned mutations.
type ChangeSet struct {
	Paths       discover.Paths
	Copies      []CopyAction
	JSONPath    string
	JSONPatch   bool
	LuaPaths    []string
	LuaToPatch  []string
	PreviewLines []string
}

// Build validates sources and computes what would change.
func Build(paths discover.Paths) (ChangeSet, error) {
	cs := ChangeSet{Paths: paths}

	libs := paths.LibsDir()
	for _, name := range []string{discover.ZBLibNative, discover.ZBLibJar} {
		src := filepath.Join(libs, name)
		if _, err := os.Stat(src); err != nil {
			return ChangeSet{}, fmt.Errorf("missing %s in %s (update workshop mod ZombieBuddy)", name, libs)
		}
		dst := filepath.Join(paths.PZ, name)
		cs.Copies = append(cs.Copies, CopyAction{Src: src, Dst: dst})
		cs.PreviewLines = append(cs.PreviewLines, fmt.Sprintf("copy %q -> %q", src, dst))
	}

	cs.JSONPath = paths.LauncherJSONPath()
	needJSON, err := launcher.NeedsPatch(cs.JSONPath)
	if err != nil {
		return ChangeSet{}, fmt.Errorf("launcher JSON: %w", err)
	}
	cs.JSONPatch = needJSON
	if needJSON {
		cs.PreviewLines = append(cs.PreviewLines, fmt.Sprintf("prepend %q to vmArgs in %q", launcher.AgentArg, cs.JSONPath))
	} else {
		cs.PreviewLines = append(cs.PreviewLines, fmt.Sprintf("skip %q (agent already present)", cs.JSONPath))
	}

	cs.LuaPaths = paths.LuaOptionsPaths()
	for _, p := range cs.LuaPaths {
		need, err := lua.NeedsPatchFile(p)
		if err != nil {
			return ChangeSet{}, fmt.Errorf("lua options %s: %w", p, err)
		}
		if need {
			cs.LuaToPatch = append(cs.LuaToPatch, p)
			cs.PreviewLines = append(cs.PreviewLines, fmt.Sprintf("insert ZombieBuddy nil guard in %q", p))
		} else {
			cs.PreviewLines = append(cs.PreviewLines, fmt.Sprintf("skip %q (guard already present or no applySettings)", p))
		}
	}
	if len(cs.LuaPaths) == 0 {
		cs.PreviewLines = append(cs.PreviewLines, "no ZombieBuddy_Options.lua found under workshop (41/42)")
	}

	return cs, nil
}
