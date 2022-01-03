package pprofs

type Option func(c *capturer)

func WithProfiles(profiles ...profile) Option {
	return func(c *capturer) {
		c.profiles = profiles
	}
}

func WithTrigger(trigger Trigger) Option {
	return func(c *capturer) {
		c.trigger = trigger
	}
}

func WithStorage(storage Storage) Option {
	return func(c *capturer) {
		c.storage = storage
	}
}

func WithLogger(logger Logger) Option {
	return func(c *capturer) {
		c.logger = logger
	}
}
