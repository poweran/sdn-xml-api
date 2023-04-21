package main

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"sdn-xml-api/internal/database/repository"
	"strings"
	"sync"
)

type state struct {
	updating bool
	names    []repository.Person
}

var (
	mutex    sync.RWMutex
	appState state
)

func initRouter(db *sql.DB) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		if _, updating := getState(); updating {
			w.WriteHeader(http.StatusServiceUnavailable)
			writeStringToRW(w, `{"result": false, "info": "updating"}`)
			return
		}

		if err := repository.Update(db); err != nil {
			log.Println("Error updating:", err)
			setState(state{updating: false})
			w.WriteHeader(http.StatusServiceUnavailable)
			writeStringToRW(w, `{"result": false, "info": "service unavailable"}`)
			return
		}

		setState(state{updating: false})

		w.WriteHeader(http.StatusOK)
		writeStringToRW(w, `{"result": true, "info": "", "code": 200}`)
	})

	r.HandleFunc("/state", func(w http.ResponseWriter, r *http.Request) {
		s, updating := getState()

		if updating {
			w.WriteHeader(http.StatusServiceUnavailable)
			writeStringToRW(w, `{"result": false, "info": "updating"}`)
			return
		}

		if len(s.names) == 0 {
			w.WriteHeader(http.StatusOK)
			writeStringToRW(w, `{"result": false, "info": "empty"}`)
			return
		}

		w.WriteHeader(http.StatusOK)
		writeStringToRW(w, `{"result": true, "info": "ok"}`)
	})

	r.HandleFunc("/get_names", func(w http.ResponseWriter, r *http.Request) {
		s, updating := getState()

		if updating {
			w.WriteHeader(http.StatusServiceUnavailable)
			writeStringToRW(w, `{"result": false, "info": "updating"}`)
			return
		}

		name := r.FormValue("name")
		matchType := strings.ToLower(r.FormValue("type"))

		if name == "" {
			w.WriteHeader(http.StatusBadRequest)
			writeStringToRW(w, `{"result": false, "info": "missing name parameter"}`)
			return
		}

		var result []repository.Person

		for _, p := range s.names {
			if matchType == "" {
				result = append(result, p)
			} else if matchType == "strong" && p.FirstName+" "+p.LastName == name {
				result = append(result, p)
			} else if matchType == "weak" && (strings.Contains(p.FirstName, name) || strings.Contains(p.LastName, name)) {
				result = append(result, p)
			}
		}

		if len(result) == 0 {
			w.WriteHeader(http.StatusOK)
			writeStringToRW(w, `{"result": false, "info": "no matches"}`)
			return
		}

		b, err := json.Marshal(result)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			writeStringToRW(w, `{"result": false, "info": "error marshalling result"}`)
			return
		}

		w.WriteHeader(http.StatusOK)
		writeBytesToRW(w, b)
	})

	return r
}

func setState(s state) {
	mutex.Lock()
	appState = s
	mutex.Unlock()
}

func getState() (state, bool) {
	mutex.RLock()
	defer mutex.RUnlock()
	return appState, appState.updating
}

func writeBytesToRW(w http.ResponseWriter, b []byte) {
	if _, err := w.Write(b); err != nil {
		log.Println(err.Error())
	}
}

func writeStringToRW(w http.ResponseWriter, s string) {
	writeBytesToRW(w, []byte(s))
}
