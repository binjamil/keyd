package service

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/binjamil/keyd/core"
	"github.com/binjamil/keyd/transact"
	"github.com/gorilla/mux"
)

var transactionLogger transact.TransactionLogger

func GetHandler(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	value, err := core.Get(key)
	if errors.Is(err, core.ErrorNoSuchKey) {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Write([]byte(value))
	log.Printf("GET key=%s\n", key)
}

func PutHandler(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	value, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	err = core.Put(key, string(value))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusCreated)
	transactionLogger.WritePut(key, string(value))
	log.Printf("PUT key=%s value=%s\n", key, string(value))
}

func DeleteHandler(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	err := core.Delete(key)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	transactionLogger.WriteDelete(key)
	log.Printf("DELETE key=%s\n", key)
}

func InitializeTransactionLog() error {
	var err error

	transactionLogger, err = transact.NewFileTransactionLogger("/tmp/transaction.log")
	if err != nil {
		return fmt.Errorf("failed to create event logger: %w", err)
	}

	events, errors := transactionLogger.ReadEvents()
	e, ok, count := transact.Event{}, true, 0

	for ok && err == nil {
		select {
		case err, ok = <-errors: // Retrieve any errors

		case e, ok = <-events:
			switch e.EventType {
			case transact.EventDelete: // Got a DELETE event!
				err = core.Delete(e.Key)
				count++
			case transact.EventPut: // Got a PUT event!
				err = core.Put(e.Key, e.Value)
				count++
			}
		}
	}

	log.Printf("%d events replayed\n", count)
	transactionLogger.Run()
	return err
}
