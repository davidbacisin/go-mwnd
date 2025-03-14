package mwnd

type Samples struct {
	mean           float64
	size, i, count int
	ring           []float64
}

func New(size int) *Samples {
	return &Samples{
		size: size,
		ring: make([]float64, size),
	}
}

func (s *Samples) Count() int {
	return s.count
}

func (s *Samples) Record(v float64) {
	old := s.ring[s.i]
	s.ring[s.i] = v
	s.i = (s.i + 1) % s.size
	if s.count == s.size {
		delta := v - old
		s.mean += delta / float64(s.count)
	} else {
		s.count++
		delta := v - s.mean
		s.mean += delta / float64(s.count)
	}
}

func (s *Samples) Mean() float64 {
	return s.mean
}
