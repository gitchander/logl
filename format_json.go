package logl

import (
	"bytes"
	"time"
)

const (
	// timeFormatJSON = "2006-01-02 15:04:05"
	// timeFormatJSON = time.RFC3339
	timeFormatJSON = time.RFC3339Nano
)

type formatterJSON struct {
	br bytes.Buffer
}

func (f *formatterJSON) Format(r *Record) []byte {
	b := &(f.br)
	b.Reset()
	b.WriteByte('{')
	f.encodePair("time", r.Time.Format(timeFormatJSON))
	b.WriteByte(',')
	f.encodePair("level", r.Level.String())
	b.WriteByte(',')
	f.encodePair("message", r.Message)
	b.WriteByte('}')
	b.WriteByte('\n')
	return b.Bytes()
}

func FormatJSON() Formatter {
	return new(formatterJSON)
}

func (f *formatterJSON) encodePair(name, value string) {
	f.encodeString(name)
	f.br.WriteByte(':')
	f.encodeString(value)
}

func (f *formatterJSON) encodeString(s string) {
	b := &(f.br)
	b.WriteByte('"')
	writeString(b, s)
	b.WriteByte('"')
}
