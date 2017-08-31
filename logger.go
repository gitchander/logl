package logl

import (
	"fmt"
	"sync"
	"time"
)

type Config struct {
	Handler Handler
	Level   Level
}

type Logger struct {
	locker  sync.Locker
	handler Handler
	level   Level
}

func New(c Config) *Logger {
	return &Logger{
		locker:  new(sync.Mutex),
		handler: c.Handler,
		level:   c.Level,
	}
}

func (l *Logger) SetHandler(handler Handler) {
	l.locker.Lock()
	if handler != nil {
		l.handler = handler
	} else {
		l.handler = FakeHandler
	}
	l.locker.Unlock()
}

func (l *Logger) Level() (level Level) {
	l.locker.Lock()
	level = l.level
	l.locker.Unlock()
	return
}

func (l *Logger) SetLevel(level Level) {
	l.locker.Lock()
	l.level = level
	l.locker.Unlock()
}

func (l *Logger) handleMessage(level Level, format *string, vs ...interface{}) {
	l.locker.Lock()
	defer l.locker.Unlock() // unlocks even if the handler call panic

	if level <= l.level {

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
}

func (l *Logger) Message(level Level, vs ...interface{}) {
	l.handleMessage(level, nil, vs...)
}

func (l *Logger) Messagef(level Level, format string, vs ...interface{}) {
	l.handleMessage(level, &format, vs...)
}

func (l *Logger) Critical(vs ...interface{}) {
	l.handleMessage(LevelCritical, nil, vs...)
}

func (l *Logger) Criticalf(format string, vs ...interface{}) {
	l.handleMessage(LevelCritical, &format, vs...)
}

func (l *Logger) Error(vs ...interface{}) {
	l.handleMessage(LevelError, nil, vs...)
}

func (l *Logger) Errorf(format string, vs ...interface{}) {
	l.handleMessage(LevelError, &format, vs...)
}

func (l *Logger) Warning(vs ...interface{}) {
	l.handleMessage(LevelWarning, nil, vs...)
}

func (l *Logger) Warningf(format string, vs ...interface{}) {
	l.handleMessage(LevelWarning, &format, vs...)
}

func (l *Logger) Info(vs ...interface{}) {
	l.handleMessage(LevelInfo, nil, vs...)
}

func (l *Logger) Infof(format string, vs ...interface{}) {
	l.handleMessage(LevelInfo, &format, vs...)
}

func (l *Logger) Debug(vs ...interface{}) {
	l.handleMessage(LevelDebug, nil, vs...)
}

func (l *Logger) Debugf(format string, vs ...interface{}) {
	l.handleMessage(LevelDebug, &format, vs...)
}
