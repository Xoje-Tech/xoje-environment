package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	cursor  int
	choices []string
	Choice  string // Field to record the final selection
}

func NewModel() Model {
	return Model{
		choices: []string{"diagnose", "install gentle-ai", "exit"},
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter", " ":
			// Record selection
			m.Choice = m.choices[m.cursor]
			// Quit the TUI to return control to main
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) View() string {
	var s strings.Builder
	s.WriteString("=================================\n")
	s.WriteString("  xoje-environment CLI & TUI  \n")
	s.WriteString("=================================\n\n")

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		s.WriteString(fmt.Sprintf("%s %s\n", cursor, choice))
	}

	s.WriteString("\nPress Enter to select, q/ctrl+c to exit.\n")
	return s.String()
}
