package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Xoje-Tech/xoje-environment/internal/cli"
	"github.com/Xoje-Tech/xoje-environment/internal/config"
	"github.com/Xoje-Tech/xoje-environment/internal/doctor"
	"github.com/Xoje-Tech/xoje-environment/internal/install"
	"github.com/Xoje-Tech/xoje-environment/internal/tui"
	"github.com/Xoje-Tech/xoje-environment/internal/tui/styles"
	tea "github.com/charmbracelet/bubbletea"
)

const Version = "1.1.0"

func main() {
	configPath, err := config.DefaultConfigPath()
	if err != nil {
		fmt.Printf("Error resolving config path: %v\n", err)
		os.Exit(1)
	}

	if err := Run(os.Args[1:], configPath); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

// Run executes one CLI invocation using the supplied configuration path.
func Run(args []string, configPath string) error {
	cfg, err := config.LoadOrCreate(configPath)
	if err != nil {
		return fmt.Errorf("initializing configuration: %w", err)
	}

	cmd, toolName, err := cli.Parse(args)
	if err != nil {
		return fmt.Errorf("arguments error: %w", err)
	}

	isTUI := len(args) == 0

	for {
		err := handleCommand(cmd, toolName, cfg, isTUI)
		if err != nil {
			return err
		}

		if !isTUI {
			break
		}

		fmt.Println("\nPress Enter to return to menu...")
		fmt.Scanln()
		cmd = "tui"
		toolName = ""
	}

	return nil
}

func handleCommand(cmd, toolName string, cfg *config.Config, isTUI bool) error {
	switch cmd {
	case "version":
		fmt.Printf("xoje version %s\n", Version)
		return nil

	case "doctor":
		configPath, _ := config.DefaultConfigPath()
		checks := doctor.DefaultChecks(cfg.InstalledTools, configPath)
		runner := doctor.NewRunner(checks)
		results := runner.Run()

		if !isTUI {
			fmt.Print(formatDoctorReport(results))
			return nil
		}

		m := tui.NewModel(cfg.InstalledTools)
		m.Choice = "doctor"
		p := tea.NewProgram(m)
		if res, err := p.Run(); err != nil {
			return fmt.Errorf("doctor command error: %v", err)
		} else {
			finalModel := res.(tui.Model)
			if finalModel.Choice != "" {
				return dispatchTUIAction(finalModel.Choice, cfg)
			}
		}
		return nil

	case "install":
		return install.Install(toolName, cfg)

	case "update":
		return install.Update(toolName, cfg)

	case "update-all":
		return install.UpdateAll(cfg)

	case "tui":
		m := tui.NewModel(cfg.InstalledTools)
		p := tea.NewProgram(m)
		if res, err := p.Run(); err != nil {
			return fmt.Errorf("TUI error: %v", err)
		} else {
			finalModel := res.(tui.Model)
			if finalModel.Choice != "" {
				return dispatchTUIAction(finalModel.Choice, cfg)
			}
		}
		return nil

	default:
		return fmt.Errorf("unknown command: %s", cmd)
	}
}

func formatDoctorReport(results []doctor.Result) string {
	var report strings.Builder
	report.WriteString(styles.TitleStyle.Render("--- 🩺 System Doctor ---"))
	report.WriteString("\n\n")

	for _, result := range results {
		statusIcon := styles.SuccessStyle.Render("✅")
		switch result.Status {
		case doctor.StatusWarn:
			statusIcon = styles.WarnStyle.Render("⚠️")
		case doctor.StatusFail:
			statusIcon = styles.ErrorStyle.Render("❌")
		}

		report.WriteString(fmt.Sprintf("%s %s\n", statusIcon, styles.IrisStyle.Render(result.Name)))
		report.WriteString(fmt.Sprintf("  %s\n", styles.MutedStyle.Render(result.Detail)))
		if result.Remedy != "" {
			report.WriteString(fmt.Sprintf("  %s %s\n", styles.GoldStyle.Render("Remedy:"), result.Remedy))
		}
		report.WriteString("\n")
	}

	readiness := doctor.AggregateReadiness(results)
	readinessText := fmt.Sprintf("READINESS: %s", readiness)
	switch readiness {
	case doctor.ReadinessDegraded:
		readinessText = styles.WarnStyle.Render(readinessText)
	case doctor.ReadinessNotReady:
		readinessText = styles.ErrorStyle.Render(readinessText)
	default:
		readinessText = styles.SuccessStyle.Render(readinessText)
	}
	report.WriteString(readinessText)

	return styles.BoxStyle.Render(report.String()) + "\n"
}

func dispatchTUIAction(choice string, cfg *config.Config) error {
	parts := strings.Split(choice, " ")
	cmd := parts[0]
	tool := ""
	if len(parts) > 1 {
		tool = parts[1]
	}

	switch cmd {
	case "install":
		return install.Install(tool, cfg)
	case "update":
		return install.Update(tool, cfg)
	case "update-all":
		return install.UpdateAll(cfg)
	case "exit":
		os.Exit(0)
	}
	return nil
}
