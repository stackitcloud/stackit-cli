package config

import (
	"fmt"
	"os"
)

// createFolderIfNotExists creates a folder if it does not exist.
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

// readFileIfExists reads the contents of a file and returns it as a string, along with a boolean indicating if the file exists.
// If the file does not exist, it returns an empty string and no error.
// If the file exists but cannot be read, it returns an error.
func readFileIfExists(filePath string) (contents string, exists bool, err error) {
	_, err = os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", false, nil
		}
		return "", true, err
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", true, fmt.Errorf("read file: %w", err)
	}

	return string(content), true, nil
}
