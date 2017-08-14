package logl

import (
	"io"
)

type Config struct {
	Output  WriteFlusher
	Prefix  string
	Level   Level
	Flag    int
	NotSafe bool // for multiple goroutine not safe
}

// Flags
const (
	Ldate         = 1 << iota // the date in the local time zone: 2009/01/23
	Ltime                     // the time in the local time zone: 01:23:23
	Lmicroseconds             // microsecond resolution: 01:23:23.123123.  assumes Ltime.
)

type WriteFlusher interface {
	io.Writer
	Flush() error
}

func OnlyWriter(w io.Writer) WriteFlusher {
	return fakeFlusher{w}
}

type fakeFlusher struct {
	io.Writer
}

func (fakeFlusher) Flush() error {
	return nil
}

var _ WriteFlusher = fakeFlusher{}
