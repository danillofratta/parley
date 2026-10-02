package receivetelegramupdate

import "time"

type Clock interface {
	Now() time.Time
}
