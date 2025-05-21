package mwnd_test

import (
	"math/rand/v2"
	"testing"

	"github.com/davidbacisin/go-mwnd"
)

func BenchmarkWindow_1000(b *testing.B) {
	w := mwnd.Fixed[int](1000)
	for range b.N {
		v := rand.Int()
		w.Put(v)

		if w.Min() < 0 || w.Max() < 0 || w.Mean() < 0.0 {
			b.Logf("invalid min, max, or mean")
			b.FailNow()
		}
	}
}
