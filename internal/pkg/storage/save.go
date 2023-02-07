package storage

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"

	"github.com/SerjZimmer/testovoe1/internal/pkg/file"
)

func (s *DiskStorage) Save(file *file.File) (string, error) {
	imageID, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("cannot generate image id: %w", err)
	}
	imagePath := fmt.Sprintf("%s/%s.%s", s.dir, imageID, file.Extension)

	createdFile, err := os.Create(filepath.Clean(imagePath))
	if err != nil {
		return "", fmt.Errorf("cannot create image file: %w", err)
	}
	defer func(createdFile *os.File) {
		err := createdFile.Close()
		if err != nil {
			log.Fatalln("Fail createdFile.Close")
		}
	}(createdFile)

	_, err = createdFile.ReadFrom(file)
	if err != nil {
		return "", errors.New("cannot read to file")
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.files[imageID.String()] = ImageRow{
		ID:        imageID.String(),
		Name:      file.Name,
		Type:      file.Extension,
		Path:      imagePath,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return imageID.String(), nil
}
