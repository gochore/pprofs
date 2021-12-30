package pprofs

import (
	"time"
)

type options struct {
	profiles []Profile
	trigger  Trigger
	storage  Storage
	logger   Logger
}

func defaultOptions() *options {
	return &options{
		profiles: []Profile{
			CpuProfile{},
			HeapProfile{},
		},
		trigger: NewRandomIntervalTrigger(15*time.Second, 2*time.Minute),
		storage: NewFileStorageFromEnv(),
	}
}

type Option func(opt *options)

func WithProfiles(profiles ...Profile) Option {
	return func(opt *options) {
		opt.profiles = profiles
	}
}

func WithTrigger(trigger Trigger) Option {
	return func(opt *options) {
		opt.trigger = trigger
	}
}

func WithStorage(storage Storage) Option {
	return func(opt *options) {
		opt.storage = storage
	}
}

func WithLogger(logger Logger) Option {
	return func(opt *options) {
		opt.logger = logger
	}
}
