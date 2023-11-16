package flags

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/pflag"
)

type readFromFileFlag struct {
	// Used to read file.
	// Set to os.ReadFile, except during tests
	reader func(filename string) ([]byte, error)
	value  string
}

// Ensure the implementation satisfies the expected interface
var _ pflag.Value = &readFromFileFlag{}

// ReadFromFileFlag returns a string flag.
// If it starts with "@", it is assumed to be a file path and content is read from file instead
func ReadFromFileFlag() *readFromFileFlag {
	return &readFromFileFlag{
		reader: os.ReadFile,
	}
}

func (f *readFromFileFlag) String() string {
	return f.value
}

func (f *readFromFileFlag) Set(value string) error {
	if !strings.HasPrefix(value, "@") {
		f.value = value
	} else {
		valuePath := strings.Trim(value[1:], `"'`)
		valueBytes, err := f.reader(valuePath)
		if err != nil {
			return fmt.Errorf("read data from file: %w", err)
		}
		f.value = string(valueBytes)
	}
	return nil
}

func (f *readFromFileFlag) Type() string {
	return "string"
}
