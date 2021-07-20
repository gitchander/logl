package logl

import (
	"bytes"
	"strconv"
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

var (
	//jsonEncodeString = jsonEncodeString1
	jsonEncodeString = jsonEncodeString2
)

func jsonEncodeString1(b *bytes.Buffer, str string) {
	b.WriteString(strconv.Quote(str))
}

func jsonEncodeString2(b *bytes.Buffer, str string) {
	b.WriteByte('"')
	for _, r := range str {
		switch r {
		case '\n':
			b.WriteString("\\n")
		case '\r':
			b.WriteString("\\r")
		case '\t':
			b.WriteString("\\t")
		case '"', '\\':
			b.WriteByte('\\')
			b.WriteByte(byte(r))
		default:
			b.WriteRune(r)
		}
	}
	b.WriteByte('"')
}
