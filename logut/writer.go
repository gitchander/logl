package logut

import (
	"io"
	"os"
)

type Rotator interface {
	Rotate() error
}

type LogWriter interface {
	io.WriteCloser
	Rotator
}

func NewLogWriter(wc WriterConfig) LogWriter {
	var lws []LogWriter
	if wc.EnableStdout {
		lws = append(lws, &onlyWriter{os.Stdout})
	}
	if wc.EnableFile {
		lws = append(lws, newLumberjack(wc.FileConfig))
	}
	switch len(lws) {
	case 0:
		return dummyLogWriter{}
	case 1:
		return lws[0]
	default:
		return &multiLogWriter{lws: lws}
	}
}

// ------------------------------------------------------------------------------
type dummyLogWriter struct{}

var _ LogWriter = dummyLogWriter{}

func (dummyLogWriter) Write(p []byte) (n int, err error) {
	return n, nil
}

func (dummyLogWriter) Close() error {
	return nil
}

func (dummyLogWriter) Rotate() error {
	return nil
}

// ------------------------------------------------------------------------------
type onlyWriter struct {
	io.Writer
}

var _ LogWriter = onlyWriter{}

func (onlyWriter) Close() error {
	return nil
}

func (onlyWriter) Rotate() error {
	return nil
}

// ------------------------------------------------------------------------------
type multiLogWriter struct {
	lws []LogWriter
}

var _ LogWriter = &multiLogWriter{}

func (p *multiLogWriter) Write(data []byte) (n int, err error) {
	for _, lw := range p.lws {
		k, err := lw.Write(data)
		if err != nil {
			return k, err
		}
	}
	return n, nil
}

func (p *multiLogWriter) Close() error {
	for _, lw := range p.lws {
		err := lw.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *multiLogWriter) Rotate() error {
	for _, lw := range p.lws {
		err := lw.Rotate()
		if err != nil {
			return err
		}
	}
	return nil
}
