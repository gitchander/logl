package logut

import (
	"fmt"
	"hash/maphash"

	"github.com/gitchander/logl"
)

func _() {
	var _ logl.Logger
	var _ logl.RecordWriter
	logl.MakeLoggerStdout()
}

// do not repeat.
// Do not repeat handler
// Unrep - unrepeated, unrepeatable
type UnrepRecordWriter struct {
	rw logl.RecordWriter
	h  maphash.Hash

	record     *logl.Record
	hashSum    uint64
	count      int
	flushCount int
}

var _ logl.RecordWriter = &UnrepRecordWriter{}

func NewUnrepRecordWriter(rw logl.RecordWriter) *UnrepRecordWriter {
	return &UnrepRecordWriter{
		rw: rw,
	}
}

func (p *UnrepRecordWriter) Flush() error {
	return p.flush(false)
}

func (p *UnrepRecordWriter) recordHashSum64(r *logl.Record) uint64 {
	p.h.Reset()
	p.h.WriteByte(byte(r.Level)) // write log level
	p.h.WriteString(r.Message)   // write log message
	return p.h.Sum64()
}

func (p *UnrepRecordWriter) flush(keepRecord bool) error {

	var err error
	if (p.record != nil) && (p.count > 0) {
		message := p.record.Message
		if p.count > 1 {
			message += fmt.Sprintf(" (%d repeats)", p.count)
		}
		r := &logl.Record{
			Time:    p.record.Time,
			Level:   p.record.Level,
			Message: message,
		}
		err = p.rw.WriteRecord(r)
	}

	if (p.record != nil) && keepRecord {
		p.count = 0
		p.flushCount = maxInt(p.flushCount, 1) * 10
	} else {
		// reset record
		p.record = nil
		p.hashSum = 0
		p.count = 0
		p.flushCount = 1
	}
	return err
}

func (p *UnrepRecordWriter) addNewRecord(r *logl.Record, newHashSum uint64) {
	p.record = cloneLogRecord(r)
	p.hashSum = newHashSum
	p.count = 1
	p.flushCount = 1
}

func (p *UnrepRecordWriter) WriteRecord(r *logl.Record) error {

	newHashSum := p.recordHashSum64(r)

	if p.record != nil {
		if p.hashSum == newHashSum {
			p.count++
			p.record.Time = r.Time
		} else {
			p.flush(false)
			p.addNewRecord(r, newHashSum)
		}
	} else {
		p.addNewRecord(r, newHashSum)
	}

	var err error
	if (p.record != nil) && (p.count >= p.flushCount) {
		err = p.flush(true)
	}
	return err
}

func cloneLogRecord(r *logl.Record) *logl.Record {
	if r != nil {
		clone := *r
		return &clone
	}
	return nil
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

//------------------------------------------------------------------------------
// DoNotRepeat error
// NRError
// Unrepeatable logger
// UnrepErrLogger - unrepeated error logger

type UnrepErrLogger struct {
	logger logl.Logger

	message    string
	count      int
	flushCount int
}

func NewUnrepErrLogger(logger logl.Logger) *UnrepErrLogger {
	return &UnrepErrLogger{
		logger: logger,
	}
}

func (p *UnrepErrLogger) Flush() {
	p.flush(false)
}

func (p *UnrepErrLogger) flush(keepMessage bool) {

	if (p.message != "") && (p.count > 0) {
		m := p.message
		if p.count > 1 {
			m += fmt.Sprintf(" (%d repeats)", p.count)
		}
		p.logger.Error(m)
	}

	if (p.message != "") && keepMessage {
		p.count = 0
		p.flushCount = maxInt(p.flushCount, 1) * 10
	} else {
		// reset message
		p.message = ""
		p.count = 0
		p.flushCount = 1
	}
}

func (p *UnrepErrLogger) addNewMessage(newMessage string) {
	p.message = newMessage
	p.count = 1
	p.flushCount = 1
}

func (p *UnrepErrLogger) LogError(err error) {

	if err == nil {
		return
	}

	newMessage := err.Error()

	if p.message != "" {
		if p.message == newMessage {
			p.count++
		} else {
			p.flush(false)
			p.addNewMessage(newMessage)
		}
	} else {
		p.addNewMessage(newMessage)
	}

	if (p.message != "") && (p.count >= p.flushCount) {
		p.flush(true)
	}
}
