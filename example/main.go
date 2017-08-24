package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"sync"

	"github.com/gitchander/logl"
)

func main() {
	exampleLogStdout()
	exampleLogOff()
	exampleLogFile()
	examplePanicRecover()
	exampleThreads()
}

func use(l *logl.Logger) {
	r := newRand()
	for i := 0; i < 100; i++ {
		randLogMessage(r, l)
	}
}

func exampleLogStdout() {
	c := logl.Config{
		Handler: &logl.StreamHandler{
			Output: os.Stdout,
			Format: logl.TextFormat(logl.TF_TIME),
		},
		Level: logl.LevelDebug,
	}
	l := logl.New(c)
	use(l)
	l.Error(fmt.Sprintf("my error no %d", 78))
}

func exampleLogOff() {
	c := logl.Config{
		Handler: &logl.StreamHandler{
			Output: os.Stdout,
			Format: logl.TextFormat(logl.TF_TIME),
		},
		Level: -1,
	}
	l := logl.New(c)
	use(l)
}

func exampleLogFile() {

	file, err := os.OpenFile("test.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return
	}
	defer file.Close()

	bw := bufio.NewWriter(file)
	defer bw.Flush()

	c := logl.Config{
		Handler: &logl.StreamHandler{
			Output: bw,
			Format: logl.TextFormat(logl.TF_DATE | logl.TF_MICROSECONDS),
		},
		Level: logl.LevelWarning,
	}
	l := logl.New(c)
	use(l)
}

func examplePanicRecover() {

	defer func() {
		message := recover()
		if message != nil {
			fmt.Println("defer panic:", message)
		}
	}()

	sh := &logl.StreamHandler{
		Output: os.Stdout,
		Format: logl.TextFormat(logl.TF_TIME),
	}
	fh := func(r *logl.Record) {
		sh.Handle(r)
		if r.Level == logl.LevelCritical {
			panic(r.Message)
		}
	}

	c := logl.Config{
		Handler: logl.FuncHandler(fh),
		Level:   logl.LevelError,
	}
	l := logl.New(c)
	l.Critical("my panic message")
}

func exampleThreads() {
	file, err := os.OpenFile("test1.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return
	}
	defer file.Close()

	bw := bufio.NewWriter(file)
	defer bw.Flush()

	logger := logl.New(
		logl.Config{
			Handler: logl.MultiHandler(
				logl.FakeHandler,
				&logl.StreamHandler{
					Output: bw,
					Format: logl.JsonFormat(),
				},
				&logl.StreamHandler{
					Output: os.Stdout,
					//Format: logl.JsonFormat(),
					//Format: new(customTextFormat),
					Format: logl.TextFormat(logl.TF_DATE | logl.TF_MICROSECONDS),
				},
			),
			Level: logl.LevelWarning,
		},
	)
	var wg sync.WaitGroup
	const n = 100
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(id int) {
			r := newRandSeed(int64(id))
			for j := 0; j < 100; j++ {
				var (
					level   = randLevel(r)
					message = fmt.Sprintf("id(%d): %s", id, randLine(r, randIntRange(r, 3, 8)))
				)
				logger.Message(level, message)
			}
			wg.Done()
		}(i)
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
