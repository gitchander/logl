package logl

import (
	"io"
	"time"
)

type Handler interface {
	Handle(*Record)
}

type Record struct {
	Time    time.Time
	Level   Level
	Message string
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

var FakeHandler Handler = fakeHandler{}

func MultiHandler(hs ...Handler) Handler {
	return FuncHandler(
		func(r *Record) {
			for _, h := range hs {
				h.Handle(r)
			}
		},
	)
}

func FilterHandler(filter func(*Record) bool, handler Handler) Handler {
	return FuncHandler(
		func(r *Record) {
			if filter(r) {
				handler.Handle(r)
			}
		},
	)
}

type StreamHandler struct {
	Output io.Writer
	Format Format
}

func (p *StreamHandler) Handle(r *Record) {
	data := p.Format.Format(r)
	p.Output.Write(data)
}

//func StreamHandler(out io.Writer, format Format) Handler {
//	return FuncHandler(
//		func(r *Record) {
//			data := format.Format(r)
//			out.Write(data)
//		},
//	)
//}
