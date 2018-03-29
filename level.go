package logl

import (
	"fmt"
	"strconv"
)

type Level int

const (
	LevelOff      Level = iota // log: [ ]
	LevelCritical              // log: [ Critical ]
	LevelError                 // log: [ Critical, Error ]
	LevelWarning               // log: [ Critical, Error, Warning ]
	LevelInfo                  // log: [ Critical, Error, Warning, Info ]
	LevelDebug                 // log: [ Critical, Error, Warning, Info, Debug ]
	LevelTrace                 // log: [ Critical, Error, Warning, Info, Debug, Trace ]
)

var name_Level = map[Level]string{
	LevelOff:      "off",
	LevelCritical: "critical",
	LevelError:    "error",
	LevelWarning:  "warning",
	LevelInfo:     "info",
	LevelDebug:    "debug",
	LevelTrace:    "trace",
}

var value_Level = map[string]Level{
	"off":      LevelOff,
	"critical": LevelCritical,
	"error":    LevelError,
	"warning":  LevelWarning,
	"info":     LevelInfo,
	"debug":    LevelDebug,
	"trace":    LevelTrace,
}

func (l Level) String() string {
	if name, ok := name_Level[l]; ok {
		return name
	}
	return strconv.Itoa(int(l))
}

func ParseLevel(s string) (Level, error) {
	level, ok := value_Level[s]
	if ok {
		return level, nil
	}
	return level, fmt.Errorf("invalid log level: %s", s)
}
