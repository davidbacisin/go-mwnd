package mwnd

// exponential computes exponentially weighted moving window statistics
// over the input stream.
type exponential[T Numeric] struct {
	// alpha is the weight of exponential fall-off
	alpha    float64
	mean, m2 float64
	min, max T
	size     int
}

// enforce compliance with interface
var _ Window[float64] = (*exponential[float64])(nil)

// Exponential initializes a moving window with the provided weight alpha.
func Exponential[T Numeric](alpha float64) *exponential[T] {
	return &exponential[T]{
		alpha: alpha,
		size:  0,
	}
}

// ExponentialAlphaForApproximatingFixed returns an alpha value for an exponential moving window
// that will approximate the behavior of a fixed moving window of length n.
func ExponentialAlphaForApproximatingFixed(n int) float64 {
	return 2.0 / float64(n+1)
}

// Size returns the number of values added to the Window.
func (w *exponential[T]) Size() int {
	return w.size
}

// Min returns the lowest value ever observed by the Window.
// If the Window has no values, then it returns the zero value.
//
// Time complexity of O(1).
func (w *exponential[T]) Min() T {
	return w.min
}

// Max returns the highest value ever observed by the Window.
// If the Window has no values, then it returns the zero value.
//
// Time complexity of O(1).
func (w *exponential[T]) Max() T {
	return w.max
}

// Mean returns the exponentially-weighted moving average of all
// values ever added to the Window. If the Window has no values,
// then it returns the zero value.
//
// Time complexity of O(1).
func (w *exponential[T]) Mean() float64 {
	return w.mean
}

// Variance returns the exponentially-weighted moving variance of
// all values ever added to the window. If the Window has no values,
// then it returns the zero value.
//
// Time complexity of O(1).
func (w *exponential[T]) Variance() float64 {
	if w.size == 0 {
		return 0
	}
	return w.m2 / float64(w.size)
}

// Put adds a new value to the Window.
//
// Time complexity of O(1).
func (w *exponential[T]) Put(v T) {
	w.size++
	if w.size == 1 {
		w.mean = float64(v)
		w.min = v
		w.max = v
		w.m2 = 0.0
		return
	}

	// Welford's algorithm for online variance, which is a numerically stable approach.
	delta := float64(v) - w.mean
	w.mean = w.alpha*float64(v) + (1-w.alpha)*w.mean
	delta2 := float64(v) - w.mean
	w.m2 += delta * delta2

	w.min = min(w.min, v)
	w.max = max(w.max, v)
}
