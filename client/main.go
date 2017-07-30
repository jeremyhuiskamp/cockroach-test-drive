package main

import (
	"flag"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	_ "github.com/mattes/migrate/database/postgres"
	_ "github.com/mattes/migrate/source/file"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	dbURL := flag.String(
		"db-url",
		"postgres://db1:26257/bank_test?sslmode=disable",
		"Connection URL for the database.")
	addr := flag.String(
		"listen-address",
		":8079",
		"The address to listen on for HTTP requests.")

	go func() {
		initer := dbInit{*dbURL}
		err := initer.Init()
		if err != nil {
			log.Printf("WARN: %s\n", err)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}
