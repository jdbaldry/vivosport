package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/tormoder/fit"
)

// help writes help text.
// If no writer is provided, it writes to stderr.
func help(w io.Writer) {
	if w == nil {
		w = os.Stderr
	}
	fmt.Fprintf(w, `A tool for understanding FIT data.

Usage:
  %s <file>
`, os.Args[0])
}

func main() {
	if len(os.Args) != 2 {
		help(os.Stderr)
		os.Exit(1)
	}

	f, err := filepath.Abs(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to determine path to file: %v", err)
		os.Exit(1)
	}

	b, err := ioutil.ReadFile(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read file: %v", err)
		os.Exit(1)
	}

	data, err := fit.Decode(bytes.NewReader(b))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to decode FIT data: %v", err)
		os.Exit(1)
	}

	fmt.Printf("%+v\n", data)
}
