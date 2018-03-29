package logl

import (
	"bytes"
)

type formatJSON struct {
	buf bytes.Buffer
}

func (f *formatJSON) Format(r *Record) []byte {
	buf := &(f.buf)
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

func FormatJSON() Format {
	return new(formatJSON)
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
