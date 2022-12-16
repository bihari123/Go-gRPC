package service

import (
	"context"
	"log"
	pcbook "pcbook/proto"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LaptopServer struct {
	Store LaptopStore
}

func NewLaptopServer(store LaptopStore) *LaptopServer {
	return &LaptopServer{Store: store}
}

func (server *LaptopServer) CreateLaptop(
	ctx context.Context,
	req *pcbook.CreateLaptopRequest,
) (
	*pcbook.CreateLaptopResponse,
	error,
) {
	laptop := req.GetLaptop()
	log.Printf("received a create-laptop request with the id: %s", laptop.Id)

	if len(laptop.Id) > 0 {
		_, err := uuid.Parse(laptop.Id)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "LaptopId is not valid uuid: %v", err)
		}
	} else {
		id, err := uuid.NewRandom()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "cannot generate ID: %v", id)
		}
		laptop.Id = id.String()
	}

	// save the laptop to in-memory
	err := server.Store.Save(laptop)
	if err != nil {
		return nil, err
	}
	res := &pcbook.CreateLaptopResponse{
		Id: laptop.Id,
	}
	return res, nil
}
