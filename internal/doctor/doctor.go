package doctor

import (
	"sync"
	"time"
)

type Status int

const (
	StatusPass Status = iota
	StatusWarn
	StatusFail
)

type Result struct {
	Name         string
	Status       Status
	Detail       string
	Remedy       string
	RemedyAction string
}

type Check interface {
	ID() string
	Name() string
	Run() Result
}

type Runner struct {
	checks  []Check
	timeout time.Duration
}

func NewRunner(checks []Check) *Runner {
	return &Runner{checks: checks, timeout: 3 * time.Second}
}

func (r *Runner) Run() []Result {
	var wg sync.WaitGroup
	results := make([]Result, len(r.checks))

	for i, check := range r.checks {
		wg.Add(1)
		go func(i int, c Check) {
			defer wg.Done()
			result := make(chan Result, 1)
			go func() { result <- c.Run() }()

			select {
			case results[i] = <-result:
			case <-time.After(r.timeout):
				results[i] = Result{
					Name:   c.Name(),
					Status: StatusFail,
					Detail: "Check timed out.",
					Remedy: "Retry the doctor check and inspect the underlying service if it remains slow.",
				}
			}
		}(i, check)
	}

	wg.Wait()
	return results
}

type Readiness string

const (
	ReadinessReady    Readiness = "READY"
	ReadinessDegraded Readiness = "DEGRADED"
	ReadinessNotReady Readiness = "NOT READY"
)

// AggregateReadiness reduces check results to the overall health state.
func AggregateReadiness(results []Result) Readiness {
	readiness := ReadinessReady
	for _, result := range results {
		switch result.Status {
		case StatusFail:
			return ReadinessNotReady
		case StatusWarn:
			readiness = ReadinessDegraded
		}
	}
	return readiness
}
