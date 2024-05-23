package fileutils

import (
	"fmt"
	"os"
)

func FileOutput(outputFileName, content string) error {
	fo, err := os.Create(outputFileName)
	if err != nil {
		return fmt.Errorf("create output file: %w", err)
	}

	defer func() {
		tempErr := fo.Close()
		if tempErr != nil {
			err = fmt.Errorf("close output file: %w", tempErr)
		}
	}()

	_, err = fo.WriteString(content)
	if err != nil {
		return fmt.Errorf("write content to output file: %w", err)
	}

	return nil
}
