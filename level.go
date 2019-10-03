package logl

import (
	"fmt"

	"encoding"
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

var namesLevel = map[Level]string{
	LevelOff:      "off",
	LevelCritical: "critical",
	LevelError:    "error",
	LevelWarning:  "warning",
	LevelInfo:     "info",
	LevelDebug:    "debug",
	LevelTrace:    "trace",
}

var valuesLevel = map[string]Level{
	"off":      LevelOff,
	"critical": LevelCritical,
	"error":    LevelError,
	"warning":  LevelWarning,
	"info":     LevelInfo,
	"debug":    LevelDebug,
	"trace":    LevelTrace,
}

func (l Level) String() string {
	if name, ok := namesLevel[l]; ok {
		return name
	}
	return fmt.Sprintf("Level(%d)", l)
}

func ParseLevel(s string) (Level, error) {
	name := s
	level, ok := valuesLevel[name]
	if !ok {
		return 0, fmt.Errorf("logl.ParseLevel() invalid name %q", name)
	}
	return level, nil
}

func _() {
	var l Level
	var (
		_ encoding.TextMarshaler   = l
		_ encoding.TextUnmarshaler = &l
	)
}

var (
	_ encoding.TextMarshaler   = Level(0)
	_ encoding.TextUnmarshaler = (*Level)(nil)
)

func (l Level) MarshalText() (text []byte, err error) {
	value := l
	name, ok := namesLevel[value]
	if !ok {
		return nil, fmt.Errorf("logl.Level.MarshalText() invalid value %d", value)
	}
	return []byte(name), nil
}

func (l *Level) UnmarshalText(text []byte) error {
	name := string(text)
	value, ok := valuesLevel[name]
	if !ok {
		return fmt.Errorf("logl.Level.UnmarshalText() invalid name %q", name)
	}
	*l = value
	return nil
}
