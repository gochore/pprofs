package pprofs

import (
	"math/rand"
	"time"
)

type Trigger interface {
	Wait() error
}

type IntervalTrigger struct {
	started  bool
	interval time.Duration
}

func NewIntervalTrigger(interval time.Duration) *IntervalTrigger {
	return &IntervalTrigger{
		interval: interval,
	}
}

func (t *IntervalTrigger) Wait() error {
	if !t.started {
		t.started = true
		return nil // trigger it at the beginning
	}

	time.Sleep(t.interval)
	return nil
}

type RandomIntervalTrigger struct {
	random   *rand.Rand
	started  bool
	min, max time.Duration
}

func NewRandomIntervalTrigger(min, max time.Duration) *RandomIntervalTrigger {
	return &RandomIntervalTrigger{
		random: rand.New(rand.NewSource(time.Now().UnixNano())),
		min:    min,
		max:    max,
	}
}

func (t *RandomIntervalTrigger) Wait() error {
	if !t.started {
		t.started = true
		return nil // trigger it at the beginning
	}

	time.Sleep(t.min + time.Duration(t.random.Int63n(int64(t.max-t.min))))
	return nil
}
