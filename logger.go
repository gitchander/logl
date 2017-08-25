package logl

import (
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
	l.handler = handler
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
	defer l.locker.Unlock()
	if level <= l.level {
		r := &Record{
			Time:    time.Now(),
			Level:   level,
			Message: message,
		}
		l.handler.Handle(r)
	}
}

func (l *Logger) Message(level Level, message string) {
	l.handleMessage(level, message)
}

func (l *Logger) Critical(message string) {
	l.handleMessage(LevelCritical, message)
}

func (l *Logger) Error(message string) {
	l.handleMessage(LevelError, message)
}

func (l *Logger) Warning(message string) {
	l.handleMessage(LevelWarning, message)
}

func (l *Logger) Info(message string) {
	l.handleMessage(LevelInfo, message)
}

func (l *Logger) Debug(message string) {
	l.handleMessage(LevelDebug, message)
}
