// Package config persists the project registry at ~/.leftmark/config.json:
// one entry per registered project root, keyed by its absolute path.
package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// ProjectConfig holds the settings registered for one project root.
type ProjectConfig struct {
	Ignore []string `json:"ignore"`
}

// Registry is the on-disk registry of known projects.
type Registry struct {
	Projects map[string]ProjectConfig `json:"projects"`
}

// DefaultPath returns the registry file's default location under the
// user's home directory.
func DefaultPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".leftmark", "config.json"), nil
}

// Load reads the registry at path. A missing file is not an error; it
// returns an empty Registry, so first-run behaves like "no projects
// registered yet."
func Load(path string) (Registry, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return Registry{Projects: map[string]ProjectConfig{}}, nil
	}
	if err != nil {
		return Registry{}, err
	}

	var reg Registry
	if err := json.Unmarshal(data, &reg); err != nil {
		return Registry{}, err
	}
	if reg.Projects == nil {
		reg.Projects = map[string]ProjectConfig{}
	}
	return reg, nil
}

// Save atomically replaces the file at path with reg's contents: it writes
// to a temp file in the same directory, then renames over path, so a crash
// or concurrent read never observes a half-written file.
func Save(path string, reg Registry) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(reg, "", "  ")
	if err != nil {
		return err
	}

	tmp, err := os.CreateTemp(dir, ".config-*.json.tmp")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	defer func() { _ = os.Remove(tmpPath) }()

	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}

	return os.Rename(tmpPath, path)
}

// Lookup finds the longest registered project root that is a prefix of (or
// equal to) cwd, using pure path-string comparison against reg's keys - no
// filesystem access. found is false if no registered root matches.
func Lookup(reg Registry, cwd string) (root string, cfg ProjectConfig, found bool) {
	cwd = filepath.Clean(cwd)

	for candidate, projectCfg := range reg.Projects {
		key := filepath.Clean(candidate)
		if key != cwd && !strings.HasPrefix(cwd, key+string(filepath.Separator)) {
			continue
		}
		if !found || len(key) > len(root) {
			root, cfg, found = key, projectCfg, true
		}
	}

	return root, cfg, found
}
