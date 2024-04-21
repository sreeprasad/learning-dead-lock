package main

import (
	airlines "airline-checkin-system/sp_airlines"
	"database/sql"
	"log"
	"sync"
	"time"

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

	startTime := time.Now()

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
	duration := time.Since(startTime)
	log.Printf("Total time taken: %s\n", duration)

}

func lockFlightThenAirlines(db *sql.DB) error {

	transaction, err := db.Begin()
	if err != nil {
		return err
	}
	defer transaction.Rollback()

	if _, err := transaction.Exec("SELECT * FROM public.flights WHERE id = 1 FOR UPDATE"); err != nil {
		return err
	}

	time.Sleep(1 * time.Second)

	if _, err := transaction.Exec("SELECT * from public.airlines WHERE id = 1 FOR UPDATE"); err != nil {
		return err
	}

	return transaction.Commit()

}

func lockAirlinesThenFlight(db *sql.DB) error {

	transaction, err := db.Begin()
	if err != nil {
		return err
	}
	defer transaction.Rollback()

	if _, err := transaction.Exec("SELECT * from public.airlines WHERE id = 1 FOR UPDATE"); err != nil {
		return err
	}

	time.Sleep(1 * time.Second)

	if _, err := transaction.Exec("SELECT * FROM public.flights WHERE id = 1 FOR UPDATE"); err != nil {
		return err
	}

	return transaction.Commit()

}
