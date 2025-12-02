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
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
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

// PtrGigaByteSizeDefault return the value of an int64 pointer to a string representation of gigabytes. If the pointer is nil,
// it returns the [defaultValue].
func PtrGigaByteSizeDefault(size *int64, defaultValue string) string {
	if size == nil {
		return defaultValue
	}
	return (bytesize.New(float64(*size)) * bytesize.GB).String()
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

// ConvertStringMapToInterfaceMap converts a map[string]string to a pointer to map[string]interface{}.
// Returns nil if the input map is empty.
//
//nolint:gocritic // Linter wants to have a non-pointer type for the map, but this would mean a nil check has to be done before every usage of this func.
func ConvertStringMapToInterfaceMap(m *map[string]string) *map[string]interface{} {
	if m == nil || len(*m) == 0 {
		return nil
	}
	result := make(map[string]interface{}, len(*m))
	for k, v := range *m {
		result[k] = v
	}
	return &result
}

// Base64Bytes implements yaml.Marshaler to convert []byte to base64 strings
// ref: https://carlosbecker.com/posts/go-custom-marshaling
type Base64Bytes []byte

// MarshalYAML implements yaml.Marshaler
func (b Base64Bytes) MarshalYAML() (interface{}, error) {
	if len(b) == 0 {
		return "", nil
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

type Base64PatchedServer struct {
	Id                  *string                 `json:"id,omitempty"`
	Name                *string                 `json:"name,omitempty"`
	Status              *string                 `json:"status,omitempty"`
	AvailabilityZone    *string                 `json:"availabilityZone,omitempty"`
	BootVolume          *iaas.ServerBootVolume  `json:"bootVolume,omitempty"`
	CreatedAt           *time.Time              `json:"createdAt,omitempty"`
	ErrorMessage        *string                 `json:"errorMessage,omitempty"`
	PowerStatus         *string                 `json:"powerStatus,omitempty"`
	AffinityGroup       *string                 `json:"affinityGroup,omitempty"`
	ImageId             *string                 `json:"imageId,omitempty"`
	KeypairName         *string                 `json:"keypairName,omitempty"`
	MachineType         *string                 `json:"machineType,omitempty"`
	Labels              *map[string]interface{} `json:"labels,omitempty"`
	LaunchedAt          *time.Time              `json:"launchedAt,omitempty"`
	MaintenanceWindow   *iaas.ServerMaintenance `json:"maintenanceWindow,omitempty"`
	Metadata            *map[string]interface{} `json:"metadata,omitempty"`
	Networking          *iaas.ServerNetworking  `json:"networking,omitempty"`
	Nics                *[]iaas.ServerNetwork   `json:"nics,omitempty"`
	SecurityGroups      *[]string               `json:"securityGroups,omitempty"`
	ServiceAccountMails *[]string               `json:"serviceAccountMails,omitempty"`
	UpdatedAt           *time.Time              `json:"updatedAt,omitempty"`
	UserData            *Base64Bytes            `json:"userData,omitempty"`
	Volumes             *[]string               `json:"volumes,omitempty"`
}

// ConvertToBase64PatchedServer converts an iaas.Server to Base64PatchedServer
// This is a temporary workaround to get the desired base64 encoded yaml output for userdata
// and will be replaced by a fix in the Go-SDK
// ref: https://jira.schwarz/browse/STACKITSDK-246
func ConvertToBase64PatchedServer(server *iaas.Server) *Base64PatchedServer {
	if server == nil {
		return nil
	}

	var userData *Base64Bytes
	if server.UserData != nil {
		userData = Ptr(Base64Bytes(*server.UserData))
	}

	return &Base64PatchedServer{
		Id:                  server.Id,
		Name:                server.Name,
		Status:              server.Status,
		AvailabilityZone:    server.AvailabilityZone,
		BootVolume:          server.BootVolume,
		CreatedAt:           server.CreatedAt,
		ErrorMessage:        server.ErrorMessage,
		PowerStatus:         server.PowerStatus,
		AffinityGroup:       server.AffinityGroup,
		ImageId:             server.ImageId,
		KeypairName:         server.KeypairName,
		MachineType:         server.MachineType,
		Labels:              server.Labels,
		LaunchedAt:          server.LaunchedAt,
		MaintenanceWindow:   server.MaintenanceWindow,
		Metadata:            server.Metadata,
		Networking:          server.Networking,
		Nics:                server.Nics,
		SecurityGroups:      server.SecurityGroups,
		ServiceAccountMails: server.ServiceAccountMails,
		UpdatedAt:           server.UpdatedAt,
		UserData:            userData,
		Volumes:             server.Volumes,
	}
}

// ConvertToBase64PatchedServers converts a slice of iaas.Server to a slice of Base64PatchedServer
// This is a temporary workaround to get the desired base64 encoded yaml output for userdata
// and will be replaced by a fix in the Go-SDK
// ref: https://jira.schwarz/browse/STACKITSDK-246
func ConvertToBase64PatchedServers(servers []iaas.Server) []Base64PatchedServer {
	if servers == nil {
		return nil
	}

	result := make([]Base64PatchedServer, len(servers))
	for i := range servers {
		result[i] = *ConvertToBase64PatchedServer(&servers[i])
	}

	return result
}

// GetSliceFromPointer returns the value of a pointer to a slice of type T.
// If the pointer is nil, it returns an empty slice.
func GetSliceFromPointer[T any](s *[]T) []T {
	if s == nil || *s == nil {
		return []T{}
	}
	return *s
}
