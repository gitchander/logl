package logl

import (
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

func getHandler(handler Handler) Handler {
	if handler != nil {
		return handler
	} else {
		return fakeHandler{}
	}
}
