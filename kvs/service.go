package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var transact TransactionLogger

func keyValuePutHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	value, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w,
			err.Error(),
			http.StatusInternalServerError)
		return
	}

	err = Put(key, string(value))

	if err != nil {
		http.Error(w,
			err.Error(),
			http.StatusInternalServerError)
		return
	}

	transact.WritePut(key, string(value))

	w.WriteHeader(http.StatusCreated)
}

func keyValueGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	value, err := Get(key)
	if errors.Is(err, ErrorNoSuchKey) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(value))
}

func keyValueDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	err := Delete(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	transact.WriteDelete(key)

	w.WriteHeader(http.StatusOK)
}

func initializeTransactionLog(source SourceLogType) error {
	var err error

	transact, err := transactionLogFactory(source)
	if err != nil {
		return fmt.Errorf("failed to create event logger: %w", err)
	}

	events, errors := transact.ReadEvents()
	e, ok := Event{}, true

	for ok && err == nil {
		select {
		case err, ok = <-errors:
		case e, ok = <-events:
			switch e.EventType {
			case EventDelete:
				err = Delete(e.Key)
			case EventPut:
				err = Put(e.Key, e.Value)
			}
		}
	}

	transact.Run()

	return err
}

type SourceLogType byte

const (
	_                  = iota
	File SourceLogType = iota
	Postgres
)

func transactionLogFactory(source SourceLogType) (TransactionLogger, error) {
	switch source {
	case File:
		return NewFileTransactionLogger("transaction.log")
	case Postgres:
		return NewPostgresTransactionLogger(PostgresDBParams{
			dbName:   "kvs",
			host:     "localhost",
			user:     "admin",
			password: "password",
		})
	default:
		return nil, fmt.Errorf("unknown source log type")
	}
}

func main() {
	err := initializeTransactionLog(Postgres)
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()

	r.HandleFunc("/v1/{key}", keyValuePutHandler).Methods("PUT")
	r.HandleFunc("/v1/{key}", keyValueGetHandler).Methods("GET")
	r.HandleFunc("/v1/{key}", keyValueDeleteHandler).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", r))
}
