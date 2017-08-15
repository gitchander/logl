package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/gitchander/logl"
)

func main() {
	exampleLogStdout()
	exampleLogOff()
	exampleLogFile()
	examplePanicRecover()
}

func use(l *logl.Logger) {
	r := newRand()
	for i := 0; i < 100; i++ {
		var (
			n_word = randIntRange(r, 3, 10)
			line   = randLine(r, n_word)
		)
		switch c := r.Intn(5); c {
		case 0:
			l.Debug(line)
		case 1:
			l.Info(line)
		case 2:
			l.Warning(line)
		case 3:
			l.Error(line)
		case 4:
			l.Critical(line)
		}
	}
}

func exampleLogStdout() {
	c := logl.Config{
		Handler: logl.StreamHandler(os.Stdout, "", logl.Ltime),
		Level:   logl.LevelDebug,
	}
	l := logl.New(c)
	use(l)
	l.Errorf("my error no %d", 78)
}

func exampleLogOff() {
	c := logl.Config{
		Handler: logl.StreamHandler(os.Stdout, "", logl.Ltime),
		Level:   logl.LevelOff,
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
		Handler: logl.StreamHandler(bw, "test ", logl.Ldate|logl.Lmicroseconds),
		Level:   logl.LevelInfo,
	}
	logger := logl.New(c)
	use(logger)
}

func examplePanicRecover() {

	defer func() {
		message := recover()
		if message != nil {
			fmt.Println("defer panic:", message)
		}
	}()

	sh := logl.StreamHandler(os.Stdout, "", logl.Ltime)
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
