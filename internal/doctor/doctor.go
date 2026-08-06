package doctor

import (
	"sync"
)

type Status int

const (
	StatusPass Status = iota
	StatusWarn
	StatusFail
)

type Result struct {
	Name   string
	Status Status
	Detail string
	Remedy string
}

type Check interface {
	ID() string
	Name() string
	Run() Result
}

type Runner struct {
	checks []Check
}

func NewRunner(checks []Check) *Runner {
	return &Runner{checks: checks}
}

func (r *Runner) Run() []Result {
	var wg sync.WaitGroup
	results := make([]Result, len(r.checks))

	for i, check := range r.checks {
		wg.Add(1)
		go func(i int, c Check) {
			defer wg.Done()
			results[i] = c.Run()
		}(i, check)
	}

	wg.Wait()
	return results
}
