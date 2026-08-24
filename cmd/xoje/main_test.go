package main

import (
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
