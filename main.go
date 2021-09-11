package main

import (
	"log"
	"net/http"

	"github.com/binjamil/keyd/service"
	"github.com/gorilla/mux"
)

func main() {
	err := service.InitializeTransactionLog()
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()

	r.HandleFunc("/v1/{key}", service.GetHandler).Methods(http.MethodGet)
	r.HandleFunc("/v1/{key}", service.PutHandler).Methods(http.MethodPut)
	r.HandleFunc("/v1/{key}", service.DeleteHandler).Methods(http.MethodDelete)

	log.Fatal(http.ListenAndServe(":8080", r))
}
