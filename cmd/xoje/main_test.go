package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Xoje-Tech/xoje-environment/internal/doctor"
)

func TestRun(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	// Note: Commands that trigger TUI (doctor, tui) will fail in non-interactive tests
	// but we can at least verify they are dispatched correctly or fail gracefully.

	t.Run("Install unknown tool", func(t *testing.T) {
		err := Run([]string{"install", "unknown-tool"}, configPath)
		if err == nil {
			t.Error("expected error for unknown tool")
		}
	})

	t.Run("Update command", func(t *testing.T) {
		err := Run([]string{"update"}, configPath)
		if err != nil {
			t.Errorf("Run update failed: %v", err)
		}
	})

	t.Run("Invalid subcommand", func(t *testing.T) {
		err := Run([]string{"frobnicate"}, configPath)
		if err == nil {
			t.Error("expected error for unknown command")
		}
	})

	t.Run("Load config error", func(t *testing.T) {
		// Pass a directory as config path to trigger read error
		err := Run([]string{"doctor"}, tmpDir)
		if err == nil {
			t.Error("expected error loading config from directory path")
		}
	})
}

func TestDoctorReportIncludesAggregateReadiness(t *testing.T) {
	report := formatDoctorReport([]doctor.Result{{
		Name: "Example", Status: doctor.StatusWarn, Detail: "warning",
	}})
	if !strings.Contains(report, "READINESS: DEGRADED") {
		t.Fatalf("report missing aggregate readiness: %q", report)
	}
}

func TestFormatDoctorReportMixedResults(t *testing.T) {
	results := []doctor.Result{
		{Name: "Passing check", Status: doctor.StatusPass, Detail: "pass detail", Remedy: "pass remedy"},
		{Name: "Warning check", Status: doctor.StatusWarn, Detail: "warning detail", Remedy: "warning remedy"},
		{Name: "Failed check", Status: doctor.StatusFail, Detail: "failure detail", Remedy: "failure remedy"},
	}

	report := formatDoctorReport(results)
	for _, want := range []string{
		"✅", "⚠️", "❌",
		"pass detail", "warning detail", "failure detail",
		"pass remedy", "warning remedy", "failure remedy",
		"READINESS: NOT READY",
	} {
		if !strings.Contains(report, want) {
			t.Errorf("report missing %q: %q", want, report)
		}
	}
}

func TestRunDoctorExecutesComprehensiveDiagnosis(t *testing.T) {
	home := t.TempDir()
	binDir := filepath.Join(home, ".local", "bin")
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		t.Fatal(err)
	}
	for name, output := range map[string]string{
		"go":        "go version go1.22.0 linux/amd64",
		"node":      "v22.0.0",
		"git":       "git version 2.46.0",
		"engram":    `{"status":"ok"}`,
		"gentle-ai": "healthy",
	} {
		path := filepath.Join(binDir, name)
		content := "#!/bin/sh\nprintf '%s\\n' '" + output + "'\n"
		if err := os.WriteFile(path, []byte(content), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	t.Setenv("HOME", home)
	t.Setenv("PATH", binDir)

	configPath := filepath.Join(home, ".config", "xoje", "config.json")
	outputPath := filepath.Join(t.TempDir(), "doctor-output.txt")
	output, err := os.Create(outputPath)
	if err != nil {
		t.Fatal(err)
	}
	originalStdout := os.Stdout
	os.Stdout = output
	runErr := Run([]string{"doctor"}, configPath)
	os.Stdout = originalStdout
	if closeErr := output.Close(); closeErr != nil {
		t.Fatal(closeErr)
	}
	if runErr != nil {
		t.Fatalf("Run(doctor) error = %v", runErr)
	}

	reportBytes, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatal(err)
	}
	report := string(reportBytes)
	for _, want := range []string{
		"Go Version",
		"Node Runtime",
		"Git Runtime",
		"System PATH Configuration",
		"https://proxy.golang.org",
		"https://github.com",
		"Configuration File",
		"Managed Tools",
		"Engram Memory Health",
		"Gentle-AI Ecosystem Health",
		"READINESS:",
	} {
		if !strings.Contains(report, want) {
			t.Errorf("doctor report missing %q: %q", want, report)
		}
	}
	if got := strings.Count(report, "Network Connectivity"); got != 2 {
		t.Errorf("doctor report contains %d network results, want 2: %q", got, report)
	}
}
