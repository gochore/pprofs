package pprofs

type Option func(c *capturer)

// WithProfiles specifies the pprof profiles.
func WithProfiles(profiles ...profile) Option {
	return func(c *capturer) {
		c.profiles = profiles
	}
}

// WithTrigger specifies the trigger.
func WithTrigger(trigger Trigger) Option {
	return func(c *capturer) {
		c.trigger = trigger
	}
}

// WithStorage specifies the storage.
func WithStorage(storage Storage) Option {
	return func(c *capturer) {
		c.storage = storage
	}
}

// WithLogger specifies the logger.
func WithLogger(logger Logger) Option {
	return func(c *capturer) {
		c.logger = logger
	}
}
