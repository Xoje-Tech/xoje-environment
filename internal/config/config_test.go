package config

import (
	"path/filepath"
	"testing"
)

func TestLoadOrCreate(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	// Test creation of default config
	cfg, err := LoadOrCreate(configPath)
	if err != nil {
		t.Fatalf("LoadOrCreate failed: %v", err)
	}

	if cfg.ActivePersona != "gandalf" {
		t.Errorf("expected default persona 'gandalf', got '%s'", cfg.ActivePersona)
	}

	if len(cfg.InstalledTools) != 0 {
		t.Errorf("expected 0 installed tools, got %d", len(cfg.InstalledTools))
	}

	// Test persistence
	cfg.InstalledTools = append(cfg.InstalledTools, "gentle-ai")
	if err := cfg.Save(); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Reload and verify
	cfg2, err := LoadOrCreate(configPath)
	if err != nil {
		t.Fatalf("Reload failed: %v", err)
	}

	if len(cfg2.InstalledTools) != 1 || cfg2.InstalledTools[0] != "gentle-ai" {
		t.Errorf("reloaded config doesn't match; got %v", cfg2.InstalledTools)
	}
}

func TestDefaultConfigPath(t *testing.T) {
	home := "/home/hermes"
	t.Setenv("HOME", home)
	
	expected := filepath.Join(home, ".config", "xoje", "config.json")
	got, err := DefaultConfigPath()
	if err != nil {
		t.Fatalf("DefaultConfigPath failed: %v", err)
	}
	
	if got != expected {
		t.Errorf("expected %s, got %s", expected, got)
	}
}
