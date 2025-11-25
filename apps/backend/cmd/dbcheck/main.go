package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

func main() {
	f, err := os.Create("db_status.txt")
	if err != nil {
		os.Exit(1)
	}
	defer f.Close()

	fmt.Fprintln(f, "Starting DB check...")
	dsn := "postgres://postgres@localhost:5432/ark?sslmode=disable"
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		fmt.Fprintf(f, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())
	fmt.Fprintln(f, "Connected to DB.")

	var tableName string
	err = conn.QueryRow(context.Background(),
		"SELECT table_name FROM information_schema.tables WHERE table_name = 'asset_logs'").Scan(&tableName)

	if err != nil {
		if err == pgx.ErrNoRows {
			fmt.Fprintln(f, "Table 'asset_logs' DOES NOT exist.")
		} else {
			fmt.Fprintf(f, "Query failed: %v\n", err)
		}
	} else {
		fmt.Fprintln(f, "Table 'asset_logs' EXISTS.")
	}

	// Check schema_version
	var version int
	err = conn.QueryRow(context.Background(), "SELECT version FROM schema_version").Scan(&version)
	if err != nil {
		fmt.Fprintf(f, "Could not read schema_version: %v\n", err)
	} else {
		fmt.Fprintf(f, "Current schema version: %d\n", version)
	}
}
