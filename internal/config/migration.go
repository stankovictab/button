package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// MigrateYMLToYAML renames all .yml files in the config directory to .yaml.
// Returns the number of files migrated and any warnings for files that could
// not be renamed (e.g. because a .yaml with the same base name already exists).
func MigrateYMLToYAML() (int, []string) {
	dir, err := ConfigDir()
	if err != nil {
		return 0, []string{fmt.Sprintf("migration: could not determine config directory — %s", err)}
	}

	matches, err := filepath.Glob(filepath.Join(dir, "*.yml"))
	if err != nil || len(matches) == 0 {
		return 0, nil
	}

	migrated := 0
	var warnings []string

	for _, ymlPath := range matches {
		base := filepath.Base(ymlPath)
		newName := strings.TrimSuffix(base, ".yml") + ".yaml"
		yamlPath := filepath.Join(dir, newName)

		// If a .yaml file with the same name already exists, skip and warn.
		if _, err := os.Stat(yamlPath); err == nil {
			warnings = append(warnings, fmt.Sprintf(
				"migration: skipped %s — %s already exists",
				base, newName,
			))
			continue
		}

		if err := os.Rename(ymlPath, yamlPath); err != nil {
			warnings = append(warnings, fmt.Sprintf(
				"migration: could not rename %s — %s",
				base, err,
			))
			continue
		}

		migrated++
	}

	return migrated, warnings
}
