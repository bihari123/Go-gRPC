package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	pcbook "pcbook/proto"
	"pcbook/service"

	"google.golang.org/grpc"
)

func main() {
	port := flag.Int("port", 0, "the server port")

	flag.Parse()

	log.Printf("start server on port: %d", *port)

	laptopStore := service.NewInMemoryLaptopStore()
	imageStore := service.NewDiskImageStore("img")
	rateStore := service.NewInMemoryratingStore()
	laptopServer := service.NewLaptopServer(laptopStore, imageStore, rateStore)

	grpcServer := grpc.NewServer()

	pcbook.RegisterLaptopServiceServer(grpcServer, laptopServer)

	address := fmt.Sprintf("0.0.0.0:%d", *port)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("could not start server", err)
	}

	err = grpcServer.Serve(listener)

	if err != nil {
		log.Fatal("could not start server", err)
	}
}
