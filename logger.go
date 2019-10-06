package logl

// Concept of a simple interface:
type ILogger interface {
	Log(level Level, vs ...interface{})
	Logf(level Level, format string, vs ...interface{})

	Level() Level
	SetLevel(Level)
}

type Logger interface {
	ILogger

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
