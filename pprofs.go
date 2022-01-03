package pprofs

import (
	"errors"
	"sync"
)

var (
	ErrAlreadyEnabled = errors.New("already enabled")
	ErrInvalidOption  = errors.New("invalid option")
)

var (
	enabledCapturer struct {
		sync.Mutex
		*capturer
	}
)

func EnableCapture(options ...Option) error {
	enabledCapturer.Lock()
	defer enabledCapturer.Unlock()

	if enabledCapturer.capturer != nil {
		return ErrAlreadyEnabled
	}
	c, err := newCapturer(options...)
	if err != nil {
		return err
	}
	enabledCapturer.capturer = c

	go enabledCapturer.capturer.run()
	return nil
}
