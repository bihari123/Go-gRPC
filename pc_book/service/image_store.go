package service

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/google/uuid"
)

type ImageStore interface {
	Save(laptopId, imageType string, imageData bytes.Buffer) (string, error)
}

type DiskImageStore struct {
	mutex       sync.RWMutex
	imageFolder string
	images      map[string]*ImageInfo
}

type ImageInfo struct {
	LaptopId string
	Type     string
	Path     string
}

func NewDiskImageStore(imageFolder string) *DiskImageStore {
	return &DiskImageStore{
		imageFolder: imageFolder,
		images:      make(map[string]*ImageInfo),
	}
}

func (store *DiskImageStore) Save(
	laptopId string,
	imagetype string,
	imageData bytes.Buffer,
) (string, error) {
	imageId, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("cannot generate image id: %w", err)
	}

	imagePath := filepath.Join(store.imageFolder, fmt.Sprintf("%v.%v", imageId, imagetype))
	file, err := os.Create(imagePath)
	if err != nil {
		return "", fmt.Errorf("cannot make image path %v : %v", imagePath, err)
	}

	_, err = imageData.WriteTo(file)

	store.mutex.Lock()
	defer store.mutex.Unlock()

	store.images[imageId.String()] = &ImageInfo{
		LaptopId: laptopId,
		Type:     imagetype,
		Path:     imagePath,
	}

	return imageId.String(), nil
}
