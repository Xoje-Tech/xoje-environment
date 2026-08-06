package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/Xoje-Tech/xoje-environment/internal/doctor"
	"github.com/Xoje-Tech/xoje-environment/internal/tui/styles"
)

type Model struct {
	cursor  int
	choices []string
	Choice  string
	
	isDoctoring bool
	results     []doctor.Result
	tools       []string 
}

func NewModel(tools []string) Model {
	return Model{
		choices: []string{"doctor", "install gentle-ai", "update", "exit"},
		tools:   tools,
	}
}

func (m Model) Init() tea.Cmd {
	if m.Choice == "doctor" {
		m.isDoctoring = true
		return runDoctor(m.tools)
	}
	return nil
}

type doctorResultsMsg []doctor.Result

func runDoctor(tools []string) tea.Cmd {
	return func() tea.Msg {
		checks := []doctor.Check{
			&doctor.GoVersionCheck{Required: "1.22"},
			&doctor.EnvPathCheck{}, // Added PATH check
			&doctor.NetworkCheck{Target: "google.com"},
			&doctor.ConfigCheck{Path: "config.json"},
			&doctor.ToolsCheck{Tools: tools},
		}
		runner := doctor.NewRunner(checks)
		return doctorResultsMsg(runner.Run())
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case doctorResultsMsg:
		m.results = msg
		return m, nil

	case tea.KeyMsg:
		if m.isDoctoring {
			if msg.String() == "q" || msg.String() == "esc" || msg.String() == "ctrl+c" {
				if m.Choice == "doctor" {
					return m, tea.Quit
				}
				m.isDoctoring = false
				m.results = nil
				return m, nil
			}
			return m, nil
		}

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
			m.Choice = m.choices[m.cursor]
			if m.Choice == "doctor" {
				m.isDoctoring = true
				m.results = nil
				return m, runDoctor(m.tools)
			}
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) View() string {
	if m.isDoctoring {
		return m.doctorView()
	}

	var s strings.Builder
	s.WriteString(styles.TitleStyle.Render("================================="))
	s.WriteString("\n")
	s.WriteString(styles.TitleStyle.Render("  xoje-environment CLI & TUI  "))
	s.WriteString("\n")
	s.WriteString(styles.TitleStyle.Render("================================="))
	s.WriteString("\n\n")

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = styles.IrisStyle.Render(">")
		}
		
		label := choice
		if m.cursor == i {
			label = styles.IrisStyle.Render(choice)
		}
		s.WriteString(fmt.Sprintf("%s %s\n", cursor, label))
	}

	s.WriteString(styles.MutedStyle.Render("\nPress Enter to select, q/ctrl+c to exit.\n"))
	return s.String()
}

func (m Model) doctorView() string {
	var s strings.Builder
	s.WriteString(styles.TitleStyle.Render("--- 🩺 System Doctor ---"))
	s.WriteString("\n\n")

	if m.results == nil {
		s.WriteString(styles.IrisStyle.Render("Running diagnostics..."))
		s.WriteString("\n")
		return styles.BoxStyle.Render(s.String())
	}

	for _, res := range m.results {
		statusIcon := styles.SuccessStyle.Render("✅")
		if res.Status == doctor.StatusFail {
			statusIcon = styles.ErrorStyle.Render("❌")
		} else if res.Status == doctor.StatusWarn {
			statusIcon = styles.WarnStyle.Render("⚠️")
		}

		s.WriteString(fmt.Sprintf("%s %s\n", statusIcon, styles.IrisStyle.Render(res.Name)))
		s.WriteString(fmt.Sprintf("   %s\n", styles.MutedStyle.Render(res.Detail)))
		if res.Remedy != "" {
			s.WriteString(fmt.Sprintf("   %s %s\n", styles.GoldStyle.Render("Remedy:"), res.Remedy))
		}
		s.WriteString("\n")
	}

	s.WriteString(styles.MutedStyle.Render("Press 'q' or 'esc' to return to menu."))
	
	return styles.BoxStyle.Render(s.String())
}
