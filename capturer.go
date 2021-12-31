package pprofs

import (
	"sync"
	"time"
)

type capturer struct {
	profiles []Profile
	trigger  Trigger
	storage  Storage
	logger   Logger
}

func newCapturer(opts ...Option) *capturer {
	c := defaultCapturer()
	for _, v := range opts {
		v(c)
	}
	return c
}

func (c *capturer) run() {
	if profiles := c.profiles; len(profiles) > 0 {
		for {
			if err := c.trigger.Wait(); err != nil {
				c.logger.Printf("wait: %v", err)
				continue
			}

			wg := &sync.WaitGroup{}
			wg.Add(len(profiles))

			now := time.Now()
			for _, profile := range c.profiles {
				go func(p Profile) {
					defer wg.Done()
					name := p.Name()
					w, err := c.storage.WriteCloser(name, now)
					if err != nil {
						c.logger.Printf("new writer for %v %v: %v", name, now, err)
						return
					}
					defer w.Close()
					if err := p.Capture(w); err != nil {
						c.logger.Printf("capture %v at %v: %v", name, now, err)
					}
				}(profile)
			}
			wg.Wait()
		}
	}
}
