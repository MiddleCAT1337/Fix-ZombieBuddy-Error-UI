package discover

import (
	"fmt"
	"os"
	"path/filepath"
)

// Paths holds resolved install locations.
type Paths struct {
	Steam      string
	PZ         string
	WorkshopZB string
}

// Overrides optional manual paths (empty = auto-detect).
type Overrides struct {
	Steam      string
	PZ         string
	WorkshopZB string
}

// Resolve discovers Steam, Project Zomboid, and ZombieBuddy workshop roots.
func Resolve(o Overrides) (Paths, error) {
	steam := o.Steam
	if steam == "" {
		var err error
		steam, err = detectSteamPath()
		if err != nil {
			return Paths{}, fmt.Errorf("detect Steam: %w", err)
		}
	}
	steam = filepath.Clean(steam)

	pz := o.PZ
	if pz == "" {
		var err error
		pz, err = detectPZPath(steam)
		if err != nil {
			return Paths{}, fmt.Errorf("detect Project Zomboid: %w", err)
		}
	}
	pz = filepath.Clean(pz)

	workshop := o.WorkshopZB
	if workshop == "" {
		var err error
		workshop, err = detectZBWorkshopPath(steam)
		if err != nil {
			return Paths{}, fmt.Errorf("detect ZombieBuddy workshop: %w (subscribe: https://steamcommunity.com/sharedfiles/filedetails/?id=%s)", err, ZBWorkshopID)
		}
	}
	workshop = filepath.Clean(workshop)

	return Paths{Steam: steam, PZ: pz, WorkshopZB: workshop}, nil
}

func detectPZPath(steamPath string) (string, error) {
	defaultPath := filepath.Join(steamPath, "steamapps", filepath.FromSlash(RelPZSub))
	if dirExists(defaultPath) {
		return defaultPath, nil
	}
	libraryVDF := filepath.Join(steamPath, "steamapps", "libraryfolders.vdf")
	if fileExists(libraryVDF) {
		return findAppInLibraries(libraryVDF, filepath.FromSlash(RelPZSub))
	}
	return "", fmt.Errorf("Project Zomboid installation not found")
}

func detectZBWorkshopPath(steamPath string) (string, error) {
	defaultPath := filepath.Join(steamPath, "steamapps", filepath.FromSlash(RelZBWorkshopSub))
	if dirExists(defaultPath) {
		return defaultPath, nil
	}
	libraryVDF := filepath.Join(steamPath, "steamapps", "libraryfolders.vdf")
	if fileExists(libraryVDF) {
		return findAppInLibraries(libraryVDF, filepath.FromSlash(RelZBWorkshopSub))
	}
	return "", fmt.Errorf("workshop content folder not found")
}

func findAppInLibraries(libraryVDFPath string, subPath string) (string, error) {
	m, err := parseVDFMap(libraryVDFPath)
	if err != nil {
		return "", err
	}

	libraryfolders, ok := m["libraryfolders"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid libraryfolders.vdf")
	}

	for _, folder := range libraryfolders {
		folderMap, ok := folder.(map[string]interface{})
		if !ok {
			continue
		}
		path, ok := folderMap["path"].(string)
		if !ok {
			continue
		}
		path = filepath.Clean(filepath.FromSlash(path))
		apps, ok := folderMap["apps"].(map[string]interface{})
		if !ok {
			continue
		}
		if _, found := apps[PZAppID]; !found {
			continue
		}
		fullPath := filepath.Join(path, "steamapps", filepath.FromSlash(subPath))
		if dirExists(fullPath) {
			return fullPath, nil
		}
	}
	return "", fmt.Errorf("app %s path not found in any Steam library", PZAppID)
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

// LibsDir returns the ZombieBuddy libs folder under workshop content.
func (p Paths) LibsDir() string {
	return filepath.Join(p.WorkshopZB, filepath.FromSlash(RelZBLibs))
}

// LauncherJSONPath returns ProjectZomboid64.json in the game folder.
func (p Paths) LauncherJSONPath() string {
	return filepath.Join(p.PZ, LauncherJSON)
}

// LuaOptionsPaths returns existing ZombieBuddy_Options.lua paths (41/42).
func (p Paths) LuaOptionsPaths() []string {
	var out []string
	for _, rel := range []string{LuaOptionsRel41, LuaOptionsRel42} {
		full := filepath.Join(p.WorkshopZB, filepath.FromSlash(rel))
		if fileExists(full) {
			out = append(out, full)
		}
	}
	return out
}
