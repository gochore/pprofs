package pprofs

import (
	"math/rand"
	"time"
)

// Trigger is an abstraction of a trigger.
type Trigger interface {
	// Wait should block until need trigger capturing profile,
	// returning error means skipping.
	Wait() error
}

// IntervalTrigger will trigger periodically with constant interval.
type IntervalTrigger struct {
	started  bool
	interval time.Duration
}

// NewIntervalTrigger returns a IntervalTrigger.
func NewIntervalTrigger(interval time.Duration) *IntervalTrigger {
	return &IntervalTrigger{
		interval: interval,
	}
}

// Wait implements Trigger.
func (t *IntervalTrigger) Wait() error {
	if !t.started {
		t.started = true
		return nil // trigger it at the beginning
	}

	time.Sleep(t.interval)
	return nil
}

// RandomIntervalTrigger will trigger periodically with random interval.
type RandomIntervalTrigger struct {
	random   *rand.Rand
	started  bool
	min, max time.Duration
}

// NewRandomIntervalTrigger returns a IntervalTrigger.
func NewRandomIntervalTrigger(min, max time.Duration) *RandomIntervalTrigger {
	return &RandomIntervalTrigger{
		random: rand.New(rand.NewSource(time.Now().UnixNano())),
		min:    min,
		max:    max,
	}
}

// Wait implements Trigger.
func (t *RandomIntervalTrigger) Wait() error {
	if !t.started {
		t.started = true
		return nil // trigger it at the beginning
	}

	time.Sleep(t.min + time.Duration(t.random.Int63n(int64(t.max-t.min))))
	return nil
}
