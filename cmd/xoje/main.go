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
	"github.com/Xoje-Tech/xoje-environment/internal/verify"
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
	case "diagnose":
		fmt.Println("\n--- 🔍 Environment Diagnosis ---")
		results, allReady := verify.DiagnoseAll()
		for _, res := range results {
			if res.Available {
				fmt.Printf("  ✅ [%s] available: %s\n", res.Name, res.Path)
			} else {
				fmt.Printf("  ❌ [%s] NOT found\n", res.Name)
			}
		}
		if !allReady {
			fmt.Println("\n⚠️  Warning: Missing prerequisites. Some features may fail.")
		}
		fmt.Println("--------------------------------\n")

	case "install":
		fmt.Printf("\n--- 🛠  Installing: %s ---\n", tool)
		err := install.InstallTool(tool, "", false)
		if err != nil {
			return fmt.Errorf("installation failed: %w", err)
		}
		fmt.Printf("✅ %s successfully installed!\n", tool)

		// Update config
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
		fmt.Println("--------------------------------\n")

	case "tui":
		m := tui.NewModel()
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

		// Feedback visual de la elección
		fmt.Printf("\n🚀 Selected from TUI: %s\n", tm.Choice)
		
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
