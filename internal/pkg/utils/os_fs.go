package utils

import (
	"io/fs"
	"os"
)

var _ fs.ReadFileFS = OsFS{}

type OsFS struct{}

func (o OsFS) Open(name string) (fs.File, error) {
	return os.Open(name)
}

func (o OsFS) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}
