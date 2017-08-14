package logl

import (
	"fmt"
	"os"
	"sync"
)

type Logger struct {
	locker sync.Locker
	out    WriteFlusher
	prefix string
	level  Level
	flag   int
	buf    []byte
}

func New(c Config) *Logger {
	return &Logger{
		locker: getLocker(c.NotSafe),
		out:    c.Output,
		prefix: c.Prefix,
		level:  c.Level,
		flag:   c.Flag,
	}
}

func (l *Logger) SetOutput(out WriteFlusher) {
	l.locker.Lock()
	l.out.Flush()
	l.out = out
	l.locker.Unlock()
}

func (l *Logger) Prefix() (prefix string) {
	l.locker.Lock()
	prefix = l.prefix
	l.locker.Unlock()
	return
}

func (l *Logger) SetPrefix(prefix string) {
	l.locker.Lock()
	l.prefix = prefix
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

func (l *Logger) Flag() (flag int) {
	l.locker.Lock()
	flag = l.flag
	l.locker.Unlock()
	return
}

func (l *Logger) SetFlag(flag int) {
	l.locker.Lock()
	l.flag = flag
	l.locker.Unlock()
}

func (l *Logger) Panic(v ...interface{}) {
	const level = LEVEL_PANIC
	var message string
	l.locker.Lock()
	if level <= l.level {
		message = fmt.Sprint(v...)
		l.writeMessage(level, message)
		l.out.Flush()
	}
	l.locker.Unlock()
	panic(message)
}

func (l *Logger) Panicf(format string, v ...interface{}) {
	const level = LEVEL_PANIC
	var message string
	l.locker.Lock()
	if level <= l.level {
		message = fmt.Sprintf(format, v...)
		l.writeMessage(level, message)
		l.out.Flush()
	}
	l.locker.Unlock()
	panic(message)
}

func (l *Logger) Fatal(v ...interface{}) {
	const level = LEVEL_FATAL
	l.locker.Lock()
	if level <= l.level {
		l.writeMessage(level, fmt.Sprint(v...))
		l.out.Flush()
	}
	l.locker.Unlock()
	os.Exit(1)
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	const level = LEVEL_FATAL
	l.locker.Lock()
	if level <= l.level {
		l.writeMessage(level, fmt.Sprintf(format, v...))
		l.out.Flush()
	}
	l.locker.Unlock()
	os.Exit(1)
}

func (l *Logger) Error(v ...interface{}) {
	const level = LEVEL_ERROR
	l.locker.Lock()
	if level <= l.level {
		l.writeMessage(level, fmt.Sprint(v...))
	}
	l.locker.Unlock()
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	const level = LEVEL_ERROR
	l.locker.Lock()
	if level <= l.level {
		l.writeMessage(level, fmt.Sprintf(format, v...))
	}
	l.locker.Unlock()
}

func (l *Logger) Warning(v ...interface{}) {
	const level = LEVEL_WARNING
	l.locker.Lock()
	if level <= l.level {
		l.writeMessage(level, fmt.Sprint(v...))
	}
	l.locker.Unlock()
}

func (l *Logger) Warningf(format string, v ...interface{}) {
	const level = LEVEL_WARNING
	l.locker.Lock()
	if level <= l.level {
		l.writeMessage(level, fmt.Sprintf(format, v...))
	}
	l.locker.Unlock()
}

func (l *Logger) Info(v ...interface{}) {
	const level = LEVEL_INFO
	l.locker.Lock()
	if level <= l.level {
		l.writeMessage(level, fmt.Sprint(v...))
	}
	l.locker.Unlock()
}

func (l *Logger) Infof(format string, v ...interface{}) {
	const level = LEVEL_INFO
	l.locker.Lock()
	if level <= l.level {
		l.writeMessage(level, fmt.Sprintf(format, v...))
	}
	l.locker.Unlock()
}

func (l *Logger) Debug(v ...interface{}) {
	const level = LEVEL_DEBUG
	l.locker.Lock()
	if level <= l.level {
		l.writeMessage(level, fmt.Sprint(v...))
	}
	l.locker.Unlock()
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	const level = LEVEL_DEBUG
	l.locker.Lock()
	if level <= l.level {
		l.writeMessage(level, fmt.Sprintf(format, v...))
	}
	l.locker.Unlock()
}

func (l *Logger) writeMessage(level Level, message string) {

	if len(message) == 0 {
		return
	}

	data := l.buf[:0]

	data = append(data, l.prefix...)
	data = append_level(data, level)
	data = append_time(data, l.flag)
	data = append_message(data, message)

	l.buf = data

	l.out.Write(l.buf)
}
