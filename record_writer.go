package logl

import (
	"io"
	"sync"
)

type RecordWriter interface {
	WriteRecord(*Record) error
}

func FuncRecordWriter(f func(*Record) error) RecordWriter {
	return funcRecordWriter(f)
}

type funcRecordWriter func(*Record) error

func (f funcRecordWriter) WriteRecord(r *Record) error {
	return f(r)
}

type dummyRecordWriter struct{}

func (dummyRecordWriter) WriteRecord(*Record) error { return nil }

func DummyRecordWriter() RecordWriter {
	return dummyRecordWriter{}
}

func MultiRecordWriter(rws ...RecordWriter) RecordWriter {
	return FuncRecordWriter(
		func(r *Record) error {
			for _, rw := range rws {
				err := rw.WriteRecord(r)
				if err != nil {
					return err
				}
			}
			return nil
		},
	)
}

func FilterRecordWriter(filter func(*Record) bool, rw RecordWriter) RecordWriter {
	return FuncRecordWriter(
		func(r *Record) error {
			if filter(r) {
				return rw.WriteRecord(r)
			}
			return nil
		},
	)
}

type FormatWriter struct {
	Writer    io.Writer
	Formatter Formatter
}

func (p *FormatWriter) WriteRecord(r *Record) error {
	data := p.Formatter.Format(r)
	_, err := p.Writer.Write(data)
	return err
}

func NewStreamRW(w io.Writer, f Formatter) RecordWriter {
	return FuncRecordWriter(
		func(r *Record) error {
			data := f.Format(r)
			_, err := w.Write(data)
			return err
		},
	)
}

type syncRecordWriter struct {
	guard sync.Mutex
	rw    RecordWriter
}

func (p *syncRecordWriter) WriteRecord(r *Record) error {

	p.guard.Lock()
	defer p.guard.Unlock()

	return p.rw.WriteRecord(r)
}

func newSyncRecordWriter(rw RecordWriter) *syncRecordWriter {
	return &syncRecordWriter{rw: rw}
}

var _ RecordWriter = &syncRecordWriter{}
