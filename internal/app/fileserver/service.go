package fileserver

import (
	"github.com/SerjZimmer/testovoe1/internal/pkg/file"
	"github.com/SerjZimmer/testovoe1/internal/pkg/storage"
	"github.com/SerjZimmer/testovoe1/pkg/fileserver"
)

// UnimplementedFileServiceServer must be embedded to have forward compatible implementations.
type fileStore interface {
	Save(file *file.File) (string, error)
	Get(ID string) (*file.File, error)
	List() ([]storage.ImageRow, error)
}

type Server struct {
	fileserver.UnimplementedFileServiceServer
	storage fileStore
}

func New(storage fileStore) *Server {
	return &Server{
		storage: storage,
	}
}
