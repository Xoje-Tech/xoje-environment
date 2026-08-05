package cli

import "fmt"

func Parse(args []string) (string, string, error) {
	if len(args) == 0 {
		return "tui", "", nil
	}

	subcommand := args[0]
	switch subcommand {
	case "diagnose", "bootstrap":
		return "diagnose", "", nil
	case "install":
		if len(args) < 2 {
			return "", "", fmt.Errorf("usage: xoje install [tool_name]")
		}
		return "install", args[1], nil
	case "tui":
		return "tui", "", nil
	default:
		return "", "", fmt.Errorf("unknown subcommand: %s", subcommand)
	}
}
