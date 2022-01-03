package pprofs

import (
	"io"
	"log"
	"time"
)

func defaultCapturer() *capturer {
	return &capturer{
		profiles: []profile{
			CpuProfile{Duration: 15 * time.Second},
			heapProfile{},
		},
		trigger: NewRandomIntervalTrigger(15*time.Second, 2*time.Minute),
		storage: NewFileStorageFromEnv(),
		logger:  log.New(io.Discard, "", 0),
	}
}

type Option func(opt *capturer)

func WithProfiles(profiles ...profile) Option {
	return func(c *capturer) {
		c.profiles = profiles
	}
}

func WithTrigger(trigger Trigger) Option {
	return func(opt *capturer) {
		opt.trigger = trigger
	}
}

func WithStorage(storage Storage) Option {
	return func(opt *capturer) {
		opt.storage = storage
	}
}

func WithLogger(logger Logger) Option {
	return func(opt *capturer) {
		opt.logger = logger
	}
}
