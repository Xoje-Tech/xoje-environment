package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/Xoje-Tech/xoje-environment/internal/cli"
	"github.com/Xoje-Tech/xoje-environment/internal/config"
	"github.com/Xoje-Tech/xoje-environment/internal/install"
	"github.com/Xoje-Tech/xoje-environment/internal/tui"
	"github.com/Xoje-Tech/xoje-environment/internal/tui/styles"
)

func main() {
	configPath, err := config.DefaultConfigPath()
	if err != nil {
		fmt.Printf("❌ Error resolving config path: %v\n", err)
		os.Exit(1)
	}

	if err := Run(os.Args[1:], configPath); err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		os.Exit(1)
	}
}

func Run(args []string, configPath string) error {
	cfg, err := config.LoadOrCreate(configPath)
	if err != nil {
		return fmt.Errorf("initializing configuration: %w", err)
	}

	cmd, tool, err := cli.Parse(args)
	if err != nil {
		return fmt.Errorf("arguments error: %w", err)
	}

	return dispatch(cmd, tool, cfg)
}

func dispatch(cmd string, tool string, cfg *config.Config) error {
	switch cmd {
	case "doctor":
		m := tui.NewModel(cfg.InstalledTools)
		m.Choice = "doctor"
		p := tea.NewProgram(m)
		if _, err := p.Run(); err != nil {
			return fmt.Errorf("doctor command error: %w", err)
		}
		return nil

	case "install":
		path, exists := install.IsInstalled(tool)
		if exists {
			fmt.Printf("%s %s is already in PATH at: %s\n", styles.SuccessStyle.Render("✓"), tool, styles.MutedStyle.Render(path))
			registerTool(tool, cfg)
			return nil
		}

		fmt.Printf("\n%s %s\n", styles.TitleStyle.Render("🛠  Installing:"), styles.IrisStyle.Render(tool))
		err := install.InstallTool(tool, "", false)
		if err != nil {
			return fmt.Errorf("installation failed: %w", err)
		}
		registerTool(tool, cfg)
		return nil

	case "update":
		if tool != "" {
			fmt.Printf("\n%s %s\n", styles.TitleStyle.Render("🔄 Updating:"), styles.IrisStyle.Render(tool))
			if err := install.InstallTool(tool, "", true); err != nil {
				return fmt.Errorf("failed to update %s: %w", tool, err)
			}
			registerTool(tool, cfg)
			return nil
		}

		fmt.Printf("\n%s\n", styles.TitleStyle.Render("🔄 Updating Fleet"))
		if len(cfg.InstalledTools) == 0 {
			fmt.Println(styles.MutedStyle.Render("No tools registered in configuration."))
		}
		for _, t := range cfg.InstalledTools {
			if err := install.InstallTool(t, "", true); err != nil {
				fmt.Printf("  %s Failed to update %s: %v\n", styles.ErrorStyle.Render("❌"), t, err)
			}
		}

	case "tui":
		m := tui.NewModel(cfg.InstalledTools)
		p := tea.NewProgram(m)
		finalModel, err := p.Run()
		if err != nil {
			return fmt.Errorf("TUI Error: %w", err)
		}

		tm := finalModel.(tui.Model)
		if tm.Choice == "" || tm.Choice == "exit" {
			fmt.Println("\nBye! 👋")
			return nil
		}

		if strings.HasPrefix(tm.Choice, "install ") {
			toolName := strings.TrimPrefix(tm.Choice, "install ")
			return dispatch("install", toolName, cfg)
		}
		return dispatch(tm.Choice, "", cfg)

	default:
		return fmt.Errorf("unrecognized command: %s", cmd)
	}

	return nil
}

func registerTool(tool string, cfg *config.Config) {
	alreadyInstalled := false
	for _, t := range cfg.InstalledTools {
		if t == tool {
			alreadyInstalled = true
			break
		}
	}
	if !alreadyInstalled {
		cfg.InstalledTools = append(cfg.InstalledTools, tool)
		if err := cfg.Save(); err != nil {
			fmt.Printf("⚠️  Warning: failed to update config registry: %v\n", err)
		}
	}
}
