package main

import (
	"bytes"
	"math/rand"
	"time"

	"github.com/gitchander/logl"
)

func newRandSeed(seed int64) *rand.Rand {
	return rand.New(rand.NewSource(seed))
}

func newRandTime(t time.Time) *rand.Rand {
	return newRandSeed(t.UnixNano())
}

func newRandNow() *rand.Rand {
	return newRandTime(time.Now())
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

var levels = []logl.Level{
	logl.LevelCritical,
	logl.LevelError,
	logl.LevelWarning,
	logl.LevelInfo,
	logl.LevelDebug,
	logl.LevelTrace,
}

func randLogLevel(r *rand.Rand) logl.Level {
	return levels[r.Intn(len(levels))]
}

func randLogMessage(r *rand.Rand, l logl.Logger) {
	var (
		level   = randLogLevel(r)
		message = randLine(r, randIntRange(r, 3, 10))
	)
	l.Log(level, message)
}
