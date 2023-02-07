package fileserver

import (
	"errors"
	"io"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/SerjZimmer/testovoe1/internal/pkg/storage"
	"github.com/SerjZimmer/testovoe1/pkg/fileserver"
)

const chunkSize = 1024 * 3

func (s *Server) Download(req *fileserver.DownloadRequest, server fileserver.FileService_DownloadServer) error {
	if req.GetId() == "" {
		return status.Error(codes.InvalidArgument, "id is required")
	}

	f, err := s.storage.Get(req.GetId())
	var e *storage.NotFoundErr
	if errors.As(err, &e) {
		return status.Errorf(codes.NotFound, "cannot get file, err: %s", err.Error())
	}

	if err != nil {
		return status.Errorf(codes.Internal, "cannot get file, err: %s", err.Error())
	}

	err = server.SendHeader(f.Metadata())
	if err != nil {
		return status.Error(codes.Internal, "error during sending header")
	}

	chunk := &fileserver.DownloadResponse{Chunk: make([]byte, chunkSize)}
	var n int

	for {
		n, err = f.Read(chunk.Chunk)
		if err == io.EOF {
			break
		}
		if err != nil {
			return status.Errorf(codes.Internal, "io.ReadAll: %v", err)
		}

		chunk.Chunk = chunk.Chunk[:n]
		serverErr := server.Send(chunk)
		if serverErr != nil {
			return status.Errorf(codes.Internal, "server.Send: %v", serverErr)
		}
	}

	return nil

}
