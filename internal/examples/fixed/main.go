// fixed/main generates sample data and plots it with its [mwnd.Fixed]
// statistics using gnuplot. The gnuplot executable must be available
// in the PATH for this example to work.
package main

import (
	"fmt"
	"log/slog"
	"math"
	"os"
	"os/exec"

	"github.com/davidbacisin/go-mwnd"
	"github.com/davidbacisin/go-mwnd/internal/examples"
)

func main() {
	n := 500
	w := mwnd.Fixed[float64](n)

	basic := examples.OpenOutputFile("basic.csv")
	defer basic.Close()
	fmt.Fprintln(basic, "x y min max mean var")

	quantiles := examples.OpenOutputFile("quantiles.csv")
	defer quantiles.Close()
	fmt.Fprintln(quantiles, "x y p10 p25 p50 p75 p90 p99")

	for i := 0.0; i < 300*math.Pi; i += 0.1 {
		v := 5*math.Sin(i*0.01) + math.Sin(i*0.1)
		w.Put(v)
		examples.WriteLine(basic, i, v, w.Min(), w.Max(), w.Mean(), w.Variance())
		examples.WriteLine(quantiles, i, v, w.Quantile(0.1), w.Quantile(0.25), w.Quantile(0.5), w.Quantile(0.75), w.Quantile(0.9), w.Quantile(0.99))
	}

	// Plot basic data
	cmd := exec.Command("gnuplot", "-e", `
	set terminal pngcairo size 2000,1000 linewidth 4;
	set output 'plot.png';
	set datafile columnheaders;
	set for [i=1:8] linetype i dashtype i;
	plot [0:943] 'basic.csv' using "x":"y" with lines title 'samples',
		'basic.csv' using "x":"min" with lines title 'min',
		'basic.csv' using "x":"max" with lines title 'max',
		'basic.csv' using "x":"mean" with lines title 'mean',
		'basic.csv' using "x":"var" with lines title 'var'
	`)
	if err := cmd.Run(); err != nil {
		slog.Error("failed to run gnuplot", "error", err)
		os.Exit(1)
	}

	// Plot quantiles
	cmd = exec.Command("gnuplot", "-e", `
	set terminal pngcairo size 2000,1000 linewidth 4;
	set output 'quantiles.png';
	set datafile columnheaders;
	set for [i=1:8] linetype i dashtype i;
	plot [0:943] 'quantiles.csv' using "x":"y" with lines title 'samples',
		'quantiles.csv' using "x":"p10" with lines title 'p10',
		'quantiles.csv' using "x":"p25" with lines title 'p25',
		'quantiles.csv' using "x":"p50" with lines title 'p50',
		'quantiles.csv' using "x":"p75" with lines title 'p75',
		'quantiles.csv' using "x":"p90" with lines title 'p90',
		'quantiles.csv' using "x":"p99" with lines title 'p99'
	`)
	if err := cmd.Run(); err != nil {
		slog.Error("failed to run gnuplot", "error", err)
		os.Exit(1)
	}
}
