package utils

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/inhies/go-bytesize"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	sdkConfig "github.com/stackitcloud/stackit-sdk-go/core/config"
)

// Ptr Returns the pointer to any type T
func Ptr[T any](v T) *T {
	return &v
}

// PtrString creates a string representation of a passed object pointer or returns
// an empty string, if the passed object is _nil_.
func PtrString[T any](t *T) string {
	if t != nil {
		return fmt.Sprintf("%v", *t)
	}
	return ""
}

// PtrValue returns the dereferenced value if the pointer is not nil. Otherwise
// the types zero element is returned
func PtrValue[T any](t *T) (r T) {
	if t != nil {
		return *t
	}
	return r
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

func ValidateURLDomain(value string) error {
	urlStruct, err := url.Parse(value)
	if err != nil {
		return fmt.Errorf("parse url: %w", err)
	}
	urlHost := urlStruct.Hostname()
	if urlHost == "" {
		return fmt.Errorf("bad url")
	}

	allowedUrlDomain := viper.GetString(config.AllowedUrlDomainKey)

	if !strings.HasSuffix(urlHost, allowedUrlDomain) {
		return fmt.Errorf(`only urls belonging to domain %s are allowed`, allowedUrlDomain)
	}
	return nil
}

// ConvertTimePToDateTimeString converts a time.Time pointer to a string represented as "2006-01-02 15:04:05"
// This function will return an empty string if the input is nil
func ConvertTimePToDateTimeString(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.DateTime)
}

// PtrStringDefault return the value of a pointer [v] as string. If the pointer is nil, it returns the [defaultValue].
func PtrStringDefault[T any](v *T, defaultValue string) string {
	if v == nil {
		return defaultValue
	}
	return fmt.Sprintf("%v", *v)
}

// PtrByteSizeDefault return the value of an in64 pointer to a string representation of bytesize. If the pointer is nil,
// it returns the [defaultValue].
func PtrByteSizeDefault(size *int64, defaultValue string) string {
	if size == nil {
		return defaultValue
	}
	return bytesize.New(float64(*size)).String()
}

// Base64Encode encodes a []byte to a base64 representation as string
func Base64Encode(message []byte) string {
	b := make([]byte, base64.StdEncoding.EncodedLen(len(message)))
	base64.StdEncoding.Encode(b, message)
	return string(b)
}

func UserAgentConfigOption(cliVersion string) sdkConfig.ConfigurationOption {
	return sdkConfig.WithUserAgent(fmt.Sprintf("stackit-cli/%s", cliVersion))
}
