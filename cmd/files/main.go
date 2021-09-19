// files traverses a garmin data directory and lists all the files and their
// FIT file types.
package main

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/tormoder/fit"
)

// help writes help text.
// If no writer is provided, it writes to stderr.
func help(w io.Writer) {
	if w == nil {
		w = os.Stderr
	}
	fmt.Fprintf(w, `Traverses a garmin data directory and lists all the files and their FIT file types.

Usage:
  %s <vivosport data directory>
`, os.Args[0])
}

func main() {
	if len(os.Args) != 2 {
		help(os.Stderr)
		os.Exit(1)
	}

	fmt.Printf("Type\tPath\n")
	err := filepath.WalkDir(os.Args[1], func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return err
		}

		if d.IsDir() {
			fmt.Printf("DIR\t%s\n", path)
			if d.Name() == "METRICS" || d.Name() == "RECORDS" || d.Name() == "SLEEP" {
				err := filepath.WalkDir(filepath.Join(os.Args[1], d.Name()), func(path string, d fs.DirEntry, err error) error {
					if strings.HasSuffix(d.Name(), ".FIT") {
						b, err := ioutil.ReadFile(path)
						if err != nil {
							return fmt.Errorf("unable to read file %s: %w\n", path, err)
						}

						_, err = fit.Decode(bytes.NewReader(b))
						if err != nil {
							fmt.Fprintf(os.Stderr, "Unable to decode FIT data from file %s: %v\n", path, err)
						}
					}
					return nil
				})
				if err != nil {
					return err
				}
				return filepath.SkipDir
			}

			return nil
		}

		if strings.HasSuffix(d.Name(), ".FIT") {
			b, err := ioutil.ReadFile(path)
			if err != nil {
				return fmt.Errorf("unable to read file %s: %w\n", path, err)
			}

			data, err := fit.Decode(bytes.NewReader(b))
			if err != nil {
				return fmt.Errorf("unable to decode FIT data from file %s: %w\n", path, err)
			}

			fmt.Printf("%s\t%s\n", data.FileId.Type, path)
		}
		return nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Trouble walking directories: %v\n", err)
		os.Exit(1)
	}
}
