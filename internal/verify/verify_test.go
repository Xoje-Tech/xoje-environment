package verify

import "testing"

func TestVerifyExecutable(t *testing.T) {
	// Test a command we know exists (go is installed via mise and in PATH)
	res := VerifyExecutable("go")
	if !res.Available {
		t.Errorf("expected 'go' to be available, got false")
	}

	// Test a command we know doesn't exist
	res2 := VerifyExecutable("nonexistent-command-xyz")
	if res2.Available {
		t.Errorf("expected 'nonexistent-command-xyz' to be unavailable, got true")
	}
}

func TestDiagnoseAll(t *testing.T) {
	results, allReady := DiagnoseAll()
	
	// We expect 3 results: go, node, git
	if len(results) != 3 {
		t.Errorf("expected 3 results, got %d", len(results))
	}
	
	// Check if names match
	expectedNames := map[string]bool{"go": true, "node": true, "git": true}
	for _, res := range results {
		if !expectedNames[res.Name] {
			t.Errorf("unexpected command name in results: %s", res.Name)
		}
	}
	
	// allReady should reflect the aggregate state
	countAvailable := 0
	for _, res := range results {
		if res.Available {
			countAvailable++
		}
	}
	
	if allReady && countAvailable < 3 {
		t.Errorf("allReady is true but only %d/3 commands available", countAvailable)
	}
}
