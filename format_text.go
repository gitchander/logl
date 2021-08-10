package logl

import (
	"bytes"
	"time"
	"unicode/utf8"
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
	HasLevel        bool
	HasDate         bool
	HasTime         bool
	HasMicroseconds bool
	ShieldSpecial   bool // { '\n', '\r', '\t', '\\' }
}

func (f *FormatText) Format(r *Record) []byte {
	var b bytes.Buffer
	writeTextFormat(&b, r, f)
	return b.Bytes()
}

func writeTextFormat(b *bytes.Buffer, r *Record, f *FormatText) {

	if f.HasLevel {
		writeLevel(b, r.Level)
		b.WriteByte(' ')
	}

	if f.HasDate || f.HasTime || f.HasMicroseconds {
		writeTime(b, f, r.Time)
	}

	if f.ShieldSpecial {
		writeString(b, r.Message)
	} else {
		b.WriteString(r.Message)
	}
	if f.ShieldSpecial || !lastByteIs(r.Message, '\n') {
		b.WriteByte('\n')
	}
}

func writeLevel(b *bytes.Buffer, level Level) {
	switch level {
	case LevelCritical:
		b.Write(tag_Critical)
	case LevelError:
		b.Write(tag_Error)
	case LevelWarning:
		b.Write(tag_Warning)
	case LevelInfo:
		b.Write(tag_Info)
	case LevelDebug:
		b.Write(tag_Debug)
	case LevelTrace:
		b.Write(tag_Trace)
	}
}

const (
	dateSeparator = '/'
	timeSeparator = ':'
)

func writeTime(b *bytes.Buffer, f *FormatText, t time.Time) {
	const base = 10
	if f.HasDate {
		year, month, day := t.Date()
		b.Write(formatIntn(year, base, 4))
		b.WriteByte(dateSeparator)
		b.Write(formatIntn(int(month), base, 2))
		b.WriteByte(dateSeparator)
		b.Write(formatIntn(day, base, 2))
		b.WriteByte(' ')
	}
	if f.HasTime || f.HasMicroseconds {
		hour, min, sec := t.Clock()
		b.Write(formatIntn(hour, base, 2))
		b.WriteByte(timeSeparator)
		b.Write(formatIntn(min, base, 2))
		b.WriteByte(timeSeparator)
		b.Write(formatIntn(sec, base, 2))
		if f.HasMicroseconds {
			b.WriteByte('.')
			microseconds := t.Nanosecond() / 1e3
			b.Write(formatIntn(microseconds, base, 6))
		}
		b.WriteByte(' ')
	}
}

var digits = []byte("0123456789abcdefghijklmnopqrstuvwxyz")

func formatIntn(x int, base int, n int) []byte {
	data := make([]byte, n)
	var digitIndex int
	for i := n; i > 0; i-- {
		x, digitIndex = quoRem(x, base)
		data[i-1] = digits[digitIndex]
	}
	return data
}

func writeString(b *bytes.Buffer, s string) {

	// json: (e *encodeState) stringBytes(...)

	for _, r := range s {
		if uint32(r) < utf8.RuneSelf { // one byte represent

			x := byte(r) // rune to byte

			if safeSet[x] {
				b.WriteByte(x)
				continue
			}

			b.WriteByte('\\')
			switch x {
			case '\\', '"':
				b.WriteByte(x)
			case '\n':
				b.WriteByte('n')
			case '\r':
				b.WriteByte('r')
			case '\t':
				b.WriteByte('t')
			default:
				b.WriteString(`u00`)
				writeByteHex(b, x)
			}
		} else {
			b.WriteRune(r)
		}
	}
}
