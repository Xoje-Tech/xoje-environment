package tui

import (
	"fmt"
	"strings"

	"github.com/Xoje-Tech/xoje-environment/internal/config"
	"github.com/Xoje-Tech/xoje-environment/internal/doctor"
	"github.com/Xoje-Tech/xoje-environment/internal/install"
	"github.com/Xoje-Tech/xoje-environment/internal/tui/styles"
	tea "github.com/charmbracelet/bubbletea"
)

type state int

const (
	stateMenu state = iota
	stateDoctor
	stateInstallTools
)

type Model struct {
	state   state
	cursor  int
	choices []string
	Choice  string

	// Doctor data
	results []doctor.Result
	tools   []string

	// Tool list data
	toolStatuses []install.ToolStatus
	isLoading    bool
}

func NewModel(tools []string) Model {
	return Model{
		state:   stateMenu,
		choices: []string{"doctor", "install tools", "update fleet", "exit"},
		tools:   tools,
	}
}

func (m Model) Init() tea.Cmd {
	if m.Choice == "doctor" {
		m.state = stateDoctor
		return runDoctor(m.tools)
	}
	return nil
}

type doctorResultsMsg []doctor.Result
type toolStatusesMsg []install.ToolStatus

func runDoctor(tools []string) tea.Cmd {
	return func() tea.Msg {
		configPath, _ := config.DefaultConfigPath()
		checks := doctor.DefaultChecks(tools, configPath)
		runner := doctor.NewRunner(checks)
		return doctorResultsMsg(runner.Run())
	}
}

func fetchToolStatuses() tea.Cmd {
	return func() tea.Msg {
		var statuses []install.ToolStatus
		for _, t := range install.Registry {
			statuses = append(statuses, install.GetStatus(t))
		}
		return toolStatusesMsg(statuses)
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case doctorResultsMsg:
		m.results = msg
		return m, nil

	case toolStatusesMsg:
		m.toolStatuses = msg
		m.isLoading = false
		return m, nil

	case tea.KeyMsg:
		// Back actions
		if m.state != stateMenu {
			if msg.String() == "q" || msg.String() == "esc" {
				m.state = stateMenu
				m.results = nil
				m.toolStatuses = nil
				m.cursor = 0
				return m, nil
			}
		}

		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "q", "esc":
			if m.state == stateMenu {
				return m, tea.Quit
			}
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			max := len(m.choices) - 1
			if m.state == stateInstallTools {
				max = len(m.toolStatuses)
			} else if m.state == stateDoctor {
				max = len(m.results) - 1
			}
			if m.cursor < max {
				m.cursor++
			}
		case "enter", " ":
			if m.state == stateDoctor && len(m.results) > 0 {
				action := m.results[m.cursor].RemedyAction
				if action != "" {
					m.Choice = action
					return m, tea.Quit
				}
				return m, nil
			}

			if m.state == stateMenu {
				rawChoice := m.choices[m.cursor]
				switch rawChoice {
				case "doctor":
					m.state = stateDoctor
					return m, runDoctor(m.tools)
				case "install tools":
					m.state = stateInstallTools
					m.isLoading = true
					m.cursor = 0
					return m, fetchToolStatuses()
				case "update fleet":
					m.Choice = "update-all"
					return m, tea.Quit
				case "exit":
					m.Choice = "exit"
					return m, tea.Quit
				}
			}

			if m.state == stateInstallTools {
				if m.cursor == len(m.toolStatuses) {
					// "Update All" option
					m.Choice = "update-all"
					return m, tea.Quit
				}
				// Selection of a specific tool
				selected := m.toolStatuses[m.cursor]
				if selected.Installed {
					m.Choice = "update " + selected.Tool.Name
				} else {
					m.Choice = "install " + selected.Tool.Name
				}
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m Model) View() string {
	switch m.state {
	case stateDoctor:
		return m.doctorView()
	case stateInstallTools:
		return m.installToolsView()
	default:
		return m.menuView()
	}
}

func (m Model) menuView() string {
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

func (m Model) installToolsView() string {
	var s strings.Builder
	s.WriteString(styles.TitleStyle.Render("--- 🛠  Installable Tools ---"))
	s.WriteString("\n\n")

	if m.isLoading {
		s.WriteString(styles.IrisStyle.Render("Checking registry..."))
		s.WriteString("\n")
		return styles.BoxStyle.Render(s.String())
	}

	for i, status := range m.toolStatuses {
		cursor := " "
		if m.cursor == i {
			cursor = styles.IrisStyle.Render(">")
		}

		stateIcon := styles.MutedStyle.Render("[ ]")
		if status.Installed {
			stateIcon = styles.SuccessStyle.Render("[✓]")
		}

		name := status.Tool.Name
		if m.cursor == i {
			name = styles.IrisStyle.Bold(true).Render(name)
		}

		s.WriteString(fmt.Sprintf("%s %s %-15s %s\n", cursor, stateIcon, name, styles.MutedStyle.Render(status.Tool.Description)))
		if status.Installed {
			s.WriteString(fmt.Sprintf("    %s %s\n", styles.IrisStyle.Render("Version:"), styles.MutedStyle.Render(status.Version)))
		}
		s.WriteString("\n")
	}

	// Update All Option
	cursor := " "
	label := "Update All Fleet"
	if m.cursor == len(m.toolStatuses) {
		cursor = styles.IrisStyle.Render(">")
		label = styles.IrisStyle.Bold(true).Render(label)
	}
	s.WriteString(fmt.Sprintf("\n%s %s 🔄\n", cursor, label))

	s.WriteString(styles.MutedStyle.Render("\nPress Enter to install/update, q to return."))

	return styles.BoxStyle.Render(s.String())
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

	for i, res := range m.results {
		statusIcon := styles.SuccessStyle.Render("✅")
		if res.Status == doctor.StatusFail {
			statusIcon = styles.ErrorStyle.Render("❌")
		} else if res.Status == doctor.StatusWarn {
			statusIcon = styles.WarnStyle.Render("⚠️")
		}

		cursor := " "
		name := styles.IrisStyle.Render(res.Name)
		if m.cursor == i {
			cursor = styles.IrisStyle.Render(">")
			name = styles.IrisStyle.Bold(true).Render(res.Name)
		}
		s.WriteString(fmt.Sprintf("%s %s %s\n", cursor, statusIcon, name))
		s.WriteString(fmt.Sprintf("   %s\n", styles.MutedStyle.Render(res.Detail)))
		if res.Remedy != "" {
			s.WriteString(fmt.Sprintf("   %s %s\n", styles.GoldStyle.Render("Remedy:"), res.Remedy))
			if m.cursor == i && res.RemedyAction != "" {
				s.WriteString(fmt.Sprintf("   %s\n", styles.MutedStyle.Render("Press Enter to apply this remedy.")))
			}
		}
		s.WriteString("\n")
	}

	s.WriteString(styles.MutedStyle.Render("Use ↑/↓ to inspect remedies; Enter applies an available remedy; q returns."))

	return styles.BoxStyle.Render(s.String())
}
