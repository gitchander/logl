package logl

import (
	"bytes"
)

type Format interface {
	Format(*Record) []byte
}

type textFormat struct {
	timeFlag int

	buf []byte
}

func (p *textFormat) Format(r *Record) []byte {

	data := p.buf[:0]

	data = append_level(data, r.Level)
	data = append_time(data, p.timeFlag, r.Time)
	data = append_message(data, r.Message)
	data = append(data, '\n')

	p.buf = data

	return cloneBytes(p.buf)
}

// Time flags
const (
	TF_DATE         = 1 << iota // the date in the local time zone: 2009/01/23
	TF_TIME                     // the time in the local time zone: 01:23:23
	TF_MICROSECONDS             // microsecond resolution: 01:23:23.123123. assumes TF_TIME.
)

func TextFormat(timeFlag int) Format {
	return &textFormat{timeFlag: timeFlag}
}

func cloneBytes(a []byte) (b []byte) {
	b = make([]byte, len(a))
	copy(b, a)
	return
}

type jsonFormat struct {
	buf *bytes.Buffer
}

func (p *jsonFormat) Format(r *Record) []byte {
	var buf = p.buf
	buf.Reset()
	buf.WriteByte('{')
	jsonEncodePair(buf, "time", r.Time.Format("2006-01-02 15:04:05"))
	buf.WriteByte(',')
	jsonEncodePair(buf, "level", r.Level.String())
	buf.WriteByte(',')
	jsonEncodePair(buf, "message", r.Message)
	buf.WriteByte('}')
	buf.WriteByte('\n')
	return buf.Bytes()
}

func JsonFormat() Format {
	return &jsonFormat{new(bytes.Buffer)}
}

func jsonEncodePair(buf *bytes.Buffer, name, value string) {
	jsonEncodeString(buf, name)
	buf.WriteByte(':')
	jsonEncodeString(buf, value)
}

func jsonEncodeString(buf *bytes.Buffer, str string) {
	buf.WriteByte('"')
	for _, r := range str {
		switch r {
		case '\n':
			buf.WriteString("\\n")
		case '\r':
			buf.WriteString("\\r")
		case '\t':
			buf.WriteString("\\t")
		case '"', '\\':
			buf.WriteByte('\\')
			buf.WriteByte(byte(r))
		default:
			buf.WriteRune(r)
		}
	}
	buf.WriteByte('"')
}
