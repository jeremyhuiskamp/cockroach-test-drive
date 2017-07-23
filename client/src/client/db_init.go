package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/mattes/migrate"
)

// dbInit is a helper that ensures that the database is ready to use.
// This includes:
// - waiting until the database is accessible
// - ensuring that the database itself exists
// - ensuring that the schema is up to date
//
// This is a bit sloppy and could probably be done more cleanly.
// Things to improve:
// - use Context to make interruptible
// - take into account administrative credentials
// - take into account default starting database if we can't connect to target one
// - look for existing packages that already implement this
// - keep the connection alive and return it for app use
type dbInit struct {
	dbURL string
}

func (d dbInit) Init() error {
	err := d.ensureDBExists()
	if err != nil {
		return fmt.Errorf("unable to ensure database exists: %s", err)
	}
	err = d.initSchema()
	if err != nil {
		return fmt.Errorf("unable to initialize schema: %s", err)
	}
	return nil
}

func (d dbInit) waitForConnection() (*sql.DB, error) {
	for i := 0; i < 10; i++ {
		// TODO: why do we have to mention postgres here?
		db, err := sql.Open("postgres", d.dbURL)
		if err != nil {
			log.Printf("connection attempt to %q unsuccessful: %s\n", d.dbURL, err)
			time.Sleep(1 * time.Second)
			continue
		}

		err = db.Ping()
		if err != nil {
			log.Printf("ping attempt to %q unsuccessful: %s\n", d.dbURL, err)
			_ = db.Close()
			time.Sleep(1 * time.Second)
			continue
		}

		return db, nil
	}

	return nil, fmt.Errorf("too many unsuccessful connection attempts to %q", d.dbURL)
}

func (d dbInit) ensureDBExists() error {
	parsedURL, err := url.Parse(d.dbURL)
	if err != nil {
		return fmt.Errorf("unable to parse url %q: %s", d.dbURL, err)
	}
	dbName := strings.TrimLeft(parsedURL.Path, "/")

	db, err := d.waitForConnection()
	if err != nil {
		return fmt.Errorf("unable to establish initial connection: %s", err)
	}

	defer func() {
		_ = db.Close()
	}()

	// lib/pq doesn't seem to like parameterization here.  Hopefully no one
	// is injecting evil database names?
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + dbName)
	if err != nil {
		return fmt.Errorf("unable to create database %q: %s", dbName, err)
	}

	log.Printf("created database %q\n", dbName)
	return nil
}

func (d dbInit) initSchema() error {
	// please build the scripts into the docker container:
	m, err := migrate.New("file:///migrate", d.dbURL)
	if err != nil {
		return fmt.Errorf("unable to creation migration: %s", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("unable to execute migration: %s", err)
	}

	return nil
}
