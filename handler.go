package logl

import (
	"io"
	"sync"
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

type dummyHandler struct{}

func (dummyHandler) Handle(*Record) {}

var DummyHandler Handler = dummyHandler{}

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
	Format Formatter
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

type guardHandler struct {
	guard   sync.Mutex
	handler Handler
}

func (gh *guardHandler) Handle(r *Record) {

	gh.guard.Lock()
	defer gh.guard.Unlock()

	gh.handler.Handle(r)
}

func newGuardHandler(handler Handler) *guardHandler {
	return &guardHandler{handler: handler}
}

var _ Handler = &guardHandler{}
