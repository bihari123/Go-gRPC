package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	pcbook "pcbook/proto"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const maxImageSize = 1 << 20

type LaptopServer struct {
	LaptopStore LaptopStore
	ImageStore  ImageStore
	pcbook.UnimplementedLaptopServiceServer
}

func NewLaptopServer(laptopstore LaptopStore, imageStore ImageStore) *LaptopServer {
	return &LaptopServer{laptopstore, imageStore, pcbook.UnimplementedLaptopServiceServer{}}
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

	// time.Sleep(6 * time.Second)

	if ctx.Err() == context.Canceled {
		log.Println("Request cancelled")
		return nil, status.Error(codes.Canceled, "request cancelled")
	}

	if ctx.Err() == context.DeadlineExceeded {
		log.Println("DeadLine exceeded")
		return nil, status.Error(codes.DeadlineExceeded, "deadline exceeded")
	}

	// save the laptop to in-memory
	err := server.LaptopStore.Save(laptop)
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

func (server *LaptopServer) SearchLaptop(
	req *pcbook.SearchLaptopRequest,
	stream pcbook.LaptopService_SearchLaptopServer,
) error {
	filter := req.GetFilter()
	log.Printf("recieved a search-laptop request with filter: %v", filter)

	err := server.LaptopStore.Search(
		filter,
		func(laptop *pcbook.Laptop) error {
			res := &pcbook.SearchLaptopResponse{
				Laptop: laptop,
			}
			err := stream.Send(res)
			if err != nil {
				return err
			}

			log.Printf("sent laptop with id: %s", laptop.GetId())
			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("could not search the laptop: %w", err)
	}

	return nil
}

func (server *LaptopServer) UploadImage(stream pcbook.LaptopService_UploadImageServer) error {
	req, err := stream.Recv()
	if err != nil {
		log.Println("cannot recieve image file info ", err)
		return status.Error(codes.Unknown, "cannot recieve image info")
	}

	laptopID := req.GetInfo().GetLaptopId()
	imageType := req.GetInfo().GetImageType()

	log.Println("got the request with imageType: ", imageType, " and laptopId: ", laptopID)

	laptop, err := server.LaptopStore.Find(laptopID)
	if err != nil {
		return status.Error(
			codes.Internal,
			fmt.Sprintf("cannot find the laptop with id: %v", laptopID),
		)
	}

	if laptop == nil {
		return status.Errorf(codes.InvalidArgument, "laptop %s doesn't exist: ", laptopID)
	}

	imageData := bytes.Buffer{}

	imageSize := 0

	for {
		log.Println("waiting to recieve more data")
		req, err := stream.Recv()
		if err == io.EOF {
			log.Println("no more data to recieve")
			break
		}

		if err != nil {
			return status.Errorf(codes.Unknown, "cannot recieve chunk data: %v", err)
		}

		chunk := req.GetChunkData()
		imageSize += len(chunk)

		_, err = imageData.Write(chunk)

		if err != nil {
			return status.Errorf(codes.Internal, "cannot write chink data: %v", err)
		}

	}

	imageId, err := server.ImageStore.Save(laptopID, imageType, imageData)
	if err != nil {
		return status.Errorf(codes.Internal, "cannot save image to the store: %v", err)
	}

	res := &pcbook.UploadImageResponse{Id: imageId, Size: uint32(imageSize)}

	err = stream.SendAndClose(res)

	if err != nil {
		return err
	}

	return nil
}
