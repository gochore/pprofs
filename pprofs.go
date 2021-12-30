package pprofs

import (
	"fmt"
	"sync"
)

var (
	autoCapture struct {
		sync.Mutex
		Enabled bool
	}
)

func EnableAutoCapture(capturer *Capturer) error {
	if capturer == nil {
		return fmt.Errorf("nil capturer")
	}

	autoCapture.Lock()
	defer autoCapture.Unlock()

	if autoCapture.Enabled {
		return fmt.Errorf("already enabled auto capture")
	}
	autoCapture.Enabled = true

	go capturer.run()
	return nil
}
