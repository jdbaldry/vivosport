// db updates a local postgres database with FIT file data.
package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
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
	fmt.Fprintf(w, `Updates a local postgres database with FIT file data.

Usage:
  %s <vivosport data directory> <csv data directory>
`, os.Args[0])
}

func main() {
	if len(os.Args) != 3 {
		help(os.Stderr)
		os.Exit(1)
	}

	var activityCount, monitoringCount, recordCount, sessionCount int
	var activityInserted, monitoringInserted, recordInserted, sessionInserted int

	ctx := context.Background()

	db, err := sql.Open("postgres", "dbname=vivosport password=vivosport sslmode=disable user=vivosport")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to open DB connection: %v", err)
		os.Exit(1)
	}

	queries := pgsql.New(db)

	fitDir, err := filepath.Abs(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to determine path to vivosport data: %v\n", err)
		os.Exit(1)
	}

	csvDir, err := filepath.Abs(os.Args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to determine path to csv data: %v\n", err)
		os.Exit(1)
	}

	path := filepath.Join(csvDir, "RECORDS", "RECORDS.csv")
	f, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to open file %s: %v", path, err)
		os.Exit(1)
	}
	r := csv.NewReader(f)
	r.FieldsPerRecord = -1 // Records have variable length fields
	records, err := r.ReadAll()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read file %s as CSV: %v", path, err)
		os.Exit(1)
	}
	for i, record := range records {
		if len(record) >= 17 {
			// 7th field is distance. [[file:~/ext/jdbaldry/vivosport/csv/RECORDS/RECORDS.md::Records][Records]]
			if record[7] == "100000" || record[7] == "160900" || record[7] == "500000" {
				recordCount++
				distance, err := strconv.ParseInt(record[7], 10, 32)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Unable to convert CSV distance field in record %d to int: %v", i, err)
					os.Exit(1)
				}
				time, err := strconv.ParseInt(record[16], 10, 32)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Unable to convert CSV time field in record %d to int: %v", i, err)
					os.Exit(1)
				}

				_, err = queries.CreateRecord(ctx, pgsql.CreateRecordParams{
					Distance: sql.NullInt32{Int32: int32(distance), Valid: true},
					Time:     sql.NullInt32{Int32: int32(time), Valid: true},
				})

				if err != nil {
					if !errors.Is(err, sql.ErrNoRows) {
						fmt.Fprintf(os.Stderr, "failed to create record: %v\n", err)
						os.Exit(1)
					}
				}
				recordInserted++
			}
		}
	}

	monitoringDir := filepath.Join(fitDir, "MONITOR")
	err = filepath.WalkDir(monitoringDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == monitoringDir {
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

		monitoringBFile, err := data.MonitoringB()
		if err != nil {
			return fmt.Errorf("error extracting monitoring data: %w\n", err)
		}

		for _, monitoring := range monitoringBFile.Monitorings {
			monitoringCount++
			_, err := queries.CreateMonitoring(ctx, pgsql.CreateMonitoringParams{
				Ts:              monitoring.Timestamp,
				Calories:        int16(monitoring.Calories),
				Cycles:          sql.NullInt32{Int32: int32(monitoring.Cycles), Valid: true},
				Distance:        sql.NullFloat64{Float64: monitoring.GetDistanceScaled(), Valid: true},
				ActiveTime:      sql.NullFloat64{Float64: monitoring.GetActiveTimeScaled(), Valid: true},
				ActivityType:    int16(monitoring.ActivityType),
				ActivitySubType: int16(monitoring.ActivitySubtype),
				LocalTs:         sql.NullTime{Time: monitoring.LocalTimestamp, Valid: true},
			})
			if err != nil {
				if !errors.Is(err, sql.ErrNoRows) {
					return fmt.Errorf("failed to create monitoring: %w\n", err)
				}
			}
			monitoringInserted++
		}
		return nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Trouble walking directories: %v\n", err)
		os.Exit(1)
	}

	activitiesDir := filepath.Join(fitDir, "ACTIVITY")
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
		activityCount++
		_, err = queries.CreateActivity(ctx, pgsql.CreateActivityParams{
			StartTs:        activityFile.Activity.Timestamp.Add(-(time.Duration(activityFile.Activity.GetTotalTimerTimeScaled()) * time.Second)),
			EndTs:          activityFile.Activity.Timestamp,
			TotalTimerTime: sql.NullFloat64{Float64: activityFile.Activity.GetTotalTimerTimeScaled(), Valid: true},
			NumSessions:    sql.NullInt32{Int32: int32(activityFile.Activity.NumSessions), Valid: true},
			Type:           sql.NullInt32{Int32: int32(activityFile.Activity.Type), Valid: true},
			Event:          int16(activityFile.Activity.Event),
			EventType:      int16(activityFile.Activity.EventType),
			LocalTs:        sql.NullTime{Time: activityFile.Activity.LocalTimestamp, Valid: true},
			EventGroup:     int16(activityFile.Activity.EventGroup),
		})
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("failed to create activity: %w\n", err)
			}
		}
		activityInserted++

		for _, session := range activityFile.Sessions {
			sessionCount++
			_, err := queries.CreateSession(ctx, pgsql.CreateSessionParams{
				StartTs:          session.StartTime,
				EndTs:            session.Timestamp,
				Event:            int16(session.Event),
				EventType:        int16(session.EventType),
				Sport:            int16(session.Sport),
				SubSport:         int16(session.SubSport),
				TotalElapsedTime: sql.NullFloat64{Float64: session.GetTotalElapsedTimeScaled(), Valid: true},
				TotalTimerTime:   sql.NullFloat64{Float64: session.GetTotalTimerTimeScaled(), Valid: true},
				TotalDistance:    sql.NullFloat64{Float64: session.GetTotalDistanceScaled(), Valid: true},
				TotalCalories:    int16(session.TotalCalories),
				AvgSpeed:         sql.NullFloat64{Float64: session.GetAvgSpeedScaled(), Valid: true},
				MaxSpeed:         sql.NullFloat64{Float64: session.GetMaxSpeedScaled(), Valid: true},
				AvgHeartRate:     int16(session.AvgHeartRate),
				MaxHeartRate:     int16(session.MaxHeartRate),
			})
			if err != nil {
				if !errors.Is(err, sql.ErrNoRows) {
					return fmt.Errorf("failed to create session: %w\n", err)
				}
			}
			sessionInserted++
		}

		return nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Trouble walking directories: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("type\tcount\tinserted")
	fmt.Printf("%s\t%d\t%d\n", "activity", activityCount, activityInserted)
	fmt.Printf("%s\t%d\t%d\n", "monitoring", monitoringCount, monitoringInserted)
	fmt.Printf("%s\t%d\t%d\n", "record", recordCount, recordInserted)
	fmt.Printf("%s\t%d\t%d\n", "session", sessionCount, sessionInserted)
}
