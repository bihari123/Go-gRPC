package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	pcbook "pcbook/proto"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LaptopServer struct {
	pcbook.UnimplementedLaptopServiceServer
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

	time.Sleep(6 * time.Second)

	if ctx.Err() == context.Canceled {
		log.Println("Request cancelled")
		return nil, status.Error(codes.Canceled, "request cancelled")
	}

	if ctx.Err() == context.DeadlineExceeded {
		log.Println("DeadLine exceeded")
		return nil, status.Error(codes.DeadlineExceeded, "deadline exceeded")
	}

	// save the laptop to in-memory
	err := server.Store.Save(laptop)
	if err != nil {
		code := codes.Internal
		if errors.Is(err, fmt.Errorf("Laptop with id: %s already exists", laptop.Id)) {
			code = codes.AlreadyExists
		}
		return nil, status.Errorf(code, "cannot save laptop to the store: %v", err)
	}
	log.Printf("\nlaptop with id %v saved", laptop.Id)
	res := &pcbook.CreateLaptopResponse{
		Id: laptop.Id,
	}
	return res, nil
}

func (server *LaptopServer) mustEmbedUnimplementedLaptopServiceServer() {}
