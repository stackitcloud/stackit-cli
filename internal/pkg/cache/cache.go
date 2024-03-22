package cache

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

const (
	cacheFolder = ".stackit/cache"
)

var ErrorInvalidCacheIdentifier = fmt.Errorf("invalid cache identifier")

var identifierRegex = regexp.MustCompile("^[a-zA-Z0-9-]+$")

func GetObject(identifier string) ([]byte, error) {
	if !identifierRegex.MatchString(identifier) {
		return nil, ErrorInvalidCacheIdentifier
	}

	cacheFolderPath, err := getCacheFolderPath()
	if err != nil {
		return nil, err
	}

	return os.ReadFile(cacheFolderPath + "/" + identifier + ".txt")
}

func PutObject(identifier string, data []byte) error {
	if !identifierRegex.MatchString(identifier) {
		return ErrorInvalidCacheIdentifier
	}

	cacheFolderPath, err := getCacheFolderPath()
	if err != nil {
		return err
	}

	err = createFolderIfNotExists(cacheFolderPath)
	if err != nil {
		return err
	}

	return os.WriteFile(cacheFolderPath+"/"+identifier+".txt", data, 0o600)
}

func DeleteObject(identifier string) error {
	if !identifierRegex.MatchString(identifier) {
		return ErrorInvalidCacheIdentifier
	}

	cacheFolderPath, err := getCacheFolderPath()
	if err != nil {
		return err
	}

	if err = os.Remove(cacheFolderPath + "/" + identifier + ".txt"); !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}

func getCacheFolderPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configFolderPath := filepath.Join(home, cacheFolder)
	return configFolderPath, nil
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
