package pprofs

import (
	"sync"
	"time"
)

type Capturer struct {
	options *options
}

func NewCapturer(opts ...Option) *Capturer {
	opt := defaultOptions()
	for _, v := range opts {
		v(opt)
	}

	return &Capturer{
		options: opt,
	}
}

func (c *Capturer) run() {
	if profiles := c.options.profiles; len(profiles) > 0 {
		for {
			if err := c.options.trigger.Wait(); err != nil {
				c.options.logger.Printf("wait: %v", err)
				continue
			}

			wg := &sync.WaitGroup{}
			wg.Add(len(profiles))

			now := time.Now()
			for _, profile := range c.options.profiles {
				go func(p Profile) {
					defer wg.Done()
					name := p.Name()
					w, err := c.options.storage.WriteCloser(name, now)
					if err != nil {
						c.options.logger.Printf("new writer for %v %v: %v", name, now, err)
						return
					}
					defer w.Close()
					if err := p.Capture(w); err != nil {
						c.options.logger.Printf("capture %v at %v: %v", name, now, err)
					}
				}(profile)
			}
			wg.Wait()
		}
	}
}
