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
	w := mwnd.Exponential[float64](0.004)

	out := examples.OpenOutputFile("data.csv")
	defer out.Close()

	fmt.Fprintln(out, "x y min max mean var")
	for i := 0.0; i < 300*math.Pi; i += 0.1 {
		v := 5*math.Sin(i*0.01) + math.Sin(i*0.1)
		w.Put(v)
		examples.WriteLine(out, i, v, w.Min(), w.Max(), w.Mean(), w.Variance())
	}

	cmd := exec.Command("gnuplot", "-e", `
	set terminal pngcairo size 2000,1000 linewidth 4;
	set output 'plot.png';
	set datafile columnheaders;
	set for [i=1:8] linetype i dashtype i;
	plot [0:943] 'data.csv' using "x":"y" with lines title 'samples',
		'data.csv' using "x":"min" with lines title 'min',
		'data.csv' using "x":"max" with lines title 'max',
		'data.csv' using "x":"mean" with lines title 'mean',
		'data.csv' using "x":"var" with lines title 'var'
	`)
	if err := cmd.Run(); err != nil {
		slog.Error("failed to run gnuplot", "error", err)
		os.Exit(1)
	}
}
