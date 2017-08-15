package logl

import (
	"io"
)

type Handler interface {
	Handle(*Record)
}

func FuncHandler(fn func(*Record)) Handler {
	return funcHandler(fn)
}

type funcHandler func(*Record)

func (h funcHandler) Handle(r *Record) {
	h(r)
}

type fakeHandler struct{}

func (fakeHandler) Handle(*Record) {}

var _ Handler = fakeHandler{}

func getHandler(handler Handler) Handler {
	if handler != nil {
		return handler
	} else {
		return fakeHandler{}
	}
}

type recordWriter struct {
	out    io.Writer
	prefix string
	flag   int

	buf []byte
}

func StreamHandler(w io.Writer, prefix string, flag int) Handler {
	return &recordWriter{
		out:    w,
		prefix: prefix,
		flag:   flag,
	}
}

func (rw *recordWriter) Handle(r *Record) {

	data := rw.buf[:0]

	data = append(data, rw.prefix...)
	data = append_level(data, r.Level)
	data = append_time(data, rw.flag, r.Time)
	data = append_message(data, r.Message)

	rw.buf = data

	rw.out.Write(rw.buf)
}
