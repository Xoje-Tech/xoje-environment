package cli

import (
	"fmt"
	"strings"
)

// Parse takes command line arguments and returns (command, tool, error)
func Parse(args []string) (string, string, error) {
	if len(args) == 0 {
		return "tui", "", nil
	}

	cmd := strings.ToLower(args[0])

	switch cmd {
	case "doctor", "diagnose", "bootstrap":
		tool := ""
		if cmd == "doctor" && len(args) >= 2 && (args[1] == "--plain" || args[1] == "-p") {
			tool = "--plain"
		}
		return "doctor", tool, nil
	case "install":
		if len(args) < 2 {
			return "", "", fmt.Errorf("usage: xoje install [tool_name]")
		}
		return "install", args[1], nil
	case "update":
		tool := ""
		if len(args) >= 2 {
			tool = args[1]
		}
		return "update", tool, nil
	case "update-all":
		return "update-all", "", nil
	case "tui", "help", "--help", "-h":
		return "tui", "", nil
	case "version", "--version", "-v":
		return "version", "", nil
	default:
		return "", "", fmt.Errorf("unknown subcommand: %s", cmd)
	}
}
