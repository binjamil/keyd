package service

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/binjamil/keyd/core"
	"github.com/gorilla/mux"
)

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

	log.Printf("DELETE key=%s\n", key)
}
