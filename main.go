package main

import (
	"log"
	"net"
	"net/http"
	"os"

	pb "github.com/binjamil/keyd/grpc"
	"github.com/binjamil/keyd/service"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

func main() {
	grpcEnabled := os.Getenv("GRPC_ENABLED")

	err := service.InitializeTransactionLog()
	if err != nil {
		panic(err)
	}

	if grpcEnabled == "true" {
		// Create a gRPC server and register KeydServer to it
		s := grpc.NewServer()
		pb.RegisterKeydServer(s, &pb.ImplementedKeydServer{})

		// Open a TCP listening port
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("Failed to listen on TCP port 50051: %v", err)
		}

		// Start the gRPC server
		log.Println("gRPC server started on port 50051")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to start gRPC server: %v", err)
		}
	} else {
		// Create an HTTP router and add keyd handlers to it
		r := mux.NewRouter()
		r.HandleFunc("/v1/{key}", service.GetHandler).Methods(http.MethodGet)
		r.HandleFunc("/v1/{key}", service.PutHandler).Methods(http.MethodPut)
		r.HandleFunc("/v1/{key}", service.DeleteHandler).Methods(http.MethodDelete)

		// Start the HTTP server
		log.Println("HTTP server started on port 8000")
		if err := http.ListenAndServe(":8000", r); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}
}
