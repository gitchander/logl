package logl

import (
	"io"
	"os"
)

func MakeLogger(w io.Writer) Logger {
	f := FormatText{
		HasLevel:        true,
		HasDate:         true,
		HasTime:         true,
		HasMicroseconds: true,
		ShieldSpecial:   true,
	}.Formatter()
	return NewLoggerRW(
		&StreamRecordWriter{
			Writer:    w,
			Formatter: f,
		},
	)
}

func MakeLoggerStdout() Logger {
	return MakeLogger(os.Stdout)
}
