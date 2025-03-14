package mwnd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Mean(t *testing.T) {
	samples := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	wantMean := 5.5

	t.Run("smaller sample length", func(t *testing.T) {
		s := New(len(samples) + 1)
		for _, v := range samples {
			s.Record(v)
		}

		gotMean := s.Mean()
		assert.Equal(t, wantMean, gotMean)
	})

	t.Run("exact sample length", func(t *testing.T) {
		s := New(len(samples))
		for _, v := range samples {
			s.Record(v)
		}

		gotMean := s.Mean()
		assert.Equal(t, wantMean, gotMean)
	})

	t.Run("twice sample length", func(t *testing.T) {
		// should only use the second half of the sample set
		wantMean := float64(6+7+8+9+10) / 5

		s := New(len(samples) / 2)
		for _, v := range samples {
			s.Record(v)
		}

		gotMean := s.Mean()
		assert.Equal(t, wantMean, gotMean)
	})
}
