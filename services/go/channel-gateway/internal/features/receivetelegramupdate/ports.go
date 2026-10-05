package receivetelegramupdate

import "time"

// Clock is a dependency only this slice needs; injecting it makes tests deterministic.
type Clock interface {
	Now() time.Time
}
