package mwnd_test

import (
	"math/rand/v2"
	"testing"

	"github.com/davidbacisin/go-mwnd"
)

func BenchmarkFixed_1000(b *testing.B) {
	w := mwnd.Fixed[int](1000)
	for b.Loop() {
		v := rand.Int()
		w.Put(v)

		if w.Min() < 0 || w.Max() < 0 || w.Mean() < 0.0 {
			b.Logf("invalid min, max, or mean")
			b.FailNow()
		}
	}
}

func BenchmarkExponential(b *testing.B) {
	w := mwnd.Exponential[int](0.002)
	for b.Loop() {
		v := rand.Int()
		w.Put(v)

		if w.Min() < 0 || w.Max() < 0 || w.Mean() < 0.0 {
			b.Logf("invalid min, max, or mean")
			b.FailNow()
		}
	}
}

func BenchmarkFixed_1000_Quantiles(b *testing.B) {
	cases := []struct {
		name string
		q    float64
	}{
		{name: "first percentile", q: 0.01},
		{name: "first decile", q: 0.1},
		{name: "first quartile", q: 0.25},
		{name: "median", q: 0.5},
		{name: "third quartile", q: 0.75},
		{name: "ninth decile", q: 0.9},
		{name: "99th percentile", q: 0.99},
	}

	for _, c := range cases {
		b.Run(c.name, func(b *testing.B) {
			w := mwnd.Fixed[int](1000)
			for b.Loop() {
				v := rand.Int()
				w.Put(v)

				if w.Quantile(c.q) < 0 {
					b.Logf("invalid quantile")
					b.FailNow()
				}
			}
		})
	}
}
