package datastore

import (
	"io"
	"os"
)

type Store interface {
	Save([]byte) error
	Load() ([]byte, error)
	Get(interface{}) ([]byte, error)
	Set(interface{}) error
	Close()
}

type FileStore struct {
	Path string
	File *os.File
}

func NewFileStore(fp string) (*FileStore, error) {
	f, err := os.OpenFile(fp, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		return nil, err
	}
	return &FileStore{Path: fp, File: f}, nil
}

func (fs *FileStore) Save(data []byte) error {
	fs.File.Truncate(0)
	_, err := fs.File.Write(data)
	fs.File.Sync()
	return err
}

func (fs *FileStore) Load() ([]byte, error) {
	return io.ReadAll(fs.File)
}

func (fs *FileStore) Get(key interface{}) ([]byte, error) {
	return nil, nil
}
func (fs *FileStore) Set(data interface{}) error {
	return nil
}

func (fs *FileStore) Close() {
	fs.File.Close()
}
