package pprofs

// Logger is an abstraction of log output.
type Logger interface {
	Printf(format string, v ...interface{})
}
