package clock

import "time"

// System is the real clock; it satisfies the Clock port of the slices.
type System struct{}

func (System) Now() time.Time { return time.Now() }
