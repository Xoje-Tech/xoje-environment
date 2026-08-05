package verify

import "os/exec"

type DiagnosticResult struct {
	Name      string `json:"name"`
	Available bool   `json:"available"`
	Path      string `json:"path"`
}

func VerifyExecutable(name string) DiagnosticResult {
	path, err := exec.LookPath(name)
	if err != nil {
		return DiagnosticResult{
			Name:      name,
			Available: false,
			Path:      "",
		}
	}
	return DiagnosticResult{
		Name:      name,
		Available: true,
		Path:      path,
	}
}

func DiagnoseAll() ([]DiagnosticResult, bool) {
	commands := []string{"go", "node", "git"}
	results := make([]DiagnosticResult, 0, len(commands))
	allReady := true

	for _, cmd := range commands {
		res := VerifyExecutable(cmd)
		results = append(results, res)
		if !res.Available {
			allReady = false
		}
	}

	return results, allReady
}
