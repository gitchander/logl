package logl

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// standard
type rwLogger struct {
	atomicLevel int32

	guard sync.Mutex // to protect rw
	rw    RecordWriter
}

var _ Logger = &rwLogger{}

func NewLoggerRW(rw RecordWriter) Logger {
	if rw == nil {
		rw = DummyRecordWriter()
	}
	l := &rwLogger{rw: rw}
	l.SetLevel(LevelError)
	return l
}

func (l *rwLogger) Level() Level {
	return Level(atomic.LoadInt32(&(l.atomicLevel)))
}

func (l *rwLogger) SetLevel(level Level) {
	atomic.StoreInt32(&(l.atomicLevel), int32(level))
}

func (l *rwLogger) handleMessage(level Level, format *string, vs ...interface{}) {

	if level > l.Level() {
		return
	}

	var message string
	if format != nil {
		message = fmt.Sprintf(*format, vs...)
	} else {
		message = fmt.Sprint(vs...)
	}

	r := &Record{
		Time:    time.Now(),
		Level:   level,
		Message: message,
	}

	l.guard.Lock()
	defer l.guard.Unlock()

	l.rw.WriteRecord(r)
}

func (l *rwLogger) Log(level Level, vs ...interface{}) {
	l.handleMessage(level, nil, vs...)
}

func (l *rwLogger) Logf(level Level, format string, vs ...interface{}) {
	l.handleMessage(level, &format, vs...)
}

func (l *rwLogger) Critical(vs ...interface{}) {
	l.handleMessage(LevelCritical, nil, vs...)
}

func (l *rwLogger) Criticalf(format string, vs ...interface{}) {
	l.handleMessage(LevelCritical, &format, vs...)
}

func (l *rwLogger) Error(vs ...interface{}) {
	l.handleMessage(LevelError, nil, vs...)
}

func (l *rwLogger) Errorf(format string, vs ...interface{}) {
	l.handleMessage(LevelError, &format, vs...)
}

func (l *rwLogger) Warning(vs ...interface{}) {
	l.handleMessage(LevelWarning, nil, vs...)
}

func (l *rwLogger) Warningf(format string, vs ...interface{}) {
	l.handleMessage(LevelWarning, &format, vs...)
}

func (l *rwLogger) Info(vs ...interface{}) {
	l.handleMessage(LevelInfo, nil, vs...)
}

func (l *rwLogger) Infof(format string, vs ...interface{}) {
	l.handleMessage(LevelInfo, &format, vs...)
}

func (l *rwLogger) Debug(vs ...interface{}) {
	l.handleMessage(LevelDebug, nil, vs...)
}

func (l *rwLogger) Debugf(format string, vs ...interface{}) {
	l.handleMessage(LevelDebug, &format, vs...)
}

func (l *rwLogger) Trace(vs ...interface{}) {
	l.handleMessage(LevelTrace, nil, vs...)
}

func (l *rwLogger) Tracef(format string, vs ...interface{}) {
	l.handleMessage(LevelTrace, &format, vs...)
}
