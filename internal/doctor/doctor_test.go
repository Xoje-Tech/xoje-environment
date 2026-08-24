package doctor

import (
	"testing"
	"time"
)

type MockCheck struct {
	id     string
	name   string
	result Result
	delay  time.Duration
}

func (m *MockCheck) ID() string   { return m.id }
func (m *MockCheck) Name() string { return m.name }
func (m *MockCheck) Run() Result {
	if m.delay > 0 {
		time.Sleep(m.delay)
	}
	return m.result
}

func TestRunner_Run(t *testing.T) {
	checks := []Check{
		&MockCheck{id: "c1", name: "Check 1", result: Result{Status: StatusPass}},
		&MockCheck{id: "c2", name: "Check 2", result: Result{Status: StatusFail}},
	}

	runner := NewRunner(checks)
	results := runner.Run()

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	// Verify concurrent execution (should be faster than serial)
	slowChecks := []Check{
		&MockCheck{id: "s1", delay: 100 * time.Millisecond, result: Result{Status: StatusPass}},
		&MockCheck{id: "s2", delay: 100 * time.Millisecond, result: Result{Status: StatusPass}},
	}

	start := time.Now()
	NewRunner(slowChecks).Run()
	elapsed := time.Since(start)

	if elapsed >= 200*time.Millisecond {
		t.Errorf("execution seems sequential, took %v", elapsed)
	}
}

func TestRunnerPerCheckTimeout(t *testing.T) {
	runner := NewRunner([]Check{
		&MockCheck{id: "slow", name: "Slow check", delay: 100 * time.Millisecond, result: Result{Status: StatusPass}},
	})
	runner.timeout = 10 * time.Millisecond

	start := time.Now()
	results := runner.Run()
	if elapsed := time.Since(start); elapsed >= 80*time.Millisecond {
		t.Fatalf("per-check timeout did not bound execution: %v", elapsed)
	}
	if len(results) != 1 || results[0].Status != StatusFail || results[0].Remedy == "" {
		t.Fatalf("timeout should produce actionable failure: %#v", results)
	}
}

func TestAggregateReadiness(t *testing.T) {
	for name, tc := range map[string]struct {
		results []Result
		want    Readiness
	}{
		"ready":     {[]Result{{Status: StatusPass}}, ReadinessReady},
		"degraded":  {[]Result{{Status: StatusPass}, {Status: StatusWarn}}, ReadinessDegraded},
		"not ready": {[]Result{{Status: StatusWarn}, {Status: StatusFail}}, ReadinessNotReady},
	} {
		t.Run(name, func(t *testing.T) {
			if got := AggregateReadiness(tc.results); got != tc.want {
				t.Fatalf("AggregateReadiness() = %q, want %q", got, tc.want)
			}
		})
	}
}
