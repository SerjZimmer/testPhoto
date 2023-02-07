package storage

import (
	"bufio"
	"bytes"
	"log"
	"os"

	"github.com/pkg/errors"

	"github.com/SerjZimmer/testovoe1/internal/pkg/file"
)

func (s *DiskStorage) Get(ID string) (*file.File, error) {
	f, ok := s.files[ID]
	if !ok {
		return nil, &NotFoundErr{}
	}

	fileFromDisc, err := os.Open(f.Path)
	if err != nil {
		return nil, errors.Wrap(err, "cannot open")
	}
	defer func(fileFromDisc *os.File) {
		err := fileFromDisc.Close()
		if err != nil {
			log.Fatalln("Fail fileFromDisk.Close")
		}
	}(fileFromDisc)

	stat, err := fileFromDisc.Stat()
	if err != nil {
		return nil, errors.Wrap(err, "cannot get stat")
	}

	imageData := bytes.Buffer{}

	_, err = imageData.ReadFrom(fileFromDisc)
	if err != nil {
		return nil, errors.Wrap(err, "cannot read from")
	}

	return file.New(bufio.NewReader(&imageData), f.Name, f.Type, int(stat.Size())), nil

}
