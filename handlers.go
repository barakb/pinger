package pinger

import (
	"net/http"
	"github.com/gorilla/mux"
	"time"
	"encoding/json"
)


func Pings(w http.ResponseWriter, r *http.Request) {
	pingsResult := ParallelPing(config.Addresses, 3 * time.Second)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(pingsResult); err != nil {
		panic(err)
	}
}

func CFGGetMachines(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(config.Addresses); err != nil {
		panic(err)
	}
}

func CFGAddMachine(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	config.AddMachine(name)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(config.Addresses); err != nil {
		panic(err)
	}
}

func CFGGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(config); err != nil {
		panic(err)
	}
}


