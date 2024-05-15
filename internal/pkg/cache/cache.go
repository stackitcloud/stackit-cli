package cache

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

var (
	cacheFolderPath string

	identifierRegex             = regexp.MustCompile("^[a-zA-Z0-9-]+$")
	ErrorInvalidCacheIdentifier = fmt.Errorf("invalid cache identifier")
)

func init() {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		panic(fmt.Errorf("get user cache dir: %w", err))
	}
	cacheFolderPath = filepath.Join(cacheDir, "stackit")
}

func GetObject(identifier string) ([]byte, error) {
	if !identifierRegex.MatchString(identifier) {
		return nil, ErrorInvalidCacheIdentifier
	}

	return os.ReadFile(filepath.Join(cacheFolderPath, identifier))
}

func PutObject(identifier string, data []byte) error {
	if !identifierRegex.MatchString(identifier) {
		return ErrorInvalidCacheIdentifier
	}

	err := createFolderIfNotExists(cacheFolderPath)
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(cacheFolderPath, identifier), data, 0o600)
}

func DeleteObject(identifier string) error {
	if !identifierRegex.MatchString(identifier) {
		return ErrorInvalidCacheIdentifier
	}

	if err := os.Remove(filepath.Join(cacheFolderPath, identifier)); !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}

func createFolderIfNotExists(folderPath string) error {
	_, err := os.Stat(folderPath)
	if os.IsNotExist(err) {
		err := os.MkdirAll(folderPath, os.ModePerm)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	return nil
}
