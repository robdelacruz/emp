package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	os.Args = os.Args[1:]
	switches, parms := parseArgs(os.Args)

	db, err := openDB(switches)
	if err != nil {
		log.Fatal(err)
	}

	cmd := "serve"
	if len(parms) > 0 {
		if parms[0] == "serve" || parms[0] == "info" || parms[0] == "help" {
			cmd = parms[0]
			parms = parms[1:]
		}
	}

	switch cmd {
	case "serve":
		port := "8000"
		if len(parms) > 0 {
			port = parms[0]
		}

		http.HandleFunc("/", rootHandler(db))
		http.HandleFunc("/emp/", empHandler(db))

		fmt.Printf("Listening on %s...\n", port)
		err = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
		log.Fatal(err)
	}
}

func parseArgs(args []string) (map[string]string, []string) {
	switches := map[string]string{}
	parms := []string{}

	standaloneSwitches := []string{}
	definitionSwitches := []string{"F"}
	fNoMoreSwitches := false
	curKey := ""

	for _, arg := range args {
		if fNoMoreSwitches {
			// any arg after "--" is a standalone parameter
			parms = append(parms, arg)
		} else if arg == "--" {
			// "--" means no more switches to come
			fNoMoreSwitches = true
		} else if strings.HasPrefix(arg, "--") {
			switches[arg[2:]] = "y"
			curKey = ""
		} else if strings.HasPrefix(arg, "-") {
			if inList(definitionSwitches, arg[1:]) {
				// -a "val"
				curKey = arg[1:]
				continue
			}
			for _, ch := range arg[1:] {
				// -a, -b, -ab
				sch := string(ch)
				if inList(standaloneSwitches, sch) {
					switches[sch] = "y"
				}
			}
		} else if curKey != "" {
			switches[curKey] = arg
			curKey = ""
		} else {
			// standalone parameter
			parms = append(parms, arg)
		}
	}

	return switches, parms
}

func inList(ss []string, v string) bool {
	for _, s := range ss {
		if v == s {
			return true
		}
	}
	return false
}

func rootHandler(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Allow requests from all sites.
		w.Header().Set("Access-Control-Allow-Origin", "*")

		w.WriteHeader(http.StatusBadRequest)
	}
}

func valuesToMap(vals url.Values) map[string]interface{} {
	m := map[string]interface{}{}
	for k, vv := range vals {
		v := vv[0] // discard any extra field definitions

		// int fields
		if k == "userid" {
			m[k], _ = strconv.Atoi(v)
			continue
		}
		// string fields (default)
		m[k] = v
	}
	return m
}

func empHandler(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		switch r.Method {
		case "PUT":
			// If querystring fields are specified, update only specific fields.
			// Else, update the whole emp record, reading emp from request body.

			sre := `^/emp/(\d+)/?`
			re := regexp.MustCompile(sre)
			matches := re.FindStringSubmatch(r.URL.Path)
			if matches == nil {
				w.WriteHeader(http.StatusBadRequest)
			}
			empId, _ := strconv.ParseInt(matches[1], 10, 64)

			// Update specific fields only if passed in querystring. Ex:
			// /emp/123?firstname=Oscar&lastname=Grouch
			// request body is ignored
			values, _ := url.ParseQuery(r.URL.RawQuery)
			fields := valuesToMap(values)
			if len(fields) > 0 {
				err := updateEmpFields(db, empId, fields)
				if err != nil {
					log.Printf("empHandler() db update error (%s)\n", err)
					http.Error(w, "Failed to update Emp.", http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusOK)
				return
			}

			// request body contains emp in json format.
			// /emp/123 {...}

			var emp Emp
			err := json.NewDecoder(r.Body).Decode(&emp)
			if err != nil {
				log.Printf("empHandler() json decoding error:\n%s\n", err)
				http.Error(w, "Invalid request body.", http.StatusBadRequest)
				return
			}
			err = updateEmp(db, empId, &emp)
			if err != nil {
				log.Printf("empHandler() db update error (%s)\n", err)
				http.Error(w, "Failed to update Emp.", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
		}
	}
}
