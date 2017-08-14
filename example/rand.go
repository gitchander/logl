package main

import (
	"bytes"
	"math/rand"
	"time"
)

func newRand() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
}

func randLine(r *rand.Rand, n_word int) string {
	delims := []byte(",.;:!?")
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
