package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/SerjZimmer/testovoe1/internal/pkg/file"
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
	cli := fileserver.NewFileServiceClient(conn)

	response, err := cli.Download(
		context.Background(),
		&fileserver.DownloadRequest{Id: "a2b5c495-172a-4e99-a64d-57e5a49c26f3"},
	)
	if err != nil {
		log.Fatalln("Failed to load file:", err)
	}
	md, err := response.Header()
	if err != nil {
		log.Fatalln("Failed get headers:", err)
	}
	r, w := io.Pipe()
	f := file.NewFromMetadata(md, r)
	go copyFromResponse(w, response)

	imagePath := fmt.Sprintf("%s.%s", "testdownloaded", "jpeg")

	file, err := os.Create(filepath.Clean(imagePath))
	if err != nil {
		log.Fatalln("Failed create file:", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalln("Failed ")
		}
	}(file)

	_, err = file.ReadFrom(f)
	if err != nil {
		log.Fatalln("Failed read to file:", err)
	}
}

func copyFromResponse(w *io.PipeWriter, res fileserver.FileService_DownloadClient) {
	message := new(fileserver.DownloadResponse)
	var err error
	for {
		err = res.RecvMsg(message)
		if err == io.EOF {
			_ = w.Close()
			break
		}
		if err != nil {
			_ = w.CloseWithError(err)
			break
		}
		if len(message.GetChunk()) > 0 {
			_, err = w.Write(message.Chunk)
			if err != nil {
				_ = res.CloseSend()
				break
			}
		}
		message.Chunk = message.Chunk[:0]
	}
}
