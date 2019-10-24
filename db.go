package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func openDB(switches map[string]string) (*sql.DB, error) {
	// spot db file is read from the following, in order of priority:
	// 1. -F <file>
	// 2. EMPDB env var
	// 3. /usr/local/share/emp/emp.db (default)
	dbfile := os.Getenv("EMPDB")
	if switches["F"] != "" {
		dbfile = switches["F"]
	}
	if dbfile == "" {
		dirpath := filepath.Join(string(os.PathSeparator), "usr", "local", "share", "emp")
		os.MkdirAll(dirpath, os.ModePerm)
		dbfile = filepath.Join(dirpath, "emp.db")
	}

	fmt.Printf("dbfile: '%s'\n", dbfile)
	db, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		return nil, fmt.Errorf("openDB(): Error opening '%s' (%s)\n", dbfile, err)
	}

	ensureCreateTables(db)
	return db, nil
}

func ensureCreateTables(db *sql.DB) {
}
