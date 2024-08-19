package utils

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
)

const (
	defaultAllowedUrlDomain = "stackit.cloud"
)

// Ptr Returns the pointer to any type T
func Ptr[T any](v T) *T {
	return &v
}

// Int64Ptr returns a pointer to an int64
// Needed because the Ptr function only returns pointer to int
func Int64Ptr(i int64) *int64 {
	return &i
}

// Float64Ptr returns a pointer to a float64
// Needed because the Ptr function only returns pointer to float
func Float64Ptr(f float64) *float64 {
	return &f
}

// CmdHelp is used to explicitly set the Run function for non-leaf commands to the command help function, so that we can catch invalid commands
// This is a workaround needed due to the open issue on the Cobra repo: https://github.com/spf13/cobra/issues/706
func CmdHelp(cmd *cobra.Command, _ []string) {
	cmd.Help() //nolint:errcheck //the function doesnt return anything to satisfy the required interface of the Run function
}

// ValidateUUID validates if the provided string is a valid UUID
func ValidateUUID(value string) error {
	_, err := uuid.Parse(value)
	if err != nil {
		return fmt.Errorf("parse %s as UUID: %w", value, err)
	}
	return nil
}

// ConvertInt64PToFloat64P converts an int64 pointer to a float64 pointer
// This function will return nil if the input is nil
func ConvertInt64PToFloat64P(i *int64) *float64 {
	if i == nil {
		return nil
	}
	f := float64(*i)
	return &f
}

func ValidateURL(value string) error {
	urlStruct, err := url.Parse(value)
	if err != nil {
		return fmt.Errorf("parse url: %w", err)
	}
	urlHost := urlStruct.Hostname()
	if urlHost == "" {
		return fmt.Errorf("bad url")
	}

	allowedUrlDomain := viper.GetString(config.AllowedUrlDomainKey)

	if allowedUrlDomain == "" {
		allowedUrlDomain = defaultAllowedUrlDomain
	}

	if !strings.HasSuffix(urlHost, allowedUrlDomain) {
		return fmt.Errorf(`only urls belonging to domain %s are allowed"`, allowedUrlDomain)
	}
	return nil
}
