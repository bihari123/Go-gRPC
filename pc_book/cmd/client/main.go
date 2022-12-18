package main

import (
	"context"
	"flag"
	"io"
	"log"
	pcbook "pcbook/proto"
	"pcbook/sample"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func createLaptop(laptopClient pcbook.LaptopServiceClient) {
	laptop := sample.NewLaptop()

	req := &pcbook.CreateLaptopRequest{
		Laptop: laptop,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := laptopClient.CreateLaptop(ctx, req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.AlreadyExists {
			// not a big deal
			log.Println("Laptop already exists")
		} else {
			log.Fatal("cannot create log ", err)
		}
		return
	}

	log.Printf("created laptop with id: %s", res.Id)
}

func searchLaptop(laptopClient pcbook.LaptopServiceClient, filter *pcbook.Filter) {
	log.Println("search filter: ", filter)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pcbook.SearchLaptopRequest{Filter: filter}
	stream, err := laptopClient.SearchLaptop(ctx, req)
	if err != nil {
		log.Fatal("cannot search laptop: ", err)
	}

	for {
		res, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return
			} else {
				log.Fatal("cannot recieve response ", err)
			}
		}
		laptop := res.GetLaptop()

		log.Print("- found: ", laptop.GetId())
		log.Print(" + brand: ", laptop.GetBrand())
		log.Print(" +name: ", laptop.GetName())
		log.Print(" +ram: ", laptop.GetRam().GetValue(), laptop.GetRam().GetUnit())
		log.Print(" +price: ", laptop.GetPriceUsd())
	}
}

func main() {
	serverAddress := flag.String("address", "", "the server address")
	flag.Parse()
	log.Printf("dial server: %s\n", *serverAddress)

	conn, err := grpc.Dial(*serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatal("cannot dial server ", err)
	}

	laptopClient := pcbook.NewLaptopServiceClient(conn)

	for i := 0; i < 10; i++ {
		createLaptop(laptopClient)
	}

	filter := &pcbook.Filter{
		MaxPriceUsd: 3000,
		MinCpuCores: 4,
		MinCpuGhz:   2.5,
		MinRam:      &pcbook.Memory{Value: 8, Unit: pcbook.Memory_GIGABYTE},
	}
	searchLaptop(laptopClient, filter)
}
