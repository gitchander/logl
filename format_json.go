package logl

import (
	"bytes"
	"time"
)

const (
	// timeFormat = "2006-01-02 15:04:05"
	// timeFormat = time.RFC3339
	timeFormat = time.RFC3339Nano
)

type formatJSON struct {
	b bytes.Buffer
}

func (f *formatJSON) Format(r *Record) []byte {
	b := &(f.b)
	b.Reset()
	b.WriteByte('{')
	jsonEncodePair(b, "time", r.Time.Format(timeFormat))
	b.WriteByte(',')
	jsonEncodePair(b, "level", r.Level.String())
	b.WriteByte(',')
	jsonEncodePair(b, "message", r.Message)
	b.WriteByte('}')
	b.WriteByte('\n')
	return b.Bytes()
}

func FormatJSON() Formatter {
	return new(formatJSON)
}

func jsonEncodePair(b *bytes.Buffer, name, value string) {
	jsonEncodeString(b, name)
	b.WriteByte(':')
	jsonEncodeString(b, value)
}

func jsonEncodeString(b *bytes.Buffer, s string) {
	b.WriteByte('"')
	writeString(b, s)
	b.WriteByte('"')
}
