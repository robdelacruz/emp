package main

import (
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

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
	sqlstr := `PRAGMA foreign_keys = ON;
BEGIN TRANSACTION;
CREATE TABLE IF NOT EXISTS emp (empid INTEGER PRIMARY KEY NOT NULL, userid INTEGER, empno TEXT, firstname TEXT, lastname TEXT, title TEXT);
INSERT OR REPLACE INTO emp (empid, empno, firstname, lastname, title) VALUES (1, '123', 'Oscar', 'the Grouch', 'Actor');
END TRANSACTION;`
	_, err := db.Exec(sqlstr)
	if err != nil {
		panic(err)
	}
}

func updateEmp(db *sql.DB, empid int64, vals url.Values) {
	empFields := []string{"empno", "firstname", "lastname", "title"}

	ss := []string{}
	vv := [](interface{}){}
	for k, v := range vals {
		if listContains(empFields, k) {
			ss = append(ss, fmt.Sprintf("%s = ?", k))
			vv = append(vv, v[0])
		}
	}
	if len(ss) == 0 {
		return
	}

	setClause := strings.Join(ss, ", ")
	sqlstr := fmt.Sprintf("UPDATE emp SET %s WHERE empid = %d", setClause, empid)
	fmt.Printf("updateEmp():\n%s\n", sqlstr)

	_, err := db.Exec(sqlstr, vv...)
	if err != nil {
		panic(err)
	}

}
