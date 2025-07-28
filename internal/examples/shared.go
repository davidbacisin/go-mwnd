package examples

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"strconv"
)

func OpenOutputFile(name string) *os.File {
	out, err := os.OpenFile(name, os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		slog.Error("failed to open file", "error", err)
		os.Exit(1)
	}
	return out
}

func WriteLine(out io.Writer, values ...float64) {
	if len(values) == 0 {
		return
	}

	fmt.Fprint(out, strconv.FormatFloat(values[0], 'f', 1, 64))
	for i := 1; i < len(values); i++ {
		fmt.Fprint(out, " ", strconv.FormatFloat(values[i], 'f', 5, 64))
	}

	fmt.Fprintln(out)
}
