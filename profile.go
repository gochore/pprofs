package pprofs

import (
	"io"
	"runtime/pprof"
	"time"
)

type Profile interface {
	Name() string
	Capture(w io.Writer) (err error)
}

type CpuProfile struct {
	Duration time.Duration // 15 seconds by default
}

func (p CpuProfile) Name() string {
	return "cpu"
}

func (p CpuProfile) Capture(w io.Writer) error {
	dur := p.Duration
	if dur <= 0 {
		dur = 15 * time.Second
	}

	if err := pprof.StartCPUProfile(w); err != nil {
		return err
	}
	time.Sleep(dur)
	pprof.StopCPUProfile()
	return nil
}

type HeapProfile struct{}

func (p HeapProfile) Name() string {
	return "heap"
}

func (p HeapProfile) Capture(w io.Writer) error {
	return captureProfile(w, p.Name())
}

func captureProfile(w io.Writer, name string) error {
	return pprof.Lookup(name).WriteTo(w, 0)
}

// TODO support more profiles
