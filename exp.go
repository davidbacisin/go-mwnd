package mwnd

// ExponentialWindow computes exponentially weighted moving window statistics
// over the input stream.
//
// All operations on the ExponentialWindow take constant time. Compared to FixedWindow,
// the exponential window has the advantages of constant time Puts
// and less memory usage, since it does not store any input values.
// However, exponential windows provide fundamentally different calculations from fixed
// windows—notably, infinite impulse response instead of finite impulse response—so
// the implementations are not generally interchangeable.
type ExponentialWindow[T Numeric] struct {
	// alpha is the weight of exponential fall-off
	alpha    float64
	mean, m2 float64
	min, max T
	size     int
}

// enforce compliance with interface
var _ window[float64] = (*ExponentialWindow[float64])(nil)

// Exponential initializes a moving window with the provided weight alpha.
func Exponential[T Numeric](alpha float64) *ExponentialWindow[T] {
	return &ExponentialWindow[T]{
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
func (w *ExponentialWindow[T]) Size() int {
	return w.size
}

// Min returns the lowest value ever observed by the Window.
// If the Window has no values, then it returns the zero value.
//
// Time complexity of O(1).
func (w *ExponentialWindow[T]) Min() T {
	return w.min
}

// Max returns the highest value ever observed by the Window.
// If the Window has no values, then it returns the zero value.
//
// Time complexity of O(1).
func (w *ExponentialWindow[T]) Max() T {
	return w.max
}

// Mean returns the exponentially-weighted moving average of all
// values ever added to the Window. If the Window has no values,
// then it returns the zero value.
//
// Time complexity of O(1).
func (w *ExponentialWindow[T]) Mean() float64 {
	return w.mean
}

// Variance returns the exponentially-weighted moving variance of
// all values ever added to the window. If the Window has no values,
// then it returns the zero value.
//
// Time complexity of O(1).
func (w *ExponentialWindow[T]) Variance() float64 {
	if w.size == 0 {
		return 0
	}
	return w.m2 / float64(w.size)
}

// Put adds a new value to the Window.
//
// Time complexity of O(1).
func (w *ExponentialWindow[T]) Put(v T) {
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
