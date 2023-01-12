package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/gitchander/logl"
	"github.com/gitchander/logl/logut"
)

func main() {
	exampleLogStdout()
	exampleLogOff()
	exampleLogFile()
	exampleThreads()
	exampleUnrep()
	examplePanic()
	exampleUseLogrus()
}

func use(l logl.Logger) {
	r := newRandNow()
	for i := 0; i < 100; i++ {
		randLogMessage(r, l)
	}
}

func openLogFile(filename string) (*os.File, error) {
	return os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
}

func exampleLogStdout() {
	l := logl.NewLoggerRW(
		&logl.FormatWriter{
			Writer: os.Stdout,
			Formatter: logl.FormatText{
				HasLevel:      true,
				HasTime:       true,
				ShieldSpecial: true,
			}.Formatter(),
		},
	)
	use(l)
	l.Errorf("my error no %d", 78)
}

func exampleLogOff() {
	l := logl.NewLoggerRW(
		&logl.FormatWriter{
			Writer: os.Stdout,
			Formatter: logl.FormatText{
				HasLevel: true,
				HasTime:  true,
			}.Formatter(),
		},
	)
	l.SetLevel(logl.LevelOff)
	use(l)
}

func exampleLogFile() {

	file, err := openLogFile("test1.log")
	if err != nil {
		return
	}
	defer file.Close()

	bw := bufio.NewWriter(file)
	defer bw.Flush()

	l := logl.NewLoggerRW(
		&logl.FormatWriter{
			Writer: bw,
			Formatter: logl.FormatText{
				HasLevel:        true,
				HasDate:         true,
				HasMicroseconds: true,
				ShieldSpecial:   true,
			}.Formatter(),
		},
	)
	l.SetLevel(logl.LevelWarning)

	use(l)
}

func examplePanic() {

	sh := &logl.FormatWriter{
		Writer: os.Stdout,
		Formatter: logl.FormatText{
			HasLevel: true,
			HasTime:  true,
		}.Formatter(),
	}

	l := logl.NewLoggerRW(
		logl.FuncRecordWriter(func(r *logl.Record) error {
			err := sh.WriteRecord(r)
			if err != nil {
				return err
			}
			if r.Level == logl.LevelCritical {
				panic(r.Message)
			}
			return nil
		}),
	)
	l.SetLevel(logl.LevelInfo)

	panicMessageRecover(l)

	l.Info("Success info!")
}

func panicMessageRecover(l logl.Logger) {

	defer func() {
		message := recover()
		if message != nil {
			fmt.Println("recover panic:", message)
		}
	}()

	l.Critical("Message with panic")
}

func exampleThreads() {

	file, err := openLogFile("test2.log")
	if err != nil {
		return
	}
	defer file.Close()

	bw := bufio.NewWriter(file)
	defer bw.Flush()

	l := logl.NewLoggerRW(
		logl.MultiRecordWriter(
			logl.DummyRecordWriter(),
			&logl.FormatWriter{
				Writer:    bw,
				Formatter: logl.FormatJSON(),
			},
			&logl.FormatWriter{
				Writer: os.Stdout,
				//Formatter: logl.FormatJSON(),
				//Formatter: new(customTextFormat),
				Formatter: logl.FormatText{
					HasLevel:        true,
					HasDate:         true,
					HasTime:         true,
					HasMicroseconds: true,
					ShieldSpecial:   true,
				}.Formatter(),
			},
		),
	)
	l.SetLevel(logl.LevelTrace)

	var wg sync.WaitGroup
	const n = 10
	wg.Add(n)
	rr := newRandNow()
	for i := 0; i < n; i++ {
		newSeed := rr.Int63()
		go func(id int, seed int64) {
			r := newRandSeed(seed)
			for j := 0; j < 20; j++ {
				var (
					level = randLogLevel(r)
					// message = fmt.Sprintf("id(%d):%s", id, randLine(r, randIntRange(r, 3, 8)))
					message = randLine(r, randIntRange(r, 3, 8))
				)
				l.Log(level, message)
			}
			wg.Done()
		}(i, newSeed)
	}
	wg.Wait()
}

type customTextFormat struct {
	buf bytes.Buffer
}

func (p *customTextFormat) Format(r *logl.Record) []byte {
	p.buf.Reset()
	fmt.Fprintf(&(p.buf), "%s [%-8s] %s\n", r.Time.Format("2006/01/02 15:04:05"), r.Level, r.Message)
	return p.buf.Bytes()
}

func exampleUnrep() {

	urw := logut.NewUnrepRecordWriter(
		&logl.FormatWriter{
			Writer: os.Stdout,
			Formatter: logl.FormatText{
				HasLevel:      true,
				HasTime:       true,
				ShieldSpecial: true,
			}.Formatter(),
		})
	defer urw.Flush()

	l := logl.NewLoggerRW(urw)
	l.SetLevel(logl.LevelInfo)

	for i := 0; i < 10000; i++ {
		l.Info("Hello, World!")
	}
}

func exampleUseLogrus() {
	l := logut.LoggerByLogrus(logrus.New())
	l.Info("Hello, Logrus!")
	use(l)
}

func exampleLogRotate() {

	lc := logut.DefaultLogConfig

	lw := logut.NewLogWriter(lc.WriterConfig)
	defer lw.Close()

	logger := logl.MakeLogger(lw)
	logger.SetLevel(lc.LogLevel)

	for i := 0; i < 10; i++ {
		logger.Debugf("Message #%d", i+1)
	}

	logger.Info("Before rotate")
	lw.Rotate()
	logger.Info("After rotate")

	for i := 0; i < 10; i++ {
		logger.Debugf("Message #%d", i+1)
	}
}
