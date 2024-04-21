package main

import (
	airlines "airline-checkin-system/sp_airlines"
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/lib/pq"
)

func main() {

	connStr := "host=localhost port=6432 user=user4 dbname=mydatabase4 password=password4 sslmode=disable"
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	airlines.InitializeDBRecords(db)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := lockAirlinesThenFlight(db); err != nil {
			log.Printf("Error in lock airlines then flight: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := lockFlightThenAirlines(db); err != nil {
			log.Printf("Error in lock flight then airlines: %v", err)
		}
	}()

	wg.Wait()
	//TODO print total time taken
}

func lockFlightThenAirlines(db *sql.DB) error {
	panic("unimplemented")
}

func lockAirlinesThenFlight(db *sql.DB) error {
	panic("unimplemented")
}
