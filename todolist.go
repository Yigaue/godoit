package main

import (
	"io"
	"net/http"

	"github.com/gorilla/mux" // for creating router.
	log "github.com/sirupsen/logrus"
)

func ApiHealth(w http.ResponseWriter, r *http.Request) {
	log.Info("API health is ok")
	w.Header().Set("content-type", "application/json")
	io.WriteString(w, `{alive: true}`)
}

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetReportCaller(true)
}

func main() {
	log.Info("Starting TodoList API server")
	router := mux.NewRouter()
	router.HandleFunc("/health", ApiHealth).Methods("Get")
	http.ListenAndServe(":8000", router)
}