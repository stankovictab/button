package registry

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// RegistryEntry holds the metadata for a single registry app.
type RegistryEntry struct {
	Filename string   `json:"filename"`
	App      string   `json:"app"`
	Icon     string   `json:"icon"`
	Tags     []string `json:"tags"`
}

// registryApp is used only for unmarshalling YAML metadata.
type registryApp struct {
	App  string   `yaml:"app"`
	Icon string   `yaml:"icon"`
	Tags []string `yaml:"tags"`
}

// Registry holds the embedded filesystem and provides access to registry apps.
type Registry struct {
	fs fs.FS
}

// New creates a Registry backed by the given filesystem.
func New(f fs.FS) *Registry {
	return &Registry{fs: f}
}

// ListApps reads all embedded YAML files and returns their metadata.
func (r *Registry) ListApps() ([]RegistryEntry, error) {
	entries, err := fs.ReadDir(r.fs, ".")
	if err != nil {
		return nil, fmt.Errorf("failed to read registry: %w", err)
	}

	var result []RegistryEntry
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, ".yaml") {
			continue
		}

		data, err := fs.ReadFile(r.fs, name)
		if err != nil {
			continue
		}

		var app registryApp
		if err := yaml.Unmarshal(data, &app); err != nil || app.App == "" {
			continue
		}

		result = append(result, RegistryEntry{
			Filename: name,
			App:      app.App,
			Icon:     app.Icon,
			Tags:     app.Tags,
		})
	}

	return result, nil
}

// GetAppYAML returns the raw YAML bytes for a registry app by filename.
func (r *Registry) GetAppYAML(filename string) ([]byte, error) {
	clean := filepath.Base(filename)
	data, err := fs.ReadFile(r.fs, clean)
	if err != nil {
		return nil, fmt.Errorf("registry app %q not found: %w", clean, err)
	}
	return data, nil
}
