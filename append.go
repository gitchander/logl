package logl

import "time"

// Flags
const (
	Ldate         = 1 << iota // the date in the local time zone: 2009/01/23
	Ltime                     // the time in the local time zone: 01:23:23
	Lmicroseconds             // microsecond resolution: 01:23:23.123123.  assumes Ltime.
)

var (
	tag_Critical = []byte("CRI")
	tag_Error    = []byte("ERR")
	tag_Warning  = []byte("WAR")
	tag_Info     = []byte("INF")
	tag_Debug    = []byte("DEB")
)

func append_level(data []byte, level Level) []byte {

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
	}

	data = append(data, ' ')

	return data
}

func append_time(data []byte, flag int, t time.Time) []byte {

	if flag&(Ldate|Ltime|Lmicroseconds) == 0 {
		return data
	}

	if flag&Ldate != 0 {
		year, month, day := t.Date()
		data = append_intc(data, year, 4)
		data = append(data, '/')
		data = append_intc(data, int(month), 2)
		data = append(data, '/')
		data = append_intc(data, day, 2)
		data = append(data, ' ')
	}

	if flag&(Ltime|Lmicroseconds) != 0 {
		hour, min, sec := t.Clock()
		data = append_intc(data, hour, 2)
		data = append(data, ':')
		data = append_intc(data, min, 2)
		data = append(data, ':')
		data = append_intc(data, sec, 2)
		if flag&Lmicroseconds != 0 {
			data = append(data, '.')
			data = append_intc(data, t.Nanosecond()/1e3, 6)
		}
		data = append(data, ' ')
	}

	return data
}

func append_message(data []byte, m string) []byte {

	data = append(data, m...)

	if !lastByteIs(m, '\n') {
		data = append(data, '\n')
	}

	return data
}

func lastByteIs(s string, b byte) bool {
	if n := len(s); n > 0 {
		return s[n-1] == b
	}
	return false
}

func append_intc(data []byte, x int, count int) []byte {
	begin := len(data)
	for i := 0; i < count; i++ {
		quo, rem := quoRem(x, 10)
		data = append(data, byte('0'+rem))
		x = quo
	}
	flip(data[begin:len(data)])
	return data
}

func quoRem(x, y int) (quo, rem int) {
	quo = x / y
	rem = x - quo*y
	return
}

func flip(data []byte) {
	i, j := 0, len(data)-1
	for i < j {
		data[i], data[j] = data[j], data[i]
		i, j = i+1, j-1
	}
}
