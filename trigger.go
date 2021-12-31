package pprofs

import (
	"math/rand"
	"os"
	"time"

	"github.com/shirou/gopsutil/v3/process"
)

type Trigger interface {
	Wait() error
}

type RandomIntervalTrigger struct {
	random   *rand.Rand
	started  bool
	min, max time.Duration
}

func NewRandomIntervalTrigger(min time.Duration, max time.Duration) *RandomIntervalTrigger {
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

	d := t.min + time.Duration(t.random.Int63n(int64(t.max-t.min)))
	time.Sleep(d)
	return nil
}

type PsTrigger struct {
	interval time.Duration
	started  bool
	cpu      float64
	mem      float64
}

func NewPsTrigger(interval time.Duration, cpu float64, mem float64) *PsTrigger {
	return &PsTrigger{
		interval: interval,
		cpu:      cpu,
		mem:      mem,
	}
}

func (t *PsTrigger) Wait() error {
	if !(t.cpu > 0 || t.mem > 0) {
		select {} // block forever
	}

	p, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		return err
	}

	if !t.started {
		t.started = true
	} else {
		time.Sleep(t.interval)
	}

	for {
		if t.cpu > 0 {
			if cpu, err := p.CPUPercent(); err != nil {
				return err
			} else if cpu > t.cpu {
				return nil
			}
		}

		if t.mem > 0 {
			if mem, err := p.MemoryPercent(); err != nil {
				return err
			} else if float64(mem) > t.mem {
				return nil
			}
		}
		time.Sleep(t.interval)
	}
}
