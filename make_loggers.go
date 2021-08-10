package logl

import (
	"io"
	"os"
)

func MakeLogger(level Level, w io.Writer) Logger {
	return NewHandleLogger(
		level,
		&StreamHandler{
			Output: w,
			Format: &FormatText{
				HasLevel:        true,
				HasDate:         true,
				HasTime:         true,
				HasMicroseconds: true,
				ShieldSpecial:   true,
			},
		},
	)
}

func MakeStdoutLogger(level Level) Logger {
	return MakeLogger(level, os.Stdout)
}
