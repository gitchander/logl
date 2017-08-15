package logl

import (
	"fmt"
	"sync"
	"time"
)

type Logger struct {
	locker  sync.Locker
	handler Handler
	level   Level
}

func New(c Config) *Logger {
	return &Logger{
		locker:  getLocker(c.NotSafe),
		handler: getHandler(c.Handler),
		level:   c.Level,
	}
}

func (l *Logger) SetHandler(handler Handler) {
	l.locker.Lock()
	l.handler = getHandler(handler)
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

func (l *Logger) Critical(v ...interface{}) {
	l.writeMessage(LevelCritical, v...)
}

func (l *Logger) Criticalf(format string, v ...interface{}) {
	l.writeMessagef(LevelCritical, format, v...)
}

func (l *Logger) Error(v ...interface{}) {
	l.writeMessage(LevelError, v...)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.writeMessagef(LevelError, format, v...)
}

func (l *Logger) Warning(v ...interface{}) {
	l.writeMessage(LevelWarning, v...)
}

func (l *Logger) Warningf(format string, v ...interface{}) {
	l.writeMessagef(LevelWarning, format, v...)
}

func (l *Logger) Info(v ...interface{}) {
	l.writeMessage(LevelInfo, v...)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.writeMessagef(LevelInfo, format, v...)
}

func (l *Logger) Debug(v ...interface{}) {
	l.writeMessage(LevelDebug, v...)
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.writeMessagef(LevelDebug, format, v...)
}

func (l *Logger) writeMessage(level Level, v ...interface{}) {
	l.locker.Lock()
	if level <= l.level {
		r := &Record{
			Time:    time.Now(),
			Level:   level,
			Message: fmt.Sprint(v...),
		}
		l.handler.Handle(r)
	}
	l.locker.Unlock()
}

func (l *Logger) writeMessagef(level Level, format string, v ...interface{}) {
	l.locker.Lock()
	if level <= l.level {
		r := &Record{
			Time:    time.Now(),
			Level:   level,
			Message: fmt.Sprintf(format, v...),
		}
		l.handler.Handle(r)
	}
	l.locker.Unlock()
}
