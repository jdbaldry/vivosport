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

	b, err := ioutil.ReadFile(filepath.Join(os.Args[1], "DEVICE.FIT"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read file: %v\n", err)
		os.Exit(1)
	}

	data, err := fit.Decode(bytes.NewReader(b))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to decode FIT data: %v\n", err)
		os.Exit(1)
	}

	device, err := data.Device()
	if err != nil {
		fmt.Fprintf(os.Stderr, "FIT data was not device information: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("type\tfield\tvalue\n")
	for _, capability := range device.FieldCapabilities {
		v := reflect.ValueOf(*capability)
		for i := 0; i < v.NumField(); i++ {
			fmt.Printf("field capability\t%s\t%v\n", v.Type().Field(i).Name, v.Field(i))
		}
	}
	for _, capability := range device.FileCapabilities {
		v := reflect.ValueOf(*capability)
		for i := 0; i < v.NumField(); i++ {
			fmt.Printf("file capability\t%s\t%v\n", v.Type().Field(i).Name, v.Field(i))
		}
	}
	for _, capability := range device.MesgCapabilities {
		v := reflect.ValueOf(*capability)
		for i := 0; i < v.NumField(); i++ {
			fmt.Printf("message capability\t%s\t%v\n", v.Type().Field(i).Name, v.Field(i))
		}
	}
	for _, software := range device.Softwares {
		v := reflect.ValueOf(*software)
		for i := 0; i < v.NumField(); i++ {
			fmt.Printf("software\t%s\t%v\n", v.Type().Field(i).Name, v.Field(i))
		}
	}
	for _, capability := range device.Capabilities {
		v := reflect.ValueOf(*capability)
		for i := 0; i < v.NumField(); i++ {
			fmt.Printf("device capability\t%s\t%v\n", v.Type().Field(i).Name, v.Field(i))
		}
	}

	b, err = ioutil.ReadFile(filepath.Join(os.Args[1], "SETTINGS", "SETTINGS.FIT"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read file: %v\n", err)
		os.Exit(1)
	}

	data, err = fit.Decode(bytes.NewReader(b))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to decode FIT data: %v\n", err)
		os.Exit(1)
	}

	settings, err := data.Settings()
	if err != nil {
		fmt.Fprintf(os.Stderr, "FIT data was not settings: %v\n", err)
		os.Exit(1)
	}

	for _, setting := range settings.DeviceSettings {
		v := reflect.ValueOf(*setting)
		for i := 0; i < v.NumField(); i++ {
			fmt.Printf("setting\t%s\t%v\n", v.Type().Field(i).Name, v.Field(i))
		}
	}
}
