package logl

import (
	"bytes"
	"time"
)

var (
	tag_Critical = []byte("CRI")
	tag_Error    = []byte("ERR")
	tag_Warning  = []byte("WAR")
	tag_Info     = []byte("INF")
	tag_Debug    = []byte("DEB")
	tag_Trace    = []byte("TRA")
)

type FormatText struct {
	HasLevel      bool
	Date          bool
	Time          bool
	Microseconds  bool
	ShieldSpecial bool // { '\n', '\r', '\t', '\\' }
}

func (f *FormatText) Format(r *Record) []byte {
	var buf bytes.Buffer
	writeTextFormat(&buf, r, f)
	return buf.Bytes()
}

func writeTextFormat(buf *bytes.Buffer, r *Record, f *FormatText) {

	if f.HasLevel {
		writeLevel(buf, r.Level)
		buf.WriteByte(' ')
	}

	if f.Date || f.Time || f.Microseconds {
		writeTime(buf, f, r.Time)
	}

	if f.ShieldSpecial {
		writeShieldSpecial(buf, r.Message)
	} else {
		buf.WriteString(r.Message)
	}
	if f.ShieldSpecial || !lastByteIs(r.Message, '\n') {
		buf.WriteByte('\n')
	}
}

func writeLevel(buf *bytes.Buffer, level Level) {
	switch level {
	case LevelCritical:
		buf.Write(tag_Critical)
	case LevelError:
		buf.Write(tag_Error)
	case LevelWarning:
		buf.Write(tag_Warning)
	case LevelInfo:
		buf.Write(tag_Info)
	case LevelDebug:
		buf.Write(tag_Debug)
	case LevelTrace:
		buf.Write(tag_Trace)
	}
}

const (
	dateSeparator = '/'
	timeSeparator = ':'
)

func writeTime(buf *bytes.Buffer, f *FormatText, t time.Time) {
	const base = 10
	if f.Date {
		year, month, day := t.Date()
		buf.Write(formatIntn(year, base, 4))
		buf.WriteByte(dateSeparator)
		buf.Write(formatIntn(int(month), base, 2))
		buf.WriteByte(dateSeparator)
		buf.Write(formatIntn(day, base, 2))
		buf.WriteByte(' ')
	}
	if f.Time || f.Microseconds {
		hour, min, sec := t.Clock()
		buf.Write(formatIntn(hour, base, 2))
		buf.WriteByte(timeSeparator)
		buf.Write(formatIntn(min, base, 2))
		buf.WriteByte(timeSeparator)
		buf.Write(formatIntn(sec, base, 2))
		if f.Microseconds {
			buf.WriteByte('.')
			microseconds := t.Nanosecond() / 1e3
			buf.Write(formatIntn(microseconds, base, 6))
		}
		buf.WriteByte(' ')
	}
}

var digits = []byte("0123456789abcdefghijklmnopqrstuvwxyz")

func formatIntn(d int, base int, n int) []byte {
	data := make([]byte, n)
	for i := n; i > 0; i-- {
		quo, rem := quoRem(d, base)
		data[i-1] = digits[rem]
		d = quo
	}
	return data
}

func writeShieldSpecial(buf *bytes.Buffer, m string) {
	for _, r := range m {
		switch r {
		case '\n':
			buf.WriteByte('\\')
			buf.WriteByte('n')
		case '\r':
			buf.WriteByte('\\')
			buf.WriteByte('r')
		case '\t':
			buf.WriteByte('\\')
			buf.WriteByte('t')
		case '\\':
			buf.WriteByte('\\')
			buf.WriteByte('\\')
		default:
			buf.WriteRune(r)
		}
	}
}
