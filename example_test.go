package mwnd_test

import (
	"fmt"

	"github.com/davidbacisin/go-mwnd"
)

func ExampleWindow() {
	w := mwnd.New[int](5)

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
	fmt.Printf("Total Sum of Squares: %.2f\n", w.TotalSumSquares())

	// Output:
	// Size: 5
	// Min: 2
	// Max: 10
	// Mean: 4.80
	// Total Sum of Squares: 38.80
}
