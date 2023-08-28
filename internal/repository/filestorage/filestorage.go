package filestorage

import (
	"github.com/Baraha/gimmick-service/internal/config"
	"github.com/Baraha/gimmick-service/internal/services"
	"io"
	"log"
	"os"
)

type FileStorage struct {
	cfg config.Configuration
}

func NewFileStorage(cfg config.Configuration) services.IFileStorage {
	return &FileStorage{cfg}
}

func (fs *FileStorage) SetFile(fileFormat services.Ifile, fileName string) error {
	file, err := fileFormat.Open()
	defer file.Close()

	dst, err := os.Create(fs.cfg.FileStorage.DefaultDir + fileName)
	if err != nil {
		log.Println("error creating file", err)
		return err
	}
	defer dst.Close()
	if _, err := io.Copy(dst, file); err != nil {
		return err
	}

	return nil
}

func (fs *FileStorage) DefaultDir() string {
	return fs.cfg.FileStorage.DefaultDir
}

func (fs *FileStorage) DeleteFile(fileName string) error {
	return os.Remove(fs.cfg.FileStorage.DefaultDir + fileName)
}
