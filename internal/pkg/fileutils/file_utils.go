package fileutils

import (
	"fmt"
	"os"
)

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
