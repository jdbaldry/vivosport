package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/jdbaldry/vivosport/pgsql"
	_ "github.com/lib/pq"
)

func main() {
	ctx := context.Background()

	db, err := sql.Open("postgres", "dbname=vivosport password=vivosport sslmode=disable user=vivosport")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to open DB connection: %v", err)
		os.Exit(1)
	}

	queries := pgsql.New(db)

	activity, err := queries.CreateActivity(ctx, pgsql.CreateActivityParams{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create activity: %v", err)
		os.Exit(1)
	}

	fmt.Printf("Created activity: %#v\n", activity)
}
