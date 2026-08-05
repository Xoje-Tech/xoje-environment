package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestTUI(t *testing.T) {
	t.Run("Initial view", func(t *testing.T) {
		m := NewModel()
		view := m.View()
		if !strings.Contains(view, "xoje-environment") {
			t.Error("expected view to contain header")
		}
		if !strings.Contains(view, "> diagnose") {
			t.Error("expected cursor to be on first choice initially")
		}
	})

	t.Run("Navigation and Clamping", func(t *testing.T) {
		m := NewModel()
		
		// Move down
		raw, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
		m = raw.(Model)
		if m.cursor != 1 {
			t.Errorf("expected cursor 1 after 'j', got %d", m.cursor)
		}

		// Move down again
		raw, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
		m = raw.(Model)
		if m.cursor != 2 {
			t.Errorf("expected cursor 2 after 'j', got %d", m.cursor)
		}

		// Clamp down
		raw, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
		m = raw.(Model)
		if m.cursor != 2 {
			t.Errorf("expected cursor clamped at 2, got %d", m.cursor)
		}

		// Move up
		raw, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")})
		m = raw.(Model)
		if m.cursor != 1 {
			t.Errorf("expected cursor 1 after 'k', got %d", m.cursor)
		}
	})

	t.Run("Quit keys", func(t *testing.T) {
		m := NewModel()
		_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
		if cmd == nil {
			t.Error("expected 'q' to return a tea.Quit command")
		}
	})
}
