package logl

import (
	"bytes"
	"time"
	"unicode/utf8"
)

const (
	dateSeparator = '/'
	timeSeparator = ':'
)

type FormatText struct {
	HasLevel        bool
	HasDate         bool
	HasTime         bool
	HasMicroseconds bool
	ShieldSpecial   bool // { '\\', '"', '\n', '\r', '\t', '\f' }
}

func (f FormatText) Formatter() Formatter {
	return newFormatterText(f)
}

type formatterText struct {
	hasLevel        bool
	hasDate         bool
	hasTime         bool
	hasMicroseconds bool
	shieldSpecial   bool

	br          bytes.Buffer
	numberBytes []byte
}

func newFormatterText(f FormatText) *formatterText {
	return &formatterText{
		hasLevel:        f.HasLevel,
		hasDate:         f.HasDate,
		hasTime:         f.HasTime,
		hasMicroseconds: f.HasMicroseconds,
		shieldSpecial:   f.ShieldSpecial,

		numberBytes: make([]byte, 6),
	}
}

func (f *formatterText) Format(r *Record) []byte {

	b := &(f.br)

	b.Reset()
	f.writeTextFormat(r)
	return b.Bytes()
}

func (f *formatterText) writeTextFormat(r *Record) {

	b := &(f.br)

	if f.hasLevel {
		b.WriteString(r.Level.shortName())
		b.WriteByte(' ')
	}

	if f.hasDate || f.hasTime || f.hasMicroseconds {
		f.writeTime(r.Time)
	}

	if f.shieldSpecial {
		writeString(b, r.Message)
	} else {
		b.WriteString(r.Message)
	}
	if f.shieldSpecial || !lastByteIs(r.Message, '\n') {
		b.WriteByte('\n')
	}
}

func (f *formatterText) writeTime(t time.Time) {

	b := &(f.br)

	if f.hasDate {
		year, month, day := t.Date()
		b.Write(f.formatIntnBytes(year, 4))
		b.WriteByte(dateSeparator)
		b.Write(f.formatIntnBytes(int(month), 2))
		b.WriteByte(dateSeparator)
		b.Write(f.formatIntnBytes(day, 2))
		b.WriteByte(' ')
	}
	if f.hasTime || f.hasMicroseconds {
		hour, min, sec := t.Clock()
		b.Write(f.formatIntnBytes(hour, 2))
		b.WriteByte(timeSeparator)
		b.Write(f.formatIntnBytes(min, 2))
		b.WriteByte(timeSeparator)
		b.Write(f.formatIntnBytes(sec, 2))
		if f.hasMicroseconds {
			b.WriteByte('.')
			microseconds := t.Nanosecond() / 1e3
			b.Write(f.formatIntnBytes(microseconds, 6))
		}
		b.WriteByte(' ')
	}
}

func (f *formatterText) formatIntnBytes(x, n int) []byte {
	const base = 10
	return f.formatIntnBytesBase(x, n, base)
}

func (f *formatterText) formatIntnBytesBase(x, n, base int) []byte {
	data := f.numberBytes
	formatIntn(x, n, base, data)
	return data[(len(data) - n):]
}

var digits = []byte("0123456789abcdefghijklmnopqrstuvwxyz")

func formatIntn(x, n int, base int, data []byte) {
	var digitIndex int
	j := len(data) - 1
	for i := 0; i < n; i++ {
		x, digitIndex = quoRem(x, base)
		data[j] = digits[digitIndex]
		j--
	}
}

// json: (e *encodeState) stringBytes(...)
func writeString(b *bytes.Buffer, s string) {
	for _, r := range s {
		if uint32(r) < utf8.RuneSelf { // one byte represent
			x := byte(r) // rune to byte
			if safeSet[x] {
				b.WriteByte(x)
				continue
			}
			if sc, ok := byteSpecialChar(x); ok {
				b.WriteByte('\\')
				b.WriteByte(sc)
			} else {
				// \u0000 - 4 hex digits
				b.WriteString("\\u00")
				writeByteHex(b, x)
			}
		} else {
			b.WriteRune(r)
		}
	}
}

// \f - formfeed
// \t - horizontal tab

func byteSpecialChar(b byte) (byte, bool) {
	switch b {
	case '\\', '"':
		return b, true
	case '\n':
		return 'n', true
	case '\r':
		return 'r', true
	case '\t':
		return 't', true
	case '\f':
		return 'f', true
	default:
		return 0, false
	}
}
