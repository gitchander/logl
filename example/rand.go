package main

import (
	"bytes"
	"math/rand"
	"time"

	"github.com/gitchander/logl"
)

func newRand() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
}

func newRandSeed(seed int64) *rand.Rand {
	return rand.New(rand.NewSource(seed))
}

func randLine(r *rand.Rand, n_word int) string {
	delims := []byte{
		'\n', '\r', '\t', '\\', '"', '\'', ',', '.', ';', ':', '!', '?',
	}
	var buf bytes.Buffer
	for i := 0; i < n_word; i++ {
		if i > 0 {
			if r.Intn(100) < 30 {
				buf.WriteByte(delims[r.Intn(len(delims))])
			}
			buf.WriteByte(' ') // space
		}
		n_letter := randIntRange(r, 3, 7)
		randWord(r, n_letter, &buf)
	}
	return buf.String()
}

func randWord(r *rand.Rand, n_letter int, buf *bytes.Buffer) {
	for i := 0; i < n_letter; i++ {
		b := byte(randIntRange(r, 'a', 'z'))
		buf.WriteByte(b)
	}
}

func randIntRange(r *rand.Rand, min, max int) int {
	if min > max {
		min, max = max, min
	}
	return min + r.Intn(max-min+1)
}

func randLevel(r *rand.Rand) (level logl.Level) {
	switch k := r.Intn(5); k {
	case 0:
		level = logl.LevelCritical
	case 1:
		level = logl.LevelError
	case 2:
		level = logl.LevelWarning
	case 3:
		level = logl.LevelInfo
	case 4:
		level = logl.LevelDebug
	}
	return
}

func randLogMessage(r *rand.Rand, l *logl.Logger) {
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
