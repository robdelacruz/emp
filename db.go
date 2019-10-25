package main

import (
	"database/sql"
	"fmt"
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

type Emp struct {
	EmpId     int64  `json:"empid"`
	UserId    int64  `json:"userid"`
	Empno     string `json:"empno"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Title     string `json:"title"`
}

func (emp Emp) String() string {
	return fmt.Sprintf(`
EmpId:     %d
UserId:    %d
Empno:     %s
Firstname: %s
Lastname:  %s
Title:     %s
`, emp.EmpId, emp.UserId, emp.Empno, emp.Firstname, emp.Lastname, emp.Title)
}

func updateEmp(db *sql.DB, empId int64, emp *Emp) error {
	sqlstr := "UPDATE emp set userid = ?, empno = ?, firstname = ?, lastname = ?, title = ? WHERE empid = ?"
	fmt.Printf("updateEmp():\n%s\n", sqlstr)
	stmt, _ := db.Prepare(sqlstr)
	_, err := stmt.Exec(emp.UserId, emp.Empno, emp.Firstname, emp.Lastname, emp.Title, empId)
	if err != nil {
		return fmt.Errorf("updateEmp() -\n%v\n(%s)\n", emp, err)
	}
	return nil
}

func updateEmpFields(db *sql.DB, empId int64, fields map[string]interface{}) error {
	validFields := []string{"empno", "userid", "firstname", "lastname", "title"}

	ss := []string{}
	vv := [](interface{}){}
	for k, v := range fields {
		if inList(validFields, k) {
			ss = append(ss, fmt.Sprintf("%s = ?", k))
			vv = append(vv, v)
		}
	}
	if len(ss) == 0 {
		return nil
	}
	vv = append(vv, empId)

	setClause := strings.Join(ss, ", ")
	sqlstr := fmt.Sprintf("UPDATE emp SET %s WHERE empid = ?", setClause)
	fmt.Printf("updateEmpFields():\n%s\n", sqlstr)

	_, err := db.Exec(sqlstr, vv...)
	if err != nil {
		return err
	}
	return nil
}
