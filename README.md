# go-mwnd
[![codecov](https://codecov.io/gh/davidbacisin/go-mwnd/graph/badge.svg?token=439YYC0C5O)](https://codecov.io/gh/davidbacisin/go-mwnd) [![buy me a coffee](https://img.shields.io/badge/%E2%80%8B_buy_me_a_coffee-fd0?logo=buymeacoffee&logoColor=333)](https://buymeacoffee.com/davidbacisin)

Moving window statistics for Go. Computes mean, minimum, maximum, and population 
variance over a sliding window, supporting both fixed-size and 
exponentially-weighted windows.

## Usage 🚀
```go
// Define a fixed moving window that holds up to 5 values
w := mwnd.Fixed[int](5)

// Note that len(values) > w.Size(), so when the last value (10) is Put,
// then the first value (1) will be evicted.
values := []int{1, 5, 4, 3, 2, 10}
for _, v := range values {
    w.Put(v)
}

fmt.Printf("Size: %d\n", w.Size())
fmt.Printf("Min: %d\n", w.Min())
fmt.Printf("Max: %d\n", w.Max())
fmt.Printf("Mean: %.2f\n", w.Mean())
fmt.Printf("Variance: %.2f\n", w.Variance())
```

## Example visualization 📊
The following plot shows generated data (solid purple line) and the moving min, max, mean, and variance 
for a fixed window of 500 values:
![plot of a composite sine wave and its fixed moving window statistics](internal/examples/fixed/plot.png)

The following plot shows the same data with statistics from an exponentially-weighted moving window 
with alpha = 0.004:
![plot of a composite sine wave and its exponential moving window statistics](internal/examples/exponential/plot.png)

Note especially the differences in min, max, and variance between these two implementations: as a result,
they are not interchangeable and instead can only approximate each other.

## Limitations and Future Work 🧪
- The moving window implementations are not safe for concurrent reads or writes. Consider using 
[`sync.Mutex`](https://pkg.go.dev/sync#Mutex) or [`sync.RWMutex`](https://pkg.go.dev/sync#RWMutex) 
to synchronize access under a concurrent workload.
- **Random sampling**: In some cases, it may be impractical to include every value 
in the window. Instead, a random subset could be probabilistically sampled. The variance 
calculation would need to be corrected for bias.
