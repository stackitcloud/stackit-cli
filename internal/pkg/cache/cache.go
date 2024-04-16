package cache

import (
	"errors"
	"fmt"
	"os"
	"regexp"

	"github.com/adrg/xdg"
)

var (
	cacheFolderPath = xdg.CacheHome + "/stackit"

	identifierRegex             = regexp.MustCompile("^[a-zA-Z0-9-]+$")
	ErrorInvalidCacheIdentifier = fmt.Errorf("invalid cache identifier")
)

func GetObject(identifier string) ([]byte, error) {
	if !identifierRegex.MatchString(identifier) {
		return nil, ErrorInvalidCacheIdentifier
	}

	return os.ReadFile(cacheFolderPath + "/" + identifier)
}

func PutObject(identifier string, data []byte) error {
	if !identifierRegex.MatchString(identifier) {
		return ErrorInvalidCacheIdentifier
	}

	err := createFolderIfNotExists(cacheFolderPath)
	if err != nil {
		return err
	}

	return os.WriteFile(cacheFolderPath+"/"+identifier, data, 0o600)
}

func DeleteObject(identifier string) error {
	if !identifierRegex.MatchString(identifier) {
		return ErrorInvalidCacheIdentifier
	}

	if err := os.Remove(cacheFolderPath + "/" + identifier); !errors.Is(err, os.ErrNotExist) {
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
