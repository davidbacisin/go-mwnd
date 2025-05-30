package mwnd

type Window[T Numeric] interface {
	Size() int
	Put(T)
	Min() T
	Max() T
	Mean() float64
	Variance() float64
}
