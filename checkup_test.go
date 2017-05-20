package gopatrol

import (
	"errors"
	"testing"
	"time"
)

func TestCheckAndStore(t *testing.T) {
	f := new(fake)
	c := Checkup{
		Checkers:         []Checker{f, f},
		ConcurrentChecks: 1,
		Timestamp:        time.Now(),
		Notifier:         []Notifier{f},
	}

	_, err := c.Check()
	if err != nil {
		t.Errorf("Didn't expect an error: %v", err)
	}
	if got, want := f.checked, 2; got != want {
		t.Errorf("Expected %d checks to be executed, but had: %d", want, got)
	}

	if got, want := f.notified, 1; got != want {
		t.Errorf("Expected Notify() to be called %d time, called %d times", want, got)
	}

	// Check error handling
	f.returnErr = true
	_, err = c.Check()
	if err == nil {
		t.Error("Expected an error, didn't get one")
	}
	if got, want := err.Error(), "i'm an error; i'm an error"; got != want {
		t.Errorf(`Expected error string "%s" but got: "%s"`, want, got)
	}

	c.ConcurrentChecks = -1
	_, err = c.Check()
	if err == nil {
		t.Error("Expected an error with ConcurrentChecks < 0, didn't get one")
	}
	c.ConcurrentChecks = 0
	_, err = c.Check()
	if err == nil {
		t.Error("Expected an error with no storage, didn't get one")
	}
}

func TestComputeStats(t *testing.T) {
	s := Result{Times: []Attempt{
		{RTT: 7 * time.Second},
		{RTT: 4 * time.Second},
		{RTT: 4 * time.Second},
		{RTT: 6 * time.Second},
		{RTT: 6 * time.Second},
		{RTT: 3 * time.Second},
	}}.ComputeStats()

	if got, want := s.Total, 30*time.Second; got != want {
		t.Errorf("Expected Total=%v, got %v", want, got)
	}
	if got, want := s.Mean, 5*time.Second; got != want {
		t.Errorf("Expected Mean=%v, got %v", want, got)
	}
	if got, want := s.Median, 5*time.Second; got != want {
		t.Errorf("Expected Median=%v, got %v", want, got)
	}
	if got, want := s.Min, 3*time.Second; got != want {
		t.Errorf("Expected Min=%v, got %v", want, got)
	}
	if got, want := s.Max, 7*time.Second; got != want {
		t.Errorf("Expected Max=%v, got %v", want, got)
	}
}

func TestResultStatus(t *testing.T) {
	r := Result{Healthy: true}
	if got, want := r.Status(), Healthy; got != want {
		t.Errorf("Expected status '%s' but got: '%s'", want, got)
	}

	r = Result{Degraded: true}
	if got, want := r.Status(), Degraded; got != want {
		t.Errorf("Expected status '%s' but got: '%s'", want, got)
	}

	r = Result{Down: true}
	if got, want := r.Status(), Down; got != want {
		t.Errorf("Expected status '%s' but got: '%s'", want, got)
	}

	r = Result{}
	if got, want := r.Status(), Unknown; got != want {
		t.Errorf("Expected status '%s' but got: '%s'", want, got)
	}

	// These are invalid states, but we need to test anyway in case a
	// checker is buggy. We expect the worst of the enabled fields.
	r = Result{Down: true, Degraded: true}
	if got, want := r.Status(), Down; got != want {
		t.Errorf("(INVALID RESULT CASE) Expected status '%s' but got: '%s'", want, got)
	}
	r = Result{Degraded: true, Healthy: true}
	if got, want := r.Status(), Degraded; got != want {
		t.Errorf("(INVALID RESULT CASE) Expected status '%s' but got: '%s'", want, got)
	}
	r = Result{Down: true, Healthy: true}
	if got, want := r.Status(), Down; got != want {
		t.Errorf("(INVALID RESULT CASE) Expected status '%s' but got: '%s'", want, got)
	}
}

// func TestPriorityOver(t *testing.T) {
// 	for i, test := range []struct {
// 		status   string
// 		another  string
// 		expected bool
// 	}{
// 		{Down, Down, false},
// 		{Down, Degraded, true},
// 		{Down, Healthy, true},
// 		{Down, Unknown, true},
// 		{Degraded, Down, false},
// 		{Degraded, Degraded, false},
// 		{Degraded, Healthy, true},
// 		{Degraded, Unknown, true},
// 		{Healthy, Down, false},
// 		{Healthy, Degraded, false},
// 		{Healthy, Healthy, false},
// 		{Healthy, Unknown, true},
// 		{Unknown, Down, false},
// 		{Unknown, Degraded, false},
// 		{Unknown, Healthy, false},
// 		{Unknown, Unknown, false},
// 	} {
// 		actual := test.status.PriorityOver(test.another)
// 		if actual != test.expected {
// 			t.Errorf("Test %d: Expected %s.PriorityOver(%s)=%v, but got %v",
// 				i, test.status, test.another, test.expected, actual)
// 		}
// 	}
// }

var errTest = errors.New("i'm an error")

type fake struct {
	returnErr bool
	checked   int
	notified  int
}

func (f *fake) Check() (Result, error) {
	f.checked++
	r := Result{Timestamp: time.Now().UTC()}
	if f.returnErr {
		return r, errTest
	}
	return r, nil
}

func (f *fake) Notify(results Result) error {
	f.notified++
	return nil
}
