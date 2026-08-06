package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestTUI(t *testing.T) {
	t.Run("Initial view", func(t *testing.T) {
		m := NewModel([]string{})
		view := m.View()
		if !strings.Contains(view, "xoje-environment") {
			t.Error("expected view to contain header")
		}
		if !strings.Contains(view, "doctor") {
			t.Error("expected view to contain 'doctor' choice")
		}
	})

	t.Run("Navigation and Selection", func(t *testing.T) {
		m := NewModel([]string{})
		
		// 1. Move down to "install gentle-ai" (index 1)
		raw, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
		m = raw.(Model)
		
		// 2. Press Enter
		raw, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		m = raw.(Model)

		// Verification A: Did it return the Quit command?
		if cmd == nil {
			t.Error("expected Enter to return a tea.Quit command")
		}

		// Verification B: Is the choice recorded correctly?
		expected := "install gentle-ai"
		if m.Choice != expected {
			t.Errorf("expected Choice to be '%s', got '%s'", expected, m.Choice)
		}
	})

	t.Run("Quit keys", func(t *testing.T) {
		m := NewModel([]string{})
		_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
		if cmd == nil {
			t.Error("expected 'q' to return a tea.Quit command")
		}
	})
}
