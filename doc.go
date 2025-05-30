// Package mwnd provides high-performance on-line moving window implementations for
// computing the minimum, maximum, mean, and variance over a stream of numeric values.
//
// For the fixed-size moving window implementation, the worst-case time complexity of adding
// a new value to the window is O(log n), where n is the capacity of the window. This
// is achieved by maintaining a red-black balanced binary tree as the underlying data
// structure. Furthermore, the Put and all statistical operations avoid memory allocations
// to preserve performance in high-scale environments.
//
// For the exponentially weighted moving window implementation, all operations are
// constant time. Compared to the fixed window, the exponential window has the advantages
// of constant time Puts and less memory usage, since it does not store any input values.
// However, exponential windows provide fundamentally different calculations from fixed
// windows—notably, infinite impulse response instead of finite impulse response—so
// the implementations are not generally interchangeable.
package mwnd
