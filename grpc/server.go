package grpc

import (
	"context"
	"log"

	"github.com/binjamil/keyd/core"
	"github.com/binjamil/keyd/service"
)

type ImplementedKeydServer struct {
	UnimplementedKeydServer
}

func (s *ImplementedKeydServer) Get(ctx context.Context, r *GetRequest) (*GetResponse, error) {
	value, err := core.Get(r.Key)

	log.Printf("GET key=%s\n", r.Key)
	return &GetResponse{Value: value}, err
}

func (s *ImplementedKeydServer) Put(ctx context.Context, r *PutRequest) (*PutResponse, error) {
	err := core.Put(r.Key, r.Value)
	service.TransactionLogger.WritePut(r.Key, r.Value)

	log.Printf("PUT key=%s value=%s\n", r.Key, r.Value)
	return &PutResponse{}, err
}

func (s *ImplementedKeydServer) Delete(ctx context.Context, r *DeleteRequest) (*DeleteResponse, error) {
	err := core.Delete(r.Key)
	service.TransactionLogger.WriteDelete(r.Key)

	log.Printf("DELETE key=%s\n", r.Key)
	return &DeleteResponse{}, err
}
