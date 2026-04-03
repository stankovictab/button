package config

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// MultiBindKeys is [][]string that unmarshals from either a flat sequence
// [j, k] (treated as a single bind [[j, k]]) or a nested sequence
// [[j, k], [ArrowUp, ArrowDown]] (multiple binds).
type MultiBindKeys [][]string

func (m *MultiBindKeys) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind != yaml.SequenceNode {
		return fmt.Errorf("expected sequence, got kind %v", value.Kind)
	}
	if len(value.Content) == 0 {
		*m = nil
		return nil
	}
	// Flat single bind: [j, k] — first child is a scalar
	if value.Content[0].Kind == yaml.ScalarNode {
		var flat []string
		if err := value.Decode(&flat); err != nil {
			return err
		}
		if len(flat) > 0 {
			*m = [][]string{flat}
		}
		return nil
	}
	// Nested multiple binds: [[j, k], [ArrowUp, ArrowDown]]
	var nested [][]string
	if err := value.Decode(&nested); err != nil {
		return err
	}
	result := nested[:0]
	for _, bind := range nested {
		if len(bind) > 0 {
			result = append(result, bind)
		}
	}
	if len(result) > 0 {
		*m = result
	}
	return nil
}

// Shortcut represents a single keyboard shortcut with optional per-platform overrides.
// During resolution, the platform-specific field (Linux/MacOS) takes priority over Keys.
// Each field holds one or more key binds; multiple binds are alternatives for the same action.
type Shortcut struct {
	Desc  string        `yaml:"desc" json:"desc"`
	Keys  MultiBindKeys `yaml:"keys,omitempty" json:"keys"`
	Linux MultiBindKeys `yaml:"linux,omitempty" json:"linux,omitempty"`
	MacOS MultiBindKeys `yaml:"macos,omitempty" json:"macos,omitempty"`
}

// Group is a named category of shortcuts (e.g. "Navigation", "Editing").
type Group struct {
	Category  string     `yaml:"category" json:"category"`
	Shortcuts []Shortcut `yaml:"shortcuts" json:"shortcuts"`
}

// AppConfig is the top-level structure parsed from a single YAML file.
type AppConfig struct {
	App     string  `yaml:"app" json:"app"`
	Icon    string  `yaml:"icon" json:"icon"`
	Groups  []Group `yaml:"groups" json:"groups"`
	ModTime int64   `yaml:"-" json:"modTime"`
}

// AppsResponse is the payload returned to the frontend, containing
// successfully parsed apps and any warnings about skipped files.
type AppsResponse struct {
	Apps     []AppConfig `json:"apps"`
	Warnings []string    `json:"warnings"`
}
