package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunSkillsUsesCurrentWorkingDirectory(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "skills", "demo"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "skills", "manifest.json"), []byte(`{"version":1,"skills":[{"name":"demo","description":"Demo","path":"demo"}]}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "skills", "demo", "SKILL.md"), []byte("<!-- xoje-owned:v1 id=demo -->\n---\nschema: xoje.skill/v1\nid: demo\ncategory: tools\nversion: 1.0.0\nmetadata:\n  name: demo\n  summary: Demo\ncapabilities: [documentation]\ntargets: []\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	old, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(old) })

	if err := Run([]string{"skills", "validate"}, filepath.Join(t.TempDir(), "unused-config.json")); err != nil {
		t.Fatalf("Run(skills validate) error = %v", err)
	}
}
