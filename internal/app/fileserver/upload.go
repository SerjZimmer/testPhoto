package fileserver

import (
	"bytes"
	"io"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/SerjZimmer/testovoe1/internal/pkg/file"
	"github.com/SerjZimmer/testovoe1/pkg/fileserver"
)

const maxImageSize = 1 << 20

func (s *Server) Upload(stream fileserver.FileService_UploadServer) error {
	imageData := bytes.Buffer{}
	imageSize := 0
	for {
		req, err := stream.Recv() //получили кусок данных
		if err == io.EOF {
			break
		}
		if err != nil {
			return logError(status.Errorf(codes.Unknown, "cannot receive chunk data: %v", err))
		}
		chunk := req.GetChunk()
		size := len(chunk)
		imageSize += size
		if imageSize > maxImageSize {
			return logError(status.Errorf(codes.InvalidArgument, "image is too large: %d > %d", imageSize, maxImageSize))
		}

		_, err = imageData.Write(chunk)
		if err != nil {
			return logError(status.Errorf(codes.Internal, "cannot write chunk data: %v", err))
		}
	}

	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return logError(status.Error(codes.Internal, "cannot get headers"))
	}

	fileWrap := file.NewFromMetadata(md, &imageData)

	id, err := s.storage.Save(fileWrap)
	if err != nil {
		return logError(status.Errorf(codes.Internal, "cannot save file: %v", err))
	}

	err = stream.SendAndClose(&fileserver.UploadResponse{
		Id: id,
	})
	if err != nil {
		return logError(status.Errorf(codes.Internal, "cannot send response: %v", err))
	}
	log.Printf("saved image with id: %s, size: %d", id, imageSize)
	return nil
}

func logError(err error) error {
	if err != nil {
		log.Print(err)
	}
	return err
}
