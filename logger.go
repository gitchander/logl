package logl

import (
	"fmt"
	"sync"
	"time"
)

type Logger interface {
	SetLevel(Level)
	Level() Level

	Critical(vs ...interface{})
	Error(vs ...interface{})
	Warning(vs ...interface{})
	Info(vs ...interface{})
	Debug(vs ...interface{})
	Trace(vs ...interface{})

	Criticalf(format string, vs ...interface{})
	Errorf(format string, vs ...interface{})
	Warningf(format string, vs ...interface{})
	Infof(format string, vs ...interface{})
	Debugf(format string, vs ...interface{})
	Tracef(format string, vs ...interface{})
}

type dummyLogger struct{}

func (dummyLogger) SetLevel(Level) {}
func (dummyLogger) Level() Level   { return LevelOff }

func (dummyLogger) Critical(vs ...interface{}) {}
func (dummyLogger) Error(vs ...interface{})    {}
func (dummyLogger) Warning(vs ...interface{})  {}
func (dummyLogger) Info(vs ...interface{})     {}
func (dummyLogger) Debug(vs ...interface{})    {}
func (dummyLogger) Trace(vs ...interface{})    {}

func (dummyLogger) Criticalf(format string, vs ...interface{}) {}
func (dummyLogger) Errorf(format string, vs ...interface{})    {}
func (dummyLogger) Warningf(format string, vs ...interface{})  {}
func (dummyLogger) Infof(format string, vs ...interface{})     {}
func (dummyLogger) Debugf(format string, vs ...interface{})    {}
func (dummyLogger) Tracef(format string, vs ...interface{})    {}

var DummyLogger Logger = dummyLogger{}

type HandleLogger struct {
	mutex   sync.Mutex
	handler Handler
	level   Level
}

var _ Logger = &HandleLogger{}

func NewHandleLogger(handler Handler, level Level) *HandleLogger {
	return &HandleLogger{
		handler: handler,
		level:   level,
	}
}

func (l *HandleLogger) SetHandler(handler Handler) {
	l.mutex.Lock()
	if handler != nil {
		l.handler = handler
	} else {
		l.handler = FakeHandler
	}
	l.mutex.Unlock()
}

func (l *HandleLogger) Level() (level Level) {
	l.mutex.Lock()
	level = l.level
	l.mutex.Unlock()
	return
}

func (l *HandleLogger) SetLevel(level Level) {
	l.mutex.Lock()
	l.level = level
	l.mutex.Unlock()
}

func (l *HandleLogger) handleMessage(level Level, format *string, vs ...interface{}) {

	l.mutex.Lock()
	if level > l.level {
		l.mutex.Unlock()
		return
	}
	defer l.mutex.Unlock() // unlocks even if the handler call panic

	r := &Record{
		Time:  time.Now(),
		Level: level,
	}
	if format != nil {
		r.Message = fmt.Sprintf(*format, vs...)
	} else {
		r.Message = fmt.Sprint(vs...)
	}

	l.handler.Handle(r)
}

func (l *HandleLogger) Message(level Level, vs ...interface{}) {
	l.handleMessage(level, nil, vs...)
}

func (l *HandleLogger) Messagef(level Level, format string, vs ...interface{}) {
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
