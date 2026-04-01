package config

// Shortcut represents a single keyboard shortcut with optional per-platform overrides.
// During resolution, the platform-specific field (Linux/MacOS) takes priority over Keys.
type Shortcut struct {
	Desc  string   `yaml:"desc" json:"desc"`
	Keys  []string `yaml:"keys,omitempty" json:"keys"`
	Linux []string `yaml:"linux,omitempty" json:"linux,omitempty"`
	MacOS []string `yaml:"macos,omitempty" json:"macos,omitempty"`
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
