package tui

import (
	"testing"

	"github.com/Xoje-Tech/xoje-environment/internal/doctor"
	"github.com/Xoje-Tech/xoje-environment/internal/install"
	tea "github.com/charmbracelet/bubbletea"
)

func TestTUISelection(t *testing.T) {
	m := NewModel([]string{"gentle-ai"})

	// Move down twice to reach "update fleet".
	m, _ = updateModel(m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	m, _ = updateModel(m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	m, cmd := updateModel(m, tea.KeyMsg{Type: tea.KeyEnter})

	if m.Choice != "update-all" {
		t.Errorf("Expected Choice 'update-all', got '%s'", m.Choice)
	}
	if cmd == nil {
		t.Error("Expected a non-nil command (tea.Quit)")
	}

	m = NewModel([]string{})
	m.state = stateInstallTools
	m.toolStatuses = []install.ToolStatus{{Tool: install.Registry[0], Installed: false}}
	m, _ = updateModel(m, tea.KeyMsg{Type: tea.KeyEnter})

	expected := "install " + install.Registry[0].Name
	if m.Choice != expected {
		t.Errorf("Expected Choice '%s', got '%s'", expected, m.Choice)
	}

	m = NewModel([]string{})
	m.state = stateInstallTools
	m.toolStatuses = []install.ToolStatus{{Tool: install.Registry[0], Installed: true}}
	m, _ = updateModel(m, tea.KeyMsg{Type: tea.KeyEnter})

	expectedUpdate := "update " + install.Registry[0].Name
	if m.Choice != expectedUpdate {
		t.Errorf("Expected Choice '%s', got '%s'", expectedUpdate, m.Choice)
	}
}

func TestDoctorRemedySelection(t *testing.T) {
	m := NewModel(nil)
	m.state = stateDoctor
	m.results = []doctor.Result{
		{Name: "PATH", Status: doctor.StatusFail, Remedy: "Edit your shell profile."},
		{Name: "Tools", Status: doctor.StatusWarn, Remedy: "Restore missing tools.", RemedyAction: "update-all"},
	}

	m, cmd := updateModel(m, tea.KeyMsg{Type: tea.KeyEnter})
	if m.Choice != "" || cmd != nil {
		t.Fatalf("non-actionable remedy executed: choice=%q cmd=%v", m.Choice, cmd)
	}

	m, _ = updateModel(m, tea.KeyMsg{Type: tea.KeyDown})
	m, cmd = updateModel(m, tea.KeyMsg{Type: tea.KeyEnter})
	if m.Choice != "update-all" {
		t.Fatalf("choice = %q, want update-all", m.Choice)
	}
	if cmd == nil {
		t.Fatal("actionable remedy should quit the TUI for dispatch")
	}
}

func updateModel(m Model, msg tea.Msg) (Model, tea.Cmd) {
	newModel, cmd := m.Update(msg)
	return newModel.(Model), cmd
}
