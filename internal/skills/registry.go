package skills

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const registrySchema = "xoje.skill-registry/v1"

type Registry struct {
	Schema string  `json:"schema"`
	Skills []Skill `json:"skills"`
}

// GenerateRegistry derives both registry views from validated canonical source.
func GenerateRegistry(repoRoot string) (*Registry, error) {
	manifest, err := Load(repoRoot)
	if err != nil {
		return nil, err
	}
	registry := &Registry{Schema: registrySchema, Skills: manifest.Skills}
	data, err := json.MarshalIndent(registry, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal skill registry: %w", err)
	}
	data = append(data, '\n')
	atl := filepath.Join(repoRoot, ".atl")
	if err := os.MkdirAll(atl, 0o755); err != nil {
		return nil, fmt.Errorf("create registry directory: %w", err)
	}
	if err := os.WriteFile(filepath.Join(atl, "skill-registry.json"), data, 0o644); err != nil {
		return nil, fmt.Errorf("write JSON registry: %w", err)
	}
	markdown := renderMarkdown(registry)
	if err := os.WriteFile(filepath.Join(atl, "skill-registry.md"), []byte(markdown), 0o644); err != nil {
		return nil, fmt.Errorf("write Markdown registry: %w", err)
	}
	return registry, nil
}

// VerifyRegistry reports source or derived-output drift and never repairs it.
func VerifyRegistry(repoRoot string) error {
	manifest, err := Load(repoRoot)
	if err != nil {
		return fmt.Errorf("registry drift: canonical source changed or invalid: %w", err)
	}
	want := &Registry{Schema: registrySchema, Skills: manifest.Skills}
	jsonWant, err := json.MarshalIndent(want, "", "  ")
	if err != nil {
		return err
	}
	jsonWant = append(jsonWant, '\n')
	jsonPath := filepath.Join(repoRoot, ".atl", "skill-registry.json")
	got, err := os.ReadFile(jsonPath)
	if err != nil {
		return fmt.Errorf("registry drift: read JSON registry: %w", err)
	}
	if !bytes.Equal(got, jsonWant) {
		return fmt.Errorf("registry drift: %s differs from canonical source", jsonPath)
	}
	mdPath := filepath.Join(repoRoot, ".atl", "skill-registry.md")
	mdGot, err := os.ReadFile(mdPath)
	if err != nil {
		return fmt.Errorf("registry drift: read Markdown registry: %w", err)
	}
	if string(mdGot) != renderMarkdown(want) {
		return fmt.Errorf("registry drift: %s differs from canonical source", mdPath)
	}
	return nil
}

func renderMarkdown(registry *Registry) string {
	var b strings.Builder
	b.WriteString("# Xoje Skill Registry\n\n")
	b.WriteString("| ID | Category | Version | Path | SHA-256 |\n|---|---|---|---|---|\n")
	for _, skill := range registry.Skills {
		fmt.Fprintf(&b, "| %s | %s | %s | `%s` | `%s` |\n", skill.ID, skill.Category, skill.Version, skill.Path, skill.SourceSHA256)
	}
	return b.String()
}

func fingerprint(data []byte) string { sum := sha256.Sum256(data); return hex.EncodeToString(sum[:]) }
