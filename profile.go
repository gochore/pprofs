package pprofs

import (
	"io"
	"runtime"
	"runtime/pprof"
	"time"
)

// CpuProfile captures the CPU profile.
func CpuProfile() cpuProfile {
	return cpuProfile{
		duration: 15 * time.Second,
	}
}

// WithDuration specifies the duration of the CPU profile.
func (p cpuProfile) WithDuration(d time.Duration) cpuProfile {
	p.duration = d
	return p
}

// HeapProfile captures the heap profile.
func HeapProfile() heapProfile {
	return heapProfile{}
}

// MutexProfile captures the mutex profile.
func MutexProfile() mutexProfile {
	return mutexProfile{}
}

// BlockProfile captures the block profile.
func BlockProfile() blockProfile {
	return blockProfile{
		rate: 1,
	}
}

// WithRate specifies the rate of the block profile.
func (p blockProfile) WithRate(rate int) blockProfile {
	p.rate = rate
	return p
}

// GoroutineProfile captures the goroutine profile.
func GoroutineProfile() goroutineProfile {
	return goroutineProfile{}
}

// ThreadcreateProfile captures the threadcreate profile.
func ThreadcreateProfile() threadcreateProfile {
	return threadcreateProfile{}
}

type profile interface {
	name() string
	capture(io.Writer) error
}

type cpuProfile struct {
	duration time.Duration
}

func (p cpuProfile) name() string {
	return "cpu"
}

func (p cpuProfile) capture(w io.Writer) error {
	dur := p.duration
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

type heapProfile struct{}

func (p heapProfile) name() string {
	return "heap"
}

func (p heapProfile) capture(w io.Writer) error {
	return captureProfile(w, p.name())
}

type mutexProfile struct{}

func (p mutexProfile) name() string {
	return "mutex"
}

func (p mutexProfile) capture(w io.Writer) error {
	return captureProfile(w, p.name())
}

type blockProfile struct {
	rate int
}

func (p blockProfile) name() string {
	return "block"
}

func (p blockProfile) capture(w io.Writer) error {
	runtime.SetBlockProfileRate(p.rate)
	return captureProfile(w, p.name())
}

type goroutineProfile struct{}

func (p goroutineProfile) name() string {
	return "goroutine"
}

func (p goroutineProfile) capture(w io.Writer) error {
	return captureProfile(w, p.name())
}

type threadcreateProfile struct{}

func (p threadcreateProfile) name() string {
	return "threadcreate"
}

func (p threadcreateProfile) capture(w io.Writer) error {
	return captureProfile(w, p.name())
}

func captureProfile(w io.Writer, name string) error {
	return pprof.Lookup(name).WriteTo(w, 0)
}
