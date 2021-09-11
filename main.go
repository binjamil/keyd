package main

import (
	"log"
	"net/http"
	"os"

	"github.com/binjamil/keyd/service"
	"github.com/gorilla/mux"
)

func main() {
	tlsCert := os.Getenv("TLS_CERTIFICATE")
	tlsKey := os.Getenv("TLS_KEY")

	err := service.InitializeTransactionLog()
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()

	r.HandleFunc("/v1/{key}", service.GetHandler).Methods(http.MethodGet)
	r.HandleFunc("/v1/{key}", service.PutHandler).Methods(http.MethodPut)
	r.HandleFunc("/v1/{key}", service.DeleteHandler).Methods(http.MethodDelete)

	if tlsCert != "" && tlsKey != "" {
		log.Println("TLS configuration found. Server started over https...")
		log.Fatal(http.ListenAndServeTLS(":8080", tlsCert, tlsKey, r))
	} else {
		log.Println("TLS configuration not found. Server started over http...")
		log.Fatal(http.ListenAndServe(":8080", r))
	}
}
