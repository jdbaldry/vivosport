package main

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/jdbaldry/vivosport/pgsql"
	_ "github.com/lib/pq"
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

	ctx := context.Background()

	db, err := sql.Open("postgres", "dbname=vivosport password=vivosport sslmode=disable user=vivosport")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to open DB connection: %v", err)
		os.Exit(1)
	}

	queries := pgsql.New(db)

	dir, err := filepath.Abs(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to determine path to vivosport data: %v\n", err)
		os.Exit(1)
	}

	activitiesDir := filepath.Join(dir, "ACTIVITY")
	err = filepath.WalkDir(activitiesDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == activitiesDir {
			return nil
		}

		b, err := ioutil.ReadFile(path)
		if err != nil {
			return fmt.Errorf("unable to read file: %w\n", err)
		}

		data, err := fit.Decode(bytes.NewReader(b))
		if err != nil {
			return fmt.Errorf("unable to decode FIT data: %w\n", err)
		}

		activityFile, err := data.Activity()
		if err != nil {
			return fmt.Errorf("FIT data was not an activity: %w\n", err)
		}
		inserted, err := queries.CreateActivity(ctx, pgsql.CreateActivityParams{
			StartTs:        activityFile.Activity.Timestamp.Add(-(time.Duration(activityFile.Activity.GetTotalTimerTimeScaled()) * time.Second)),
			EndTs:          activityFile.Activity.Timestamp,
			TotalTimerTime: sql.NullInt32{Int32: int32(activityFile.Activity.GetTotalTimerTimeScaled()), Valid: true},
			NumSessions:    sql.NullInt32{Int32: int32(activityFile.Activity.NumSessions), Valid: true},
			Type:           sql.NullInt32{Int32: int32(activityFile.Activity.Type), Valid: true},
			Event:          int16(activityFile.Activity.Event),
			EventType:      int16(activityFile.Activity.EventType),
			LocalTs:        sql.NullTime{Time: activityFile.Activity.LocalTimestamp, Valid: true},
			EventGroup:     int16(activityFile.Activity.EventGroup),
		})
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil
			}
			return fmt.Errorf("failed to create activity: %w\n", err)
		}
		fmt.Printf("Created activity: %#v\n", inserted)

		return nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Trouble walking directories: %v\n", err)
		os.Exit(1)
	}
}
