package logl

import (
	"time"
)

type Record struct {
	Time    time.Time
	Level   Level
	Message string
}
