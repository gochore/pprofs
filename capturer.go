package pprofs

import (
	"fmt"
	"io"
	"log"
	"sync"
	"time"
)

type capturer struct {
	profiles []profile
	trigger  Trigger
	storage  Storage
	logger   Logger
}

func defaultCapturer() *capturer {
	return &capturer{
		profiles: []profile{
			CpuProfile(),
			HeapProfile(),
		},
		trigger: NewRandomIntervalTrigger(15*time.Second, 2*time.Minute),
		storage: NewFileStorageFromEnv(),
		logger:  log.New(io.Discard, "", 0),
	}
}

func newCapturer(opts ...Option) (*capturer, error) {
	c := defaultCapturer()
	for _, v := range opts {
		v(c)
	}
	if err := c.validate(); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *capturer) run() {
	for {
		if err := c.trigger.Wait(); err != nil {
			c.logger.Printf("wait: %v", err)
			continue
		}

		wg := &sync.WaitGroup{}
		wg.Add(len(c.profiles))

		now := time.Now()
		for _, p := range c.profiles {
			go func(p profile) {
				defer wg.Done()
				name := p.name()
				w, err := c.storage.WriteCloser(name, now)
				if err != nil {
					c.logger.Printf("new writer for %v %v: %v", name, now, err)
					return
				}
				defer w.Close()
				if err := p.capture(w); err != nil {
					c.logger.Printf("capture %v at %v: %v", name, now, err)
				}
			}(p)
		}
		wg.Wait()
	}
}

func (c *capturer) validate() error {
	if len(c.profiles) == 0 {
		return fmt.Errorf("%w: empty profiles", ErrInvalidOption)
	}
	exists := map[string]struct{}{}
	for _, v := range c.profiles {
		if v == nil {
			return fmt.Errorf("%w: nil profile", ErrInvalidOption)
		}
		name := v.name()
		if _, ok := exists[name]; ok {
			return fmt.Errorf("%w: duplicated profile %v", ErrInvalidOption, name)
		}
		exists[v.name()] = struct{}{}
	}

	if c.trigger == nil {
		return fmt.Errorf("%w: nil trigger", ErrInvalidOption)
	}

	if c.storage == nil {
		return fmt.Errorf("%w: nil storage", ErrInvalidOption)
	}

	if c.logger == nil {
		return fmt.Errorf("%w: nil logger", ErrInvalidOption)
	}

	return nil
}
