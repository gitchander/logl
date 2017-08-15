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
		locker:  getLocker(false),
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

func (l *Logger) handleMessage(level Level, message string) {
	l.locker.Lock()
	if level <= l.level {
		r := &Record{
			Time:    time.Now(),
			Level:   level,
			Message: message,
		}
		l.handler.Handle(r)
	}
	l.locker.Unlock()
}

func (l *Logger) Critical(v ...interface{}) {
	l.handleMessage(LevelCritical, fmt.Sprint(v...))
}

func (l *Logger) Criticalf(format string, v ...interface{}) {
	l.handleMessage(LevelCritical, fmt.Sprintf(format, v...))
}

func (l *Logger) Error(v ...interface{}) {
	l.handleMessage(LevelError, fmt.Sprint(v...))
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.handleMessage(LevelError, fmt.Sprintf(format, v...))
}

func (l *Logger) Warning(v ...interface{}) {
	l.handleMessage(LevelWarning, fmt.Sprint(v...))
}

func (l *Logger) Warningf(format string, v ...interface{}) {
	l.handleMessage(LevelWarning, fmt.Sprintf(format, v...))
}

func (l *Logger) Info(v ...interface{}) {
	l.handleMessage(LevelInfo, fmt.Sprint(v...))
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.handleMessage(LevelInfo, fmt.Sprintf(format, v...))
}

func (l *Logger) Debug(v ...interface{}) {
	l.handleMessage(LevelDebug, fmt.Sprint(v...))
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.handleMessage(LevelDebug, fmt.Sprintf(format, v...))
}
