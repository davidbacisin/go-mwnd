// Package mwnd provides a high-performance on-line moving window implementation for
// computing the minimum, maximum, mean, and total sum of squared differences from
// the mean over a stream of numeric values.
//
// The worst-case time complexity of adding a new value to the window is O(log n),
// where n is the capacity of the window. This is achieved by maintaining a red-black
// balanced binary tree as the underlying data structure. Furthermore, the Put and
// all statistical operations avoid memory allocations to preserve performance in
// high-scale environments.
package mwnd
