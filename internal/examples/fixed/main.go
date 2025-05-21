// fixed/main generates sample data and plots it with its [mwnd.Fixed]
// statistics using gnuplot. The gnuplot executable must be available
// in the PATH for this example to work.
package main

import (
	"fmt"
	"io"
	"log/slog"
	"math"
	"os"
	"os/exec"
	"strconv"

	"github.com/davidbacisin/go-mwnd"
)

func main() {
	n := 500
	w := mwnd.Fixed[float64](n)

	out, err := os.OpenFile("data.csv", os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		slog.Error("failed to open file", "error", err)
		os.Exit(1)
	}
	defer out.Close()

	fmt.Fprintln(out, "x y min max mean")
	for i := 0.0; i < 300*math.Pi; i += 0.1 {
		v := 5*math.Sin(i*0.01) + math.Sin(i*0.1)
		w.Put(v)
		writeLine(out, i, v, w.Min(), w.Max(), w.Mean())
	}

	cmd := exec.Command("gnuplot", "-e", `
	set terminal pngcairo size 2000,1000 linewidth 4;
	set output 'plot.png';
	set datafile columnheaders;
	set for [i=1:8] linetype i dashtype i;
	plot [0:943] 'data.csv' using "x":"y" with lines title 'samples',
		'data.csv' using "x":"min" with lines title 'min',
		'data.csv' using "x":"max" with lines title 'max',
		'data.csv' using "x":"mean" with lines title 'mean'
	`)
	if err := cmd.Run(); err != nil {
		slog.Error("failed to run gnuplot", "error", err)
		os.Exit(1)
	}
}

func writeLine(out io.Writer, values ...float64) {
	if len(values) == 0 {
		return
	}

	fmt.Fprint(out, strconv.FormatFloat(values[0], 'f', 1, 64))
	for i := 1; i < len(values); i++ {
		fmt.Fprint(out, " ", strconv.FormatFloat(values[i], 'f', 5, 64))
	}

	fmt.Fprintln(out)
}
