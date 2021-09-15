package main

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

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
  %s <vivosport data directory>
`, os.Args[0])
}

func main() {
	if len(os.Args) != 2 {
		help(os.Stderr)
		os.Exit(1)
	}

	dir, err := filepath.Abs(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to determine path to vivosport data: %v\n", err)
		os.Exit(1)
	}

	sessions := []*fit.SessionMsg{}
	activitiesDir := filepath.Join(dir, "ACTIVITY")
	filepath.WalkDir(activitiesDir, func(path string, d fs.DirEntry, err error) error {
		if path == activitiesDir {
			return nil
		}

		b, err := ioutil.ReadFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to read file: %v\n", err)
			os.Exit(1)
		}

		data, err := fit.Decode(bytes.NewReader(b))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to decode FIT data: %v\n", err)
			os.Exit(1)
		}

		activity, err := data.Activity()
		if err != nil {
			fmt.Fprintf(os.Stderr, "FIT data was not an activity: %v\n", err)
			os.Exit(1)
		}

		sessions = append(sessions, activity.Sessions...)
		return nil
	})

	sort.Slice(sessions, func(p, q int) bool {
		return sessions[p].Timestamp.Unix() > sessions[q].Timestamp.Unix()
	})

	type column struct {
		title  string
		format string
		value  func(s *fit.SessionMsg) interface{}
	}

	fmt.Println("Sessions")
	columns := []column{
		{"start", "%s", func(s *fit.SessionMsg) interface{} { return s.StartTime }},
		{"end", "%s", func(s *fit.SessionMsg) interface{} { return s.Timestamp }},
		{"time", "%.02fmins", func(s *fit.SessionMsg) interface{} { return s.GetTotalTimerTimeScaled() / 60 }},
		{"sport", "%s", func(s *fit.SessionMsg) interface{} { return s.Sport }},
		{"HR (avg)", "%d", func(s *fit.SessionMsg) interface{} { return s.AvgHeartRate }},
		{"HR (max)", "%d", func(s *fit.SessionMsg) interface{} { return s.MaxHeartRate }},
		{"speed (avg)", "%.02fm/s", func(s *fit.SessionMsg) interface{} { return s.GetAvgSpeedScaled() }},
		{"distance", "%.02fm", func(s *fit.SessionMsg) interface{} { return s.GetTotalDistanceScaled() }},
		{"calories", "%d", func(s *fit.SessionMsg) interface{} { return s.TotalCalories }},
	}
	for i, c := range columns {
		if i > 0 {
			fmt.Printf("\t")
		}
		fmt.Printf("%s", c.title)
	}
	fmt.Printf("\n")
	for _, s := range sessions {
		for i, c := range columns {
			if i > 0 {
				fmt.Printf("\t")
			}
			fmt.Printf(c.format, c.value(s))
		}
		fmt.Printf("\n")
	}
}
