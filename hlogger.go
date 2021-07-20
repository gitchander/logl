package logl

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type HandleLogger struct {
	atomicLevel int32

	guard   sync.Mutex // protect handler
	handler Handler
}

var _ Logger = &HandleLogger{}

func NewHandleLogger(level Level, handler Handler) *HandleLogger {
	var l HandleLogger
	l.SetLevel(level)
	l.SetHandler(handler)
	return &l
}

func (l *HandleLogger) Level() Level {
	return Level(atomic.LoadInt32(&(l.atomicLevel)))
}

func (l *HandleLogger) SetLevel(level Level) {
	atomic.StoreInt32(&(l.atomicLevel), int32(level))
}

func (l *HandleLogger) SetHandler(handler Handler) {

	if handler == nil {
		handler = DummyHandler
	}

	l.guard.Lock()
	defer l.guard.Unlock()

	l.handler = handler
}

func (l *HandleLogger) handleMessage(level Level, format *string, vs ...interface{}) {

	if level > l.Level() {
		return
	}

	var r = Record{
		Time:  time.Now(),
		Level: level,
	}
	if format != nil {
		r.Message = fmt.Sprintf(*format, vs...)
	} else {
		r.Message = fmt.Sprint(vs...)
	}

	l.guard.Lock()
	defer l.guard.Unlock()

	l.handler.Handle(&r)
}

func (l *HandleLogger) Log(level Level, vs ...interface{}) {
	l.handleMessage(level, nil, vs...)
}

func (l *HandleLogger) Logf(level Level, format string, vs ...interface{}) {
	l.handleMessage(level, &format, vs...)
}

func (l *HandleLogger) Critical(vs ...interface{}) {
	l.handleMessage(LevelCritical, nil, vs...)
}

func (l *HandleLogger) Criticalf(format string, vs ...interface{}) {
	l.handleMessage(LevelCritical, &format, vs...)
}

func (l *HandleLogger) Error(vs ...interface{}) {
	l.handleMessage(LevelError, nil, vs...)
}

func (l *HandleLogger) Errorf(format string, vs ...interface{}) {
	l.handleMessage(LevelError, &format, vs...)
}

func (l *HandleLogger) Warning(vs ...interface{}) {
	l.handleMessage(LevelWarning, nil, vs...)
}

func (l *HandleLogger) Warningf(format string, vs ...interface{}) {
	l.handleMessage(LevelWarning, &format, vs...)
}

func (l *HandleLogger) Info(vs ...interface{}) {
	l.handleMessage(LevelInfo, nil, vs...)
}

func (l *HandleLogger) Infof(format string, vs ...interface{}) {
	l.handleMessage(LevelInfo, &format, vs...)
}

func (l *HandleLogger) Debug(vs ...interface{}) {
	l.handleMessage(LevelDebug, nil, vs...)
}

func (l *HandleLogger) Debugf(format string, vs ...interface{}) {
	l.handleMessage(LevelDebug, &format, vs...)
}

func (l *HandleLogger) Trace(vs ...interface{}) {
	l.handleMessage(LevelTrace, nil, vs...)
}

func (l *HandleLogger) Tracef(format string, vs ...interface{}) {
	l.handleMessage(LevelTrace, &format, vs...)
}
