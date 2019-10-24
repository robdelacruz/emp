package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
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
			if listContains(definitionSwitches, arg[1:]) {
				// -a "val"
				curKey = arg[1:]
				continue
			}
			for _, ch := range arg[1:] {
				// -a, -b, -ab
				sch := string(ch)
				if listContains(standaloneSwitches, sch) {
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

func listContains(ss []string, v string) bool {
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

func empHandler(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		switch r.Method {
		case "PUT":
			r.ParseForm()
			for k, v := range r.Form {
				fmt.Fprintf(w, "%s: %s\n", k, v)
			}

			w.WriteHeader(http.StatusOK)
		}
	}
}
