package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	server "github.com/SerjZimmer/testovoe1/internal/app/fileserver"
	"github.com/SerjZimmer/testovoe1/internal/pkg/storage"
	"github.com/SerjZimmer/testovoe1/pkg/fileserver"
)

func main() {
	// Create a listener on TCP port
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	fileStorage := storage.New("img")
	srv := server.New(fileStorage)
	// Create a gRPC server object
	s := grpc.NewServer()
	// Attach the Greeter service to the server
	fileserver.RegisterFileServiceServer(s, srv)
	// Serve gRPC server
	log.Println("Serving gRPC on localhost:8081")
	log.Fatalln(s.Serve(lis))

}
