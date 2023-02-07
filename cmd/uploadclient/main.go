package main

import (
	"bufio"
	"context"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"github.com/SerjZimmer/testovoe1/internal/pkg/file"
	"github.com/SerjZimmer/testovoe1/pkg/fileserver"
)

const ctxTime = 5
const defaultLen = 2048

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
	cli := fileserver.NewFileServiceClient(conn)

	uploadImage(cli, "test", "./test.jpeg")
}

func uploadImage(client fileserver.FileServiceClient, name string, imagePath string) {
	fileFromDisc, err := os.Open(filepath.Clean(imagePath))
	if err != nil {
		log.Fatal("cannot open image fileFromDisc: ", err)
	}
	defer func(fileFromDisc *os.File) {
		err := fileFromDisc.Close()
		if err != nil {
			log.Fatalln("Fail fileFromDisk.Close")
		}
	}(fileFromDisc)

	ctx, cancel := context.WithTimeout(context.Background(), ctxTime*time.Second)
	defer cancel()

	stat, err := fileFromDisc.Stat()
	if err != nil {
		log.Fatal("cannot get stat: ", err)
	}

	fileToSend := file.New(bufio.NewReader(fileFromDisc), name, "jpeg", int(stat.Size()))

	ctx = metadata.NewOutgoingContext(ctx, fileToSend.Metadata())

	stream, err := client.Upload(ctx)
	if err != nil {
		log.Fatal("cannot upload image: ", err)
	}

	buffer := make([]byte, defaultLen)

	for {
		n, err := fileToSend.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("cannot read chunk to buffer: ", err)
		}

		req := &fileserver.UploadRequest{
			Chunk: buffer[:n],
		}

		err = stream.Send(req)
		if err != nil {
			log.Fatal("cannot send chunk to server: ", err)
		}

	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal("cannot receive response: ", err)
	}

	log.Printf("image uploaded with id: %s", res.GetId())
}
