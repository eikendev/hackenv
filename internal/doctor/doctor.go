// Package doctor runs read-only checks of the host prerequisites.
package doctor

import (
	"fmt"
)

// Status is the outcome of a single check.
type Status int

const (
	// Pass means the prerequisite is satisfied.
	Pass Status = iota
	// Fail means the prerequisite is missing or broken.
	Fail
	// Skip means the check does not apply on this host.
	Skip
)

func (s Status) String() string {
	switch s {
	case Pass:
		return "PASS"
	case Fail:
		return "FAIL"
	case Skip:
		return "SKIP"
	default:
		return fmt.Sprintf("Status(%d)", int(s))
	}
}

// Result is the outcome of a check together with details for the user.
type Result struct {
	Name   string
	Status Status
	Detail string
	Hint   string
}

// Check is a single named prerequisite check.
type Check struct {
	Name string
	Run  func() Result
}

func pass(name, detail string) Result {
	return Result{Name: name, Status: Pass, Detail: detail}
}

func fail(name, detail, hint string) Result {
	return Result{Name: name, Status: Fail, Detail: detail, Hint: hint}
}

func skip(name, detail string) Result {
	return Result{Name: name, Status: Skip, Detail: detail}
}

// RunAll executes every check and returns the results in order.
func RunAll(checks []Check) []Result {
	results := make([]Result, 0, len(checks))
	for _, check := range checks {
		results = append(results, check.Run())
	}
	return results
}

// HasFailures reports whether any result failed.
func HasFailures(results []Result) bool {
	for _, result := range results {
		if result.Status == Fail {
			return true
		}
	}
	return false
}
