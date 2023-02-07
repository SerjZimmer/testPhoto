package fileserver

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/SerjZimmer/testovoe1/pkg/fileserver"
)

func (s *Server) FileList(ctx context.Context, server *fileserver.FileListRequest) (*fileserver.FileListResponse, error) {
	fileRows, err := s.storage.List()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot list files, err: %s", err.Error())
	}

	result := make([]*fileserver.FileInfo, 0, len(fileRows))
	for _, v := range fileRows {
		result = append(result, &fileserver.FileInfo{
			Id:        v.ID,
			Name:      v.Name,
			CreatedAt: timestamppb.New(v.CreatedAt),
			UpdatedAt: timestamppb.New(v.UpdatedAt),
		})
	}
	return &fileserver.FileListResponse{Files: result}, nil
}
