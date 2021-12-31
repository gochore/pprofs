package pprofs

import (
	"fmt"
	"sync"
)

var (
	autoCapture struct {
		sync.Mutex
		c *capturer
	}
)

func EnableAutoCapture(options ...Option) error {
	autoCapture.Lock()
	defer autoCapture.Unlock()

	if autoCapture.c != nil {
		return fmt.Errorf("already enabled auto capture")
	}
	autoCapture.c = newCapturer(options...)

	go autoCapture.c.run()
	return nil
}
