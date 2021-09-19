// settings displays human readable device settings.
package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"

	"github.com/tormoder/fit"
)

// help writes help text.
// If no writer is provided, it writes to stderr.
func help(w io.Writer) {
	if w == nil {
		w = os.Stderr
	}
	fmt.Fprintf(w, `Displays human readable device settings.

Usage:
  %s <vivosport data directory>
`, os.Args[0])
}

func main() {
	if len(os.Args) != 2 {
		help(os.Stderr)
		os.Exit(1)
	}
	b, err := ioutil.ReadFile(filepath.Join(os.Args[1], "SETTINGS", "SETTINGS.FIT"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read file: %v\n", err)
		os.Exit(1)
	}

	data, err := fit.Decode(bytes.NewReader(b))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to decode FIT data: %v\n", err)
		os.Exit(1)
	}

	settings, err := data.Settings()
	if err != nil {
		fmt.Fprintf(os.Stderr, "FIT data was not settings: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("field\tvalue\n")
	for _, setting := range settings.DeviceSettings {
		v := reflect.ValueOf(*setting)
		for i := 0; i < v.NumField(); i++ {
			fmt.Printf("%s\t%v\n", v.Type().Field(i).Name, v.Field(i))
		}
	}
}
