package logut

import (
	"fmt"

	"github.com/gitchander/logl"
	"github.com/sirupsen/logrus"
)

type logrusLogger struct {
	l *logrus.Logger
}

var _ logl.Logger = &logrusLogger{}

func LoggerByLogrus(l *logrus.Logger) logl.Logger {
	return &logrusLogger{l}
}

func (p *logrusLogger) Log(level logl.Level, vs ...interface{}) {
	gr := convertLevelGL_To_GR(level)
	p.l.Log(gr, vs...)
}

func (p *logrusLogger) Logf(level logl.Level, format string, vs ...interface{}) {
	gr := convertLevelGL_To_GR(level)
	p.l.Logf(gr, format, vs...)
}

func (p *logrusLogger) Level() logl.Level {
	gr := p.l.GetLevel()
	gl := convertLevelGR_To_GL(gr)
	return gl
}

func (p *logrusLogger) SetLevel(level logl.Level) {
	gr := convertLevelGL_To_GR(level)
	p.l.SetLevel(gr)
}

func (p *logrusLogger) Critical(vs ...interface{}) {
	p.l.Fatal(vs...)
}

func (p *logrusLogger) Error(vs ...interface{}) {
	p.l.Error(vs...)
}

func (p *logrusLogger) Warning(vs ...interface{}) {
	p.l.Warn(vs...)
}

func (p *logrusLogger) Info(vs ...interface{}) {
	p.l.Info(vs...)
}

func (p *logrusLogger) Debug(vs ...interface{}) {
	p.l.Debug(vs...)
}

func (p *logrusLogger) Trace(vs ...interface{}) {
	p.l.Trace(vs...)
}

func (p *logrusLogger) Criticalf(format string, vs ...interface{}) {
	p.l.Fatalf(format, vs...)
}

func (p *logrusLogger) Errorf(format string, vs ...interface{}) {
	p.l.Errorf(format, vs...)
}

func (p *logrusLogger) Warningf(format string, vs ...interface{}) {
	p.l.Warnf(format, vs...)
}

func (p *logrusLogger) Infof(format string, vs ...interface{}) {
	p.l.Infof(format, vs...)
}

func (p *logrusLogger) Debugf(format string, vs ...interface{}) {
	p.l.Debugf(format, vs...)
}

func (p *logrusLogger) Tracef(format string, vs ...interface{}) {
	p.l.Tracef(format, vs...)
}

// ------------------------------------------------------------------------------
var logrusAllLevels = []logrus.Level{
	logrus.PanicLevel,
	logrus.FatalLevel,
	logrus.ErrorLevel,
	logrus.WarnLevel,
	logrus.InfoLevel,
	logrus.DebugLevel,
	logrus.TraceLevel,
}

var loglAllLevels = []logl.Level{
	logl.LevelOff,
	logl.LevelCritical,
	logl.LevelError,
	logl.LevelWarning,
	logl.LevelInfo,
	logl.LevelDebug,
	logl.LevelTrace,
}

// GR - logrus.Level
// GL - logl.Level
// level GR -> GL

func convertLevelGL_To_GR(gl logl.Level) (gr logrus.Level) {
	switch gl {
	case logl.LevelOff:
		gr = logrus.FatalLevel
	case logl.LevelCritical:
		gr = logrus.FatalLevel
	case logl.LevelError:
		gr = logrus.ErrorLevel
	case logl.LevelWarning:
		gr = logrus.WarnLevel
	case logl.LevelInfo:
		gr = logrus.InfoLevel
	case logl.LevelDebug:
		gr = logrus.DebugLevel
	case logl.LevelTrace:
		gr = logrus.TraceLevel
	default:
		err := fmt.Errorf("invalid %s level %d", "logl", gl)
		panic(err)
	}
	return gr
}

func convertLevelGR_To_GL(gr logrus.Level) (gl logl.Level) {
	switch gr {
	case logrus.PanicLevel:
		gl = logl.LevelCritical
	case logrus.FatalLevel:
		gl = logl.LevelCritical
	case logrus.ErrorLevel:
		gl = logl.LevelError
	case logrus.WarnLevel:
		gl = logl.LevelWarning
	case logrus.InfoLevel:
		gl = logl.LevelInfo
	case logrus.DebugLevel:
		gl = logl.LevelDebug
	case logrus.TraceLevel:
		gl = logl.LevelTrace
	default:
		err := fmt.Errorf("invalid %s level %d", "logrus", gr)
		panic(err)
	}
	return gl
}
