package pprofs

import (
	"fmt"
	"sync"
)

var (
	enabledCapturer struct {
		sync.Mutex
		c *capturer
	}
)

func EnableCapture(options ...Option) error {
	enabledCapturer.Lock()
	defer enabledCapturer.Unlock()

	if enabledCapturer.c != nil {
		return fmt.Errorf("already enabled capture")
	}
	enabledCapturer.c = newCapturer(options...)

	go enabledCapturer.c.run()
	return nil
}
