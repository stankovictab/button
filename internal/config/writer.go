package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// appNameToFilename converts an app display name to a YAML filename.
// e.g. "My Cool App" → "my-cool-app.yaml"
func appNameToFilename(name string) string {
	lower := strings.ToLower(strings.TrimSpace(name))
	dashed := strings.ReplaceAll(lower, " ", "-")
	return dashed + ".yaml"
}

// flowStrings is a []string that serialises as a YAML flow sequence,
// e.g. [Ctrl, ","] instead of a multi-line block sequence.
type flowStrings []string

func (f flowStrings) MarshalYAML() (interface{}, error) {
	node := &yaml.Node{
		Kind:  yaml.SequenceNode,
		Style: yaml.FlowStyle,
		Tag:   "!!seq",
	}
	for _, s := range f {
		node.Content = append(node.Content, &yaml.Node{
			Kind:  yaml.ScalarNode,
			Tag:   "!!str",
			Value: s,
		})
	}
	return node, nil
}

// multiBindYAML serialises [][]string:
//   - 0 binds → omitted via omitempty
//   - 1 bind  → flat flow sequence: [j, k]
//   - 2+ binds → block sequence of flow sequences:
//   - [j, k]
//   - [ArrowUp, ArrowDown]
type multiBindYAML [][]string

func (m multiBindYAML) MarshalYAML() (interface{}, error) {
	if len(m) == 0 {
		return nil, nil
	}
	if len(m) == 1 {
		return flowStrings(m[0]).MarshalYAML()
	}
	outer := &yaml.Node{
		Kind: yaml.SequenceNode,
		Tag:  "!!seq",
	}
	for _, bind := range m {
		inner := &yaml.Node{
			Kind:  yaml.SequenceNode,
			Style: yaml.FlowStyle,
			Tag:   "!!seq",
		}
		for _, s := range bind {
			inner.Content = append(inner.Content, &yaml.Node{
				Kind:  yaml.ScalarNode,
				Tag:   "!!str",
				Value: s,
			})
		}
		outer.Content = append(outer.Content, inner)
	}
	return outer, nil
}

// yamlShortcut mirrors Shortcut but uses multiBindYAML so key slices are
// written as inline flow arrays (single bind) or block+flow (multiple binds).
type yamlShortcut struct {
	Desc  string        `yaml:"desc"`
	Keys  multiBindYAML `yaml:"keys,omitempty"`
	Linux multiBindYAML `yaml:"linux,omitempty"`
	MacOS multiBindYAML `yaml:"macos,omitempty"`
}

// yamlGroup mirrors Group but uses yamlShortcut.
type yamlGroup struct {
	Category  string         `yaml:"category"`
	Shortcuts []yamlShortcut `yaml:"shortcuts"`
}

// yamlAppConfig is the subset of AppConfig that gets persisted to YAML
// (excludes runtime-only fields like ModTime).
// Default is intentionally omitted — any user save strips the built-in registry flag.
type yamlAppConfig struct {
	App    string      `yaml:"app"`
	Icon   string      `yaml:"icon"`
	Tags   []string    `yaml:"tags,omitempty"`
	Groups []yamlGroup `yaml:"groups"`
}

// marshalApp produces clean YAML bytes for the given AppConfig.
// Key slices are written as inline flow sequences, e.g. ["Ctrl", ","].
func marshalApp(app AppConfig) ([]byte, error) {
	groups := make([]yamlGroup, len(app.Groups))
	for i, g := range app.Groups {
		shortcuts := make([]yamlShortcut, len(g.Shortcuts))
		for j, s := range g.Shortcuts {
			shortcuts[j] = yamlShortcut{
				Desc:  s.Desc,
				Keys:  multiBindYAML(s.Keys),
				Linux: multiBindYAML(s.Linux),
				MacOS: multiBindYAML(s.MacOS),
			}
		}
		groups[i] = yamlGroup{
			Category:  g.Category,
			Shortcuts: shortcuts,
		}
	}
	out := yamlAppConfig{
		App:    app.App,
		Icon:   app.Icon,
		Tags:   app.Tags,
		Groups: groups,
	}
	return yaml.Marshal(out)
}

// CreateApp writes a new YAML file for the given app config.
// Returns an error if a file for this app name already exists.
func CreateApp(app AppConfig) error {
	dir, err := ConfigDir()
	if err != nil {
		return err
	}

	filename := appNameToFilename(app.App)
	path := filepath.Join(dir, filename)

	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("app file %q already exists", filename)
	}

	data, err := marshalApp(app)
	if err != nil {
		return fmt.Errorf("failed to marshal app config: %w", err)
	}

	return os.WriteFile(path, data, 0644)
}

// UpdateApp replaces the YAML file for an app. If the app name changed,
// the old file is removed and a new one is created.
// Returns a non-empty warning string if the new filename would overwrite
// a different existing file — in that case the caller should confirm first
// by calling UpdateApp again with force=true (via UpdateAppForce).
func UpdateApp(oldAppName string, app AppConfig, force bool) (warning string, err error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}

	oldFilename := appNameToFilename(oldAppName)
	newFilename := appNameToFilename(app.App)
	oldPath := filepath.Join(dir, oldFilename)
	newPath := filepath.Join(dir, newFilename)

	// If name changed and new file already exists (and it's not the same file), warn.
	if oldFilename != newFilename {
		if _, err := os.Stat(newPath); err == nil && !force {
			return fmt.Sprintf(
				"A file named %q already exists. Saving will overwrite it.",
				newFilename,
			), nil
		}
	}

	data, err := marshalApp(app)
	if err != nil {
		return "", fmt.Errorf("failed to marshal app config: %w", err)
	}

	// Write new file first, then remove old if different.
	if err := os.WriteFile(newPath, data, 0644); err != nil {
		return "", fmt.Errorf("failed to write app config: %w", err)
	}

	if oldFilename != newFilename {
		// Remove old file (ignore error if it doesn't exist).
		os.Remove(oldPath)
	}

	return "", nil
}

// AppFilePath returns the absolute path to the YAML file for the given app name.
func AppFilePath(appName string) (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, appNameToFilename(appName)), nil
}

// DeleteApp removes the YAML file for the given app name.
func DeleteApp(appName string) error {
	dir, err := ConfigDir()
	if err != nil {
		return err
	}

	filename := appNameToFilename(appName)
	path := filepath.Join(dir, filename)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("app file %q does not exist", filename)
	}

	return os.Remove(path)
}
