package skills

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeProjectManifest(t *testing.T, root, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(root, "skills"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "skills", "manifest.json"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func writeSkill(t *testing.T, root, name, description string) {
	t.Helper()
	path := filepath.Join(root, "skills", name, "SKILL.md")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	content := "<!-- xoje-owned:v1 id=" + name + " -->\n---\nschema: xoje.skill/v1\nid: " + name + "\ncategory: tools\nversion: 1.0.0\nmetadata:\n  name: " + name + "\n  summary: " + description + "\ncapabilities: [documentation]\ntargets: []\n---\n\n# " + name + "\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestLoadValidManifest(t *testing.T) {
	root := t.TempDir()
	writeProjectManifest(t, root, `{"version":1,"skills":[{"name":"go-conventions","description":"Go conventions","path":"go-conventions"}]}`)
	writeSkill(t, root, "go-conventions", "Go conventions")
	manifest, err := Load(root)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if len(manifest.Skills) != 1 || manifest.Skills[0].ID != "go-conventions" {
		t.Fatalf("Load() manifest = %#v", manifest)
	}
}
func TestValidateRejectsMissingSkillDocument(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "skills", "missing"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := Validate(root); err == nil || !strings.Contains(err.Error(), "SKILL.md") {
		t.Fatalf("Validate() error = %v, want missing SKILL.md", err)
	}
}
func TestValidateRejectsDuplicateNamesAndTraversal(t *testing.T) {
	root := t.TempDir()
	writeSkill(t, root, "same", "One")
	if err := os.MkdirAll(filepath.Join(root, "skills", "Same"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := Validate(root); err == nil || !strings.Contains(err.Error(), "lowercase") {
		t.Fatalf("Validate() error = %v, want lowercase id", err)
	}
}
func TestValidateRequiresFrontmatterNameAndDescription(t *testing.T) {
	root := t.TempDir()
	writeSkill(t, root, "go-conventions", "Go conventions")
	path := filepath.Join(root, "skills", "go-conventions", "SKILL.md")
	if err := os.WriteFile(path, []byte("bad\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := Validate(root); err == nil || !strings.Contains(err.Error(), "marker") {
		t.Fatalf("Validate() error = %v, want marker error", err)
	}
}

func TestValidateRejectsDuplicateTopLevelMetadata(t *testing.T) {
	root := t.TempDir()
	writeSkill(t, root, "demo", "Demo")
	path := filepath.Join(root, "skills", "demo", "SKILL.md")
	content := "<!-- xoje-owned:v1 id=demo -->\n---\nschema: xoje.skill/v1\nid: demo\ncategory: tools\nversion: 1.0.0\nmetadata:\n  name: demo\n  summary: Demo\nmetadata:\n  name: duplicate\ncapabilities: [documentation]\ntargets: []\n---\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := Validate(root); err == nil || !strings.Contains(err.Error(), `duplicate field "metadata"`) {
		t.Fatalf("Validate() error = %v, want duplicate metadata error", err)
	}
}

func TestValidateRejectsNumericPrereleaseLeadingZero(t *testing.T) {
	root := t.TempDir()
	writeSkill(t, root, "demo", "Demo")
	path := filepath.Join(root, "skills", "demo", "SKILL.md")
	content := "<!-- xoje-owned:v1 id=demo -->\n---\nschema: xoje.skill/v1\nid: demo\ncategory: tools\nversion: 1.0.0-01\nmetadata:\n  name: demo\n  summary: Demo\ncapabilities: [documentation]\ntargets: []\n---\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := Validate(root); err == nil || !strings.Contains(err.Error(), "not SemVer") {
		t.Fatalf("Validate() error = %v, want SemVer error", err)
	}
}

func TestListReturnsStableNameOrder(t *testing.T) {
	root := t.TempDir()
	writeSkill(t, root, "zeta", "Z")
	writeSkill(t, root, "alpha", "A")
	got, err := List(root)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if got[0].ID != "alpha" || got[1].ID != "zeta" {
		t.Fatalf("List() order = %#v", got)
	}
}
