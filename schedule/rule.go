package schedule

import "time"

/// Only can repeat for one period now
type Rule struct {
	Duration time.Duration
}

func (r *Rule) GetNextTime(prevTime time.Time) time.Time {
	return prevTime.Add(r.Duration)
}
