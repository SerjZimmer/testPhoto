package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/SerjZimmer/testovoe1/pkg/fileserver"
)

func main() {
	conn, err := grpc.DialContext(
		context.Background(),
		"localhost:8081",
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	client := fileserver.NewFileServiceClient(conn)
	response, err := client.FileList(context.Background(), &fileserver.FileListRequest{})
	fmt.Println(response, err)
}
