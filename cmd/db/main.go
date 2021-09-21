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

	var activitiesCount, activityLapsCount, activityRecordsCount, activitySessionsCount, monitoringsCount, recordsCount int
	var activitiesInserted, activityLapsInserted, activityRecordsInserted, activitySessionsInserted, monitoringsInserted, recordsInserted int

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
				recordsCount++
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
				recordsInserted++
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
			monitoringsCount++
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
			monitoringsInserted++
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
			return fmt.Errorf("unable to read file %s: %w\n", path, err)
		}

		data, err := fit.Decode(bytes.NewReader(b))
		if err != nil {
			return fmt.Errorf("unable to decode FIT data in file %s: %w\n", path, err)
		}

		activityFile, err := data.Activity()
		if err != nil {
			return fmt.Errorf("FIT data in file %s was not an activity: %w\n", path, err)
		}
		activitiesCount++
		id, err := queries.CreateActivity(ctx, pgsql.CreateActivityParams{
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
		activitiesInserted++

		for _, session := range activityFile.Sessions {
			activitySessionsCount++
			_, err := queries.CreateActivitySession(ctx, pgsql.CreateActivitySessionParams{
				Activity:         id,
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
			activitySessionsInserted++
		}
		for _, lap := range activityFile.Laps {
			activityLapsCount++
			_, err := queries.CreateActivityLap(ctx, pgsql.CreateActivityLapParams{
				Activity:         id,
				StartTs:          lap.StartTime,
				EndTs:            lap.Timestamp,
				Event:            int16(lap.Event),
				EventType:        int16(lap.EventType),
				Sport:            int16(lap.Sport),
				SubSport:         int16(lap.SubSport),
				TotalElapsedTime: sql.NullFloat64{Float64: lap.GetTotalElapsedTimeScaled(), Valid: true},
				TotalTimerTime:   sql.NullFloat64{Float64: lap.GetTotalTimerTimeScaled(), Valid: true},
				TotalDistance:    sql.NullFloat64{Float64: lap.GetTotalDistanceScaled(), Valid: true},
				TotalCalories:    int16(lap.TotalCalories),
				AvgSpeed:         sql.NullFloat64{Float64: lap.GetAvgSpeedScaled(), Valid: true},
				MaxSpeed:         sql.NullFloat64{Float64: lap.GetMaxSpeedScaled(), Valid: true},
				AvgHeartRate:     int16(lap.AvgHeartRate),
				MaxHeartRate:     int16(lap.MaxHeartRate),
			})
			if err != nil {
				if !errors.Is(err, sql.ErrNoRows) {
					return fmt.Errorf("failed to create lap: %w\n", err)
				}
			}
			activityLapsInserted++
		}
		for _, record := range activityFile.Records {
			activityRecordsCount++
			_, err := queries.CreateActivityRecord(ctx, pgsql.CreateActivityRecordParams{
				Activity:  id,
				Ts:        record.Timestamp,
				Altitude:  int16(record.Altitude),
				HeartRate: int16(record.HeartRate),
				Cadence:   int16(record.Cadence),
				Distance:  sql.NullFloat64{Float64: record.GetDistanceScaled(), Valid: true},
				Speed:     sql.NullFloat64{Float64: record.GetSpeedScaled(), Valid: true},
				Cycles:    int16(record.Cycles),
			})
			if err != nil {
				if !errors.Is(err, sql.ErrNoRows) {
					return fmt.Errorf("failed to create record: %w\n", err)
				}
			}
			activityRecordsInserted++
		}

		return nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Trouble walking directories: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("monitoringsCount\tcount\tinserted")
	fmt.Printf("%s\t%d\t%d\n", "activities", activitiesCount, activitiesInserted)
	fmt.Printf("%s\t%d\t%d\n", "activity_laps", activityLapsCount, activityLapsInserted)
	fmt.Printf("%s\t%d\t%d\n", "activity_records", activityRecordsCount, activityRecordsInserted)
	fmt.Printf("%s\t%d\t%d\n", "activity_sessions", activitySessionsCount, activitySessionsInserted)
	fmt.Printf("%s\t%d\t%d\n", "monitorings", monitoringsCount, monitoringsInserted)
	fmt.Printf("%s\t%d\t%d\n", "records", recordsCount, recordsInserted)
}
