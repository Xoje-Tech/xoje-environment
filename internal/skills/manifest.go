package skills

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

type legacyManifest struct {
	Version int           `json:"version"`
	Skills  []legacySkill `json:"skills"`
}
type legacySkill struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Path        string            `json:"path"`
	Metadata    map[string]string `json:"metadata"`
}

func checkLegacyManifest(repoRoot string, manifest *Manifest) error {
	path := filepath.Join(repoRoot, "skills", "manifest.json")
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("read skills manifest %q: %w", path, err)
	}
	var legacy legacyManifest
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&legacy); err != nil {
		return fmt.Errorf("parse skills manifest %q: %w", path, err)
	}
	if legacy.Version != 1 || len(legacy.Skills) != len(manifest.Skills) {
		return fmt.Errorf("skills manifest is inconsistent with canonical sources")
	}
	sort.Slice(legacy.Skills, func(i, j int) bool { return legacy.Skills[i].Name < legacy.Skills[j].Name })
	for i, item := range legacy.Skills {
		if i >= len(manifest.Skills) || item.Name != manifest.Skills[i].ID || item.Description != manifest.Skills[i].Description || item.Path != filepath.Base(filepath.Dir(filepath.FromSlash(manifest.Skills[i].Path))) {
			return fmt.Errorf("skills manifest is inconsistent with canonical skill %q", item.Name)
		}
	}
	return nil
}
