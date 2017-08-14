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
		switch c := r.Intn(4); c {
		case 0:
			l.Debug(line)
		case 1:
			l.Info(line)
		case 2:
			l.Warning(line)
		case 3:
			l.Error(line)
		}
	}
}

func exampleLogStdout() {
	c := logl.Config{
		Output:  logl.OnlyWriter(os.Stdout),
		Prefix:  "",
		Level:   logl.LEVEL_DEBUG,
		Flag:    logl.Ltime,
		NotSafe: false,
	}
	l := logl.New(c)
	use(l)
}

func exampleLogOff() {
	c := logl.Config{
		Output:  logl.OnlyWriter(os.Stdout),
		Prefix:  "",
		Level:   logl.LEVEL_OFF,
		Flag:    logl.Ltime,
		NotSafe: false,
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
		Output:  bw,
		Prefix:  "test ",
		Level:   logl.LEVEL_INFO,
		Flag:    logl.Ldate | logl.Lmicroseconds,
		NotSafe: true,
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
	c := logl.Config{
		Output:  logl.OnlyWriter(os.Stdout),
		Prefix:  "",
		Level:   logl.LEVEL_ERROR,
		Flag:    logl.Ltime,
		NotSafe: false,
	}
	l := logl.New(c)
	l.Panic("my panic message")
}
