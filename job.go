package mmbot

// JobFunc is job action function.
type JobFunc func(*Robot) error

// Job is a scheduled task.
type Job struct {
	// Schedule pattern.
	// See https://godoc.org/github.com/robfig/cron

	// NOTE: It is different from cron, there is also seconds field.
	Schedule string

	// Job function.
	Action JobFunc
}
