package skills

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeCanonicalSkill(t *testing.T, root, id, category, version string, body string) {
	t.Helper()
	path := filepath.Join(root, "skills", id, "SKILL.md")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	content := "<!-- xoje-owned:v1 id=" + id + " -->\n---\n" +
		"schema: xoje.skill/v1\n" +
		"id: " + id + "\n" +
		"category: " + category + "\n" +
		"version: " + version + "\n" +
		"metadata:\n  name: " + id + "\n  summary: A summary\n" +
		"capabilities: [documentation]\n" +
		"targets: []\n---\n\n" + body + "\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestLoadRejectsCanonicalContractViolations(t *testing.T) {
	tests := []struct {
		name, content, want string
	}{
		{"marker", "<!-- wrong id -->\n---\nschema: xoje.skill/v1\nid: demo\ncategory: tools\nversion: 1.0.0\nmetadata:\n  name: demo\n  summary: Demo\ncapabilities: [documentation]\ntargets: []\n---\n", "marker"},
		{"unknown field", "<!-- xoje-owned:v1 id=demo -->\n---\nschema: xoje.skill/v1\nid: demo\ncategory: tools\nversion: 1.0.0\nmetadata:\n  name: demo\n  summary: Demo\ncapabilities: [documentation]\ntargets: []\nextra: nope\n---\n", "unknown"},
		{"version", "<!-- xoje-owned:v1 id=demo -->\n---\nschema: xoje.skill/v2\nid: demo\ncategory: tools\nversion: 1.0.0\nmetadata:\n  name: demo\n  summary: Demo\ncapabilities: [documentation]\ntargets: []\n---\n", "schema"},
		{"capability", "<!-- xoje-owned:v1 id=demo -->\n---\nschema: xoje.skill/v1\nid: demo\ncategory: tools\nversion: 1.0.0\nmetadata:\n  name: demo\n  summary: Demo\ncapabilities: [shell]\ntargets: []\n---\n", "capability"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			if err := os.MkdirAll(filepath.Join(root, "skills", "demo"), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(root, "skills", "demo", "SKILL.md"), []byte(tt.content), 0o644); err != nil {
				t.Fatal(err)
			}
			if _, err := Load(root); err == nil || !strings.Contains(strings.ToLower(err.Error()), tt.want) {
				t.Fatalf("Load() error = %v, want %q", err, tt.want)
			}
		})
	}
}

func TestLoadAcceptsLegacyManifestMetadata(t *testing.T) {
	root := t.TempDir()
	writeCanonicalSkill(t, root, "demo", "tools", "1.2.3", "body")
	manifest := `{"version":1,"skills":[{"name":"demo","description":"A summary","path":"demo","metadata":{"owner":"xoje-environment"}}]}`
	if err := os.WriteFile(filepath.Join(root, "skills", "manifest.json"), []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(root); err != nil {
		t.Fatalf("Load() error = %v", err)
	}
}

func TestLoadTreatsMarkdownAsInertData(t *testing.T) {
	root := t.TempDir()
	writeCanonicalSkill(t, root, "demo", "tools", "1.2.3", "```sh\nrm -rf /\n```")
	if _, err := Load(root); err != nil {
		t.Fatalf("Load() error = %v", err)
	}
}

func TestGenerateRegistryIsDeterministicAndDetectsDrift(t *testing.T) {
	root := t.TempDir()
	writeCanonicalSkill(t, root, "zeta", "tools", "1.0.0", "Z")
	writeCanonicalSkill(t, root, "alpha", "docs", "1.0.0", "A")
	if _, err := GenerateRegistry(root); err != nil {
		t.Fatal(err)
	}
	firstJSON, err := os.ReadFile(filepath.Join(root, ".atl", "skill-registry.json"))
	if err != nil {
		t.Fatal(err)
	}
	firstMD, err := os.ReadFile(filepath.Join(root, ".atl", "skill-registry.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(firstJSON) != string(mustRead(t, filepath.Join(root, ".atl", "skill-registry.json"))) {
		t.Fatal("registry changed on repeated read")
	}
	if !strings.HasSuffix(string(firstJSON), "\n") || strings.Contains(string(firstJSON), "\r") || !strings.HasSuffix(string(firstMD), "\n") {
		t.Fatal("registry output is not canonical UTF-8/LF/final-newline")
	}
	if !strings.Contains(string(firstJSON), `"id": "alpha"`) || strings.Index(string(firstJSON), `"id": "alpha"`) > strings.Index(string(firstJSON), `"id": "zeta"`) {
		t.Fatal("registry is not bytewise sorted")
	}
	if err := VerifyRegistry(root); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "skills", "alpha", "SKILL.md"), []byte("changed\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := VerifyRegistry(root); err == nil || !strings.Contains(strings.ToLower(err.Error()), "drift") {
		t.Fatalf("VerifyRegistry() error = %v, want drift", err)
	}
}

func mustRead(t *testing.T, path string) []byte {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return b
}
