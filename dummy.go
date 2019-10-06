package logl

type dummyLogger struct{}

var _ Logger = dummyLogger{}

func DummyLogger() Logger {
	return dummyLogger{}
}

func (dummyLogger) Log(level Level, vs ...interface{})                 {}
func (dummyLogger) Logf(level Level, format string, vs ...interface{}) {}
func (dummyLogger) Level() Level                                       { return LevelOff }
func (dummyLogger) SetLevel(Level)                                     {}

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
