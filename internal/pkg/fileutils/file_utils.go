package fileutils

import (
	"fmt"
	"os"
)

// WriteToFile writes the given content to a file.
// If the file already exists, it will be overwritten.
func WriteToFile(outputFileName, content string) (err error) {
	fo, err := os.Create(outputFileName)
	if err != nil {
		return fmt.Errorf("create output file: %w", err)
	}
	defer func() {
		tempErr := fo.Close()
		if tempErr != nil {
			if err != nil {
				err = fmt.Errorf("%w; close output file: %w", err, tempErr)
			} else {
				err = fmt.Errorf("close output file: %w", tempErr)
			}
		}
	}()
	_, err = fo.WriteString(content)
	if err != nil {
		return fmt.Errorf("write content to output file: %w", err)
	}
	return err
}

// ReadFileIfExists reads the contents of a file and returns it as a string, along with a boolean indicating if the file exists.
// If the file does not exist, it returns an empty string, false and no error.
// If the file exists but cannot be read, it returns an error.
func ReadFileIfExists(filePath string) (contents string, exists bool, err error) {
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

// CopyFile copies the contents of a file to another file.
// If the destination file already exists, it will be overwritten.
func CopyFile(src, dst string) (err error) {
	contents, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("read source file: %w", err)
	}

	err = WriteToFile(dst, string(contents))
	if err != nil {
		return fmt.Errorf("write destination file: %w", err)
	}

	return nil
}
