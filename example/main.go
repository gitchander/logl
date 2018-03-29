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
	examplePanic()
	exampleThreads()
}

func use(l logl.Logger) {
	r := newRandFromTime()
	for i := 0; i < 100; i++ {
		randLogMessage(r, l)
	}
}

func exampleLogStdout() {
	l := logl.NewHandleLogger(
		&logl.StreamHandler{
			Output: os.Stdout,
			Format: &logl.FormatText{
				HasLevel:      true,
				Time:          true,
				ShieldSpecial: true,
			},
		},
		logl.LevelTrace,
	)
	use(l)
	l.Error("my error no %d", 78)
}

func exampleLogOff() {
	l := logl.NewHandleLogger(
		&logl.StreamHandler{
			Output: os.Stdout,
			Format: &logl.FormatText{
				HasLevel: true,
				Time:     true,
			},
		},
		logl.LevelOff,
	)
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

	l := logl.NewHandleLogger(
		&logl.StreamHandler{
			Output: bw,
			Format: &logl.FormatText{
				Date:         true,
				Microseconds: true,
			},
		},
		logl.LevelWarning,
	)
	use(l)
}

func examplePanic() {

	sh := &logl.StreamHandler{
		Output: os.Stdout,
		Format: &logl.FormatText{
			HasLevel: true,
			Time:     true,
		},
	}

	l := logl.NewHandleLogger(
		logl.FuncHandler(func(r *logl.Record) {
			sh.Handle(r)
			if r.Level == logl.LevelCritical {
				panic(r.Message)
			}
		}),
		logl.LevelInfo,
	)

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
	file, err := os.OpenFile("test1.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return
	}
	defer file.Close()

	bw := bufio.NewWriter(file)
	defer bw.Flush()

	l := logl.NewHandleLogger(
		logl.MultiHandler(
			logl.FakeHandler,
			&logl.StreamHandler{
				Output: bw,
				Format: logl.FormatJSON(),
			},
			&logl.StreamHandler{
				Output: os.Stdout,
				//Format: logl.FormatJSON(),
				//Format: new(customTextFormat),
				Format: &logl.FormatText{
					HasLevel:      true,
					Date:          true,
					Time:          true,
					Microseconds:  true,
					ShieldSpecial: true,
				},
			},
		),
		logl.LevelWarning,
	)
	var wg sync.WaitGroup
	const n = 100
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(id int) {
			r := newRandSeed(int64(id))
			for j := 0; j < 100; j++ {
				var (
					level = randLevel(r)
					// message = fmt.Sprintf("id(%d):%s", id, randLine(r, randIntRange(r, 3, 8)))
					message = randLine(r, randIntRange(r, 3, 8))
				)
				logMessage(l, level, message)
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
