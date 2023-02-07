package storage

import (
	"sync"
	"time"
)

type ImageRow struct {
	ID        string
	Name      string
	Type      string
	Path      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type NotFoundErr struct {
}

func (e *NotFoundErr) Error() string {
	return "not found"
}

type DiskStorage struct {
	dir   string
	mutex sync.Mutex
	files map[string]ImageRow
}

func New(dir string) *DiskStorage {
	return &DiskStorage{
		dir:   dir,
		mutex: sync.Mutex{},
		files: make(map[string]ImageRow),
	}
}
