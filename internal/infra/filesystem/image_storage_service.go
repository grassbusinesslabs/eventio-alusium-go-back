package filesystem

import (
	"log"
	"os"
	"path"

	"github.com/BohdanBoriak/boilerplate-go-back/config"
)

type ImageStorageService interface {
	SaveImage(filename string, content []byte) error
	DeleteImage(filename string) error
	GetImageContent(imgPath string) ([]byte, error)
}

type imageStorageService struct {
	loc string
}

func NewImageStorageService(conf config.Configuration) imageStorageService {
	return imageStorageService{
		loc: conf.FileStorageLocation,
	}
}

func (s imageStorageService) SaveImage(filename string, content []byte) error {
	location := path.Join(s.loc, filename)
	err := writeFileToStorage(location, content)
	log.Printf("Image saved at: %s", location)
	if err != nil {
		log.Printf("writeFileToStorage(imageStorageService.SaveImage): %s", err)
		return err
	}

	return nil
}

func (s imageStorageService) DeleteImage(filename string) error {
	location := path.Join(s.loc, filename)
	err := os.Remove(location)
	if err != nil {
		log.Printf("os.Remove(imageStorageService.DeleteImage): %s", err)
		return err
	}

	return nil
}

//nolint:gosec // need a permission greater than 0600 to read files on dev/prod
func writeFileToStorage(location string, file []byte) error {
	dirLocation := path.Dir(location)
	err := os.MkdirAll(dirLocation, os.ModePerm)
	if err != nil {
		return err
	}

	err = os.WriteFile(location, file, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (s imageStorageService) GetImageContent(imgPath string) ([]byte, error) {
	location := path.Join(s.loc, imgPath)
	content, err := os.ReadFile(location)
	if err != nil {
		log.Printf("os.ReadFile(imageStorageService.GetImageContent): %s", err)
		return nil, err
	}

	return content, nil
}
