package mwnd_test

import (
	"fmt"

	"github.com/davidbacisin/go-mwnd"
)

func ExampleFixed() {
	w := mwnd.Fixed[int](5)

	// Note that len(values) > w.Size(), so the first value (1) will be evicted
	// when the last value (10) is Put.
	values := []int{1, 5, 4, 3, 2, 10}
	for _, v := range values {
		w.Put(v)
	}

	fmt.Printf("Size: %d\n", w.Size())
	fmt.Printf("Min: %d\n", w.Min())
	fmt.Printf("Max: %d\n", w.Max())
	fmt.Printf("Mean: %.2f\n", w.Mean())
	fmt.Printf("Variance: %.2f\n", w.Variance())

	// Output:
	// Size: 5
	// Min: 2
	// Max: 10
	// Mean: 4.80
	// Variance: 7.76
}

func ExampleExponential() {
	w := mwnd.Exponential[int](0.1)

	values := []int{1, 5, 4, 3, 2, 10}
	for _, v := range values {
		w.Put(v)
	}

	fmt.Printf("Size: %d\n", w.Size())
	fmt.Printf("Min: %d\n", w.Min())
	fmt.Printf("Max: %d\n", w.Max())
	fmt.Printf("Mean: %.2f\n", w.Mean())
	fmt.Printf("Variance: %.2f\n", w.Variance())

	// Output:
	// Size: 6
	// Min: 1
	// Max: 10
	// Mean: 2.63
	// Variance: 13.74
}
