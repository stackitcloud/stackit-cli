package flags

import (
	"fmt"
	"strings"

	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/pflag"
)

type uuidSliceFlag struct {
	value []string
}

// Ensure the implementation satisfies the expected interface
var _ pflag.Value = &uuidFlag{}

// UUIDSliceFlag returns a flag which must be a valid slice.
func UUIDSliceFlag() *uuidSliceFlag {
	return &uuidSliceFlag{}
}

func (f *uuidSliceFlag) String() string {
	return "[" + strings.Join(f.value, ",") + "]"
}

func (f *uuidSliceFlag) Set(value string) error {
	if value == "" {
		return fmt.Errorf("value cannot be empty")
	}

	uuids := strings.Split(value, ",")

	for i, uuid := range uuids {
		uuids[i] = strings.TrimSpace(uuid)

		err := utils.ValidateUUID(uuids[i])
		if err != nil {
			return err
		}
	}

	f.value = append(f.value, uuids...)
	return nil
}

func (f *uuidSliceFlag) Type() string {
	return "stringSlice"
}
