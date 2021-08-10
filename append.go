// +build ignore

package logl

import (
	"time"
	"unicode/utf8"
)

// Time flags
const (
	tf_DATE         = 1 << iota // the date in the local time zone: 2009/01/23
	tf_TIME                     // the time in the local time zone: 01:23:23
	tf_MICROSECONDS             // microsecond resolution: 01:23:23.123123. assumes TF_TIME.
)

func text_format(r *Record, timeFlag int) (data []byte) {
	data = appendLevel(data, r.Level)
	data = appendTime(data, timeFlag, r.Time)
	data = appendMessage(data, r.Message)
	data = append(data, '\n')
	return
}

func appendLevel(data []byte, level Level) []byte {

	switch level {
	case LevelCritical:
		data = append(data, tag_Critical...)
	case LevelError:
		data = append(data, tag_Error...)
	case LevelWarning:
		data = append(data, tag_Warning...)
	case LevelInfo:
		data = append(data, tag_Info...)
	case LevelDebug:
		data = append(data, tag_Debug...)
	case LevelTrace:
		data = append(data, tag_Trace...)
	}

	data = append(data, ' ')

	return data
}

func appendTime(data []byte, flag int, t time.Time) []byte {

	if (flag & (tf_DATE | tf_TIME | tf_MICROSECONDS)) == 0 {
		return data
	}

	if (flag & tf_DATE) != 0 {
		year, month, day := t.Date()
		data = appendIntn(data, year, 4)
		data = append(data, '/')
		data = appendIntn(data, int(month), 2)
		data = append(data, '/')
		data = appendIntn(data, day, 2)
		data = append(data, ' ')
	}

	if (flag & (tf_TIME | tf_MICROSECONDS)) != 0 {
		hour, min, sec := t.Clock()
		data = appendIntn(data, hour, 2)
		data = append(data, ':')
		data = appendIntn(data, min, 2)
		data = append(data, ':')
		data = appendIntn(data, sec, 2)
		if flag&tf_MICROSECONDS != 0 {
			data = append(data, '.')
			data = appendIntn(data, t.Nanosecond()/1e3, 6)
		}
		data = append(data, ' ')
	}

	return data
}

func appendMessage(data []byte, m string) []byte {
	runeBuf := make([]byte, utf8.UTFMax)
	for _, r := range m {
		switch r {
		case '\n':
			data = append(data, '\\', 'n')
		case '\r':
			data = append(data, '\\', 'r')
		case '\t':
			data = append(data, '\\', 't')
		case '"', '\\':
			data = append(data, '\\', byte(r))
		default:
			size := utf8.EncodeRune(runeBuf, r)
			data = append(data, runeBuf[:size]...)
		}
	}
	return data
}

func appendIntn(data []byte, x int, count int) []byte {
	const base = 10
	begin := len(data)
	var digitIndex int
	for i := 0; i < count; i++ {
		x, digitIndex = quoRem(x, base)
		data = append(data, digits[digitIndex])
	}
	flipBytes(data[begin:len(data)])
	return data
}

func flipBytes(data []byte) {
	i, j := 0, len(data)-1
	for i < j {
		data[i], data[j] = data[j], data[i]
		i, j = i+1, j-1
	}
}
