package skills

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"
)

const skillSchema = "xoje.skill/v1"
const ownershipMarker = "<!-- xoje-owned:v1 id=%s -->"

var identifierPattern = regexp.MustCompile(`^[a-z0-9-]+$`)
var semverPattern = regexp.MustCompile(`^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)(-[0-9A-Za-z-]+(\.[0-9A-Za-z-]+)*)?(\+[0-9A-Za-z-]+(\.[0-9A-Za-z-]+)*)?$`)

// Manifest is the validated set of canonical repository-owned skills.
type Manifest struct{ Skills []Skill }

type Skill struct {
	ID           string   `json:"id"`
	Name         string   `json:"-"`
	Description  string   `json:"-"`
	Category     string   `json:"category"`
	Version      string   `json:"version"`
	Metadata     Metadata `json:"metadata"`
	Capabilities []string `json:"capabilities"`
	Targets      []string `json:"targets"`
	Path         string   `json:"path"`
	SourceSHA256 string   `json:"source_sha256"`
}

type Metadata struct {
	Name    string `json:"name"`
	Summary string `json:"summary"`
}

func Load(repoRoot string) (*Manifest, error) {
	root := filepath.Join(repoRoot, "skills")
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, fmt.Errorf("read skills directory %q: %w", root, err)
	}
	result := &Manifest{}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		id := entry.Name()
		if !identifierPattern.MatchString(id) {
			return nil, fmt.Errorf("skill id %q is not lowercase path-safe", id)
		}
		path := filepath.Join(root, id, "SKILL.md")
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("skill %q: read SKILL.md: %w", id, err)
		}
		skill, err := parseSkill(id, data)
		if err != nil {
			return nil, fmt.Errorf("skill %q: %w", id, err)
		}
		skill.Path = filepath.ToSlash(filepath.Join("skills", id, "SKILL.md"))
		sum := sha256.Sum256(data)
		skill.SourceSHA256 = hex.EncodeToString(sum[:])
		result.Skills = append(result.Skills, skill)
	}
	if len(result.Skills) == 0 {
		return nil, fmt.Errorf("skills directory %q contains no canonical skills", root)
	}
	sort.Slice(result.Skills, func(i, j int) bool { return result.Skills[i].ID < result.Skills[j].ID })
	if err := checkLegacyManifest(repoRoot, result); err != nil {
		return nil, err
	}
	return result, nil
}

func Validate(repoRoot string) error { _, err := Load(repoRoot); return err }

func List(repoRoot string) ([]Skill, error) {
	m, err := Load(repoRoot)
	if err != nil {
		return nil, err
	}
	return m.Skills, nil
}

func parseSkill(id string, data []byte) (Skill, error) {
	var skill Skill
	if !utf8.Valid(data) {
		return skill, fmt.Errorf("source is not valid UTF-8")
	}
	text := strings.ReplaceAll(string(data), "\r\n", "\n")
	lines := strings.Split(text, "\n")
	if len(lines) < 4 || lines[0] != fmt.Sprintf(ownershipMarker, id) {
		return skill, fmt.Errorf("malformed ownership marker")
	}
	if lines[1] != "---" {
		return skill, fmt.Errorf("frontmatter must follow ownership marker")
	}
	end := -1
	for i := 2; i < len(lines); i++ {
		if lines[i] == "---" {
			end = i
			break
		}
	}
	if end < 0 {
		return skill, fmt.Errorf("missing frontmatter terminator")
	}
	values := map[string]string{}
	metadata := map[string]string{}
	metadataDeclared := false
	section := ""
	for i := 2; i < end; i++ {
		line := lines[i]
		if strings.TrimSpace(line) == "" {
			continue
		}
		indent := len(line) - len(strings.TrimLeft(line, " "))
		key, value, ok := strings.Cut(strings.TrimSpace(line), ":")
		if !ok || strings.TrimSpace(key) == "" {
			return skill, fmt.Errorf("invalid frontmatter line %d", i+1)
		}
		if indent > 0 {
			if section != "metadata" || indent != 2 {
				return skill, fmt.Errorf("invalid nested field %q", key)
			}
			if _, ok := metadata[key]; ok {
				return skill, fmt.Errorf("duplicate field %q", key)
			}
			metadata[key] = scalar(value)
			continue
		}
		section = key
		if key == "metadata" {
			if metadataDeclared {
				return skill, fmt.Errorf("duplicate field %q", key)
			}
			metadataDeclared = true
			if strings.TrimSpace(value) != "" {
				return skill, fmt.Errorf("metadata must be a mapping")
			}
			continue
		}
		if _, ok := values[key]; ok {
			return skill, fmt.Errorf("duplicate field %q", key)
		}
		values[key] = strings.TrimSpace(value)
	}
	allowed := map[string]bool{"schema": true, "id": true, "category": true, "version": true, "metadata": true, "capabilities": true, "targets": true}
	for key := range values {
		if !allowed[key] {
			return skill, fmt.Errorf("unknown field %q", key)
		}
	}
	for key := range metadata {
		if key != "name" && key != "summary" {
			return skill, fmt.Errorf("unknown metadata field %q", key)
		}
	}
	if values["schema"] != skillSchema {
		return skill, fmt.Errorf("unsupported schema %q", values["schema"])
	}
	if values["id"] != id {
		return skill, fmt.Errorf("frontmatter id %q does not match marker", values["id"])
	}
	if !identifierPattern.MatchString(values["category"]) {
		return skill, fmt.Errorf("category must be lowercase path-safe")
	}
	if !isSemVer(values["version"]) {
		return skill, fmt.Errorf("version %q is not SemVer", values["version"])
	}
	if strings.TrimSpace(metadata["name"]) == "" || strings.TrimSpace(metadata["summary"]) == "" {
		return skill, fmt.Errorf("metadata name and summary are required")
	}
	caps, err := listValue(values["capabilities"])
	if err != nil {
		return skill, fmt.Errorf("capabilities: %w", err)
	}
	for _, capability := range caps {
		if capability != "documentation" && capability != "workflow" && capability != "rpi-reference" {
			return skill, fmt.Errorf("unsupported capability %q", capability)
		}
	}
	targets, err := listValue(values["targets"])
	if err != nil {
		return skill, fmt.Errorf("targets: %w", err)
	}
	for _, target := range targets {
		return skill, fmt.Errorf("unsupported target %q", target)
	}
	skill = Skill{ID: id, Name: metadata["name"], Description: metadata["summary"], Category: values["category"], Version: values["version"], Metadata: Metadata{Name: metadata["name"], Summary: metadata["summary"]}, Capabilities: caps, Targets: targets}
	return skill, nil
}

func isSemVer(version string) bool {
	if !semverPattern.MatchString(version) {
		return false
	}
	base := strings.SplitN(version, "+", 2)[0]
	parts := strings.SplitN(base, "-", 2)
	if len(parts) == 1 {
		return true
	}
	for _, identifier := range strings.Split(parts[1], ".") {
		if len(identifier) > 1 && identifier[0] == '0' && strings.Trim(identifier, "0123456789") == "" {
			return false
		}
	}
	return true
}

func scalar(value string) string {
	value = strings.TrimSpace(value)
	if len(value) >= 2 && ((value[0] == '"' && value[len(value)-1] == '"') || (value[0] == '\'' && value[len(value)-1] == '\'')) {
		if value[0] == '"' {
			if v, err := strconv.Unquote(value); err == nil {
				return v
			}
		}
		return value[1 : len(value)-1]
	}
	return value
}
func listValue(value string) ([]string, error) {
	value = strings.TrimSpace(value)
	if value == "[]" {
		return []string{}, nil
	}
	if len(value) < 2 || value[0] != '[' || value[len(value)-1] != ']' {
		return nil, fmt.Errorf("must be an inline list")
	}
	var out []string
	for _, item := range strings.Split(value[1:len(value)-1], ",") {
		item = scalar(item)
		if item == "" {
			return nil, fmt.Errorf("empty value")
		}
		out = append(out, item)
	}
	sort.Strings(out)
	for i := 1; i < len(out); i++ {
		if out[i] == out[i-1] {
			return nil, fmt.Errorf("duplicate value %q", out[i])
		}
	}
	return out, nil
}
