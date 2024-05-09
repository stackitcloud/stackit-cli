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

// readFileIfExists reads the contents of a file and returns it as a string.
// If the file does not exist, it returns an empty string.
func readFileIfExists(filePath string) (string, error) {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return "", nil
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("read file: %v", err)
	}

	return string(content), nil
}
