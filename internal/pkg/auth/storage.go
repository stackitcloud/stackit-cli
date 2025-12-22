package auth

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	pkgErrors "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"

	"github.com/zalando/go-keyring"
)

// Package-level printer for debug logging in storage operations
var storagePrinter = print.NewPrinter() //nolint:unused // set via SetStoragePrinter, may be used for future debug logging

// SetStoragePrinter sets the printer used for storage debug logging
// This should be called with the main command's printer to ensure consistent verbosity
func SetStoragePrinter(p *print.Printer) {
	if p != nil {
		storagePrinter = p
	}
}

// Name of an auth-related field
type authFieldKey string

// Possible values of authentication flows
type AuthFlow string

// StorageContext represents the context in which credentials are stored
// CLI context is for the CLI's own authentication
// API context is for Terraform Provider and SDK authentication
type StorageContext string

const (
	StorageContextCLI StorageContext = "cli"
	StorageContextAPI StorageContext = "api"
)

const (
	keyringServiceCLI  = "stackit-cli"
	keyringServiceAPI  = "stackit-cli-api"
	textFileNameCLI    = "cli-auth-storage.txt"
	textFileNameAPI    = "cli-api-auth-storage.txt"
	envAccessTokenName = "STACKIT_ACCESS_TOKEN"
)

const (
	SESSION_EXPIRES_AT_UNIX authFieldKey = "session_expires_at_unix"
	ACCESS_TOKEN            authFieldKey = "access_token"
	REFRESH_TOKEN           authFieldKey = "refresh_token"
	SERVICE_ACCOUNT_TOKEN   authFieldKey = "service_account_token"
	SERVICE_ACCOUNT_EMAIL   authFieldKey = "service_account_email"
	USER_EMAIL              authFieldKey = "user_email"
	SERVICE_ACCOUNT_KEY     authFieldKey = "service_account_key" //nolint:gosec // linter false positive
	PRIVATE_KEY             authFieldKey = "private_key"
	TOKEN_CUSTOM_ENDPOINT   authFieldKey = "token_custom_endpoint"
	IDP_TOKEN_ENDPOINT      authFieldKey = "idp_token_endpoint" //nolint:gosec // linter false positive
)

const (
	authFlowType                    authFieldKey = "auth_flow_type"
	AUTH_FLOW_USER_TOKEN            AuthFlow     = "user_token"
	AUTH_FLOW_SERVICE_ACCOUNT_TOKEN AuthFlow     = "sa_token"
	AUTH_FLOW_SERVICE_ACCOUNT_KEY   AuthFlow     = "sa_key"
)

// Returns all auth field keys managed by the auth storage
var authFieldKeys = []authFieldKey{
	SESSION_EXPIRES_AT_UNIX,
	ACCESS_TOKEN,
	REFRESH_TOKEN,
	SERVICE_ACCOUNT_TOKEN,
	SERVICE_ACCOUNT_EMAIL,
	USER_EMAIL,
	SERVICE_ACCOUNT_KEY,
	PRIVATE_KEY,
	TOKEN_CUSTOM_ENDPOINT,
	IDP_TOKEN_ENDPOINT,
	authFlowType,
}

// All fields that are set when a user logs in
// These fields should match the ones in LoginUser, which is ensured by the tests
var loginAuthFieldKeys = []authFieldKey{
	SESSION_EXPIRES_AT_UNIX,
	ACCESS_TOKEN,
	REFRESH_TOKEN,
	USER_EMAIL,
}

// getKeyringServiceName returns the keyring service name for the given context and profile
func getKeyringServiceName(context StorageContext, profile string) string {
	var baseService string
	switch context {
	case StorageContextAPI:
		baseService = keyringServiceAPI
	default:
		baseService = keyringServiceCLI
	}

	if profile != config.DefaultProfileName {
		return filepath.Join(baseService, profile)
	}
	return baseService
}

// getTextFileName returns the text file name for the given context
func getTextFileName(context StorageContext) string {
	switch context {
	case StorageContextAPI:
		return textFileNameAPI
	default:
		return textFileNameCLI
	}
}

func SetAuthFlow(value AuthFlow) error {
	return SetAuthField(authFlowType, string(value))
}

func SetAuthFlowWithContext(context StorageContext, value AuthFlow) error {
	return SetAuthFieldWithContext(context, authFlowType, string(value))
}

// Sets the values in the auth storage according to the given map
func SetAuthFieldMap(keyMap map[authFieldKey]string) error {
	for key, value := range keyMap {
		err := SetAuthField(key, value)
		if err != nil {
			return fmt.Errorf("set auth field \"%s\": %w", key, err)
		}
	}
	return nil
}

// SetAuthFieldMapWithContext sets the values in the auth storage according to the given map for a specific context
func SetAuthFieldMapWithContext(context StorageContext, keyMap map[authFieldKey]string) error {
	for key, value := range keyMap {
		err := SetAuthFieldWithContext(context, key, value)
		if err != nil {
			return fmt.Errorf("set auth field \"%s\": %w", key, err)
		}
	}
	return nil
}

func SetAuthField(key authFieldKey, value string) error {
	return SetAuthFieldWithContext(StorageContextCLI, key, value)
}

// SetAuthFieldWithContext sets an auth field for a specific storage context
func SetAuthFieldWithContext(context StorageContext, key authFieldKey, value string) error {
	activeProfile, err := config.GetProfile()
	if err != nil {
		return fmt.Errorf("get profile: %w", err)
	}

	return setAuthFieldWithProfileAndContext(context, activeProfile, key, value)
}

func setAuthFieldWithProfile(profile string, key authFieldKey, value string) error {
	return setAuthFieldWithProfileAndContext(StorageContextCLI, profile, key, value)
}

func setAuthFieldWithProfileAndContext(context StorageContext, profile string, key authFieldKey, value string) error {
	err := setAuthFieldInKeyringWithContext(context, profile, key, value)
	if err != nil {
		errFallback := setAuthFieldInEncodedTextFileWithContext(context, profile, key, value)
		if errFallback != nil {
			return fmt.Errorf("write to keyring failed (%w), try writing to encoded text file: %w", err, errFallback)
		}
	}
	return nil
}

func setAuthFieldInKeyring(activeProfile string, key authFieldKey, value string) error {
	return setAuthFieldInKeyringWithContext(StorageContextCLI, activeProfile, key, value)
}

func setAuthFieldInKeyringWithContext(context StorageContext, activeProfile string, key authFieldKey, value string) error {
	keyringServiceName := getKeyringServiceName(context, activeProfile)
	return keyring.Set(keyringServiceName, string(key), value)
}

func DeleteAuthField(key authFieldKey) error {
	return DeleteAuthFieldWithContext(StorageContextCLI, key)
}

// DeleteAuthFieldWithContext deletes an auth field for a specific storage context
func DeleteAuthFieldWithContext(context StorageContext, key authFieldKey) error {
	activeProfile, err := config.GetProfile()
	if err != nil {
		return fmt.Errorf("get profile: %w", err)
	}
	return deleteAuthFieldWithProfileAndContext(context, activeProfile, key)
}

func deleteAuthFieldWithProfile(profile string, key authFieldKey) error {
	return deleteAuthFieldWithProfileAndContext(StorageContextCLI, profile, key)
}

func deleteAuthFieldWithProfileAndContext(context StorageContext, profile string, key authFieldKey) error {
	err := deleteAuthFieldInKeyringWithContext(context, profile, key)
	if err != nil {
		// if the key is not found, we can ignore the error
		if !errors.Is(err, keyring.ErrNotFound) {
			errFallback := deleteAuthFieldInEncodedTextFileWithContext(context, profile, key)
			if errFallback != nil {
				return fmt.Errorf("delete from keyring failed (%w), try deleting from encoded text file: %w", err, errFallback)
			}
		}
	}
	return nil
}

func deleteAuthFieldInEncodedTextFile(activeProfile string, key authFieldKey) error {
	return deleteAuthFieldInEncodedTextFileWithContext(StorageContextCLI, activeProfile, key)
}

func deleteAuthFieldInEncodedTextFileWithContext(context StorageContext, activeProfile string, key authFieldKey) error {
	err := createEncodedTextFileWithContext(context, activeProfile)
	if err != nil {
		return err
	}

	textFileDir := config.GetProfileFolderPath(activeProfile)
	fileName := getTextFileName(context)
	textFilePath := filepath.Join(textFileDir, fileName)

	contentEncoded, err := os.ReadFile(textFilePath)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}
	contentBytes, err := base64.StdEncoding.DecodeString(string(contentEncoded))
	if err != nil {
		return fmt.Errorf("decode file: %w", err)
	}
	content := map[authFieldKey]string{}
	err = json.Unmarshal(contentBytes, &content)
	if err != nil {
		return fmt.Errorf("unmarshal file: %w", err)
	}

	delete(content, key)

	contentBytes, err = json.Marshal(content)
	if err != nil {
		return fmt.Errorf("marshal file: %w", err)
	}
	contentEncoded = []byte(base64.StdEncoding.EncodeToString(contentBytes))
	err = os.WriteFile(textFilePath, contentEncoded, 0o600)
	if err != nil {
		return fmt.Errorf("write file: %w", err)
	}
	return nil
}

func deleteAuthFieldInKeyring(activeProfile string, key authFieldKey) error {
	return deleteAuthFieldInKeyringWithContext(StorageContextCLI, activeProfile, key)
}

func deleteAuthFieldInKeyringWithContext(context StorageContext, activeProfile string, key authFieldKey) error {
	keyringServiceName := getKeyringServiceName(context, activeProfile)
	return keyring.Delete(keyringServiceName, string(key))
}

func setAuthFieldInEncodedTextFile(activeProfile string, key authFieldKey, value string) error {
	return setAuthFieldInEncodedTextFileWithContext(StorageContextCLI, activeProfile, key, value)
}

func setAuthFieldInEncodedTextFileWithContext(context StorageContext, activeProfile string, key authFieldKey, value string) error {
	textFileDir := config.GetProfileFolderPath(activeProfile)
	fileName := getTextFileName(context)
	textFilePath := filepath.Join(textFileDir, fileName)

	err := createEncodedTextFileWithContext(context, activeProfile)
	if err != nil {
		return err
	}

	contentEncoded, err := os.ReadFile(textFilePath)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}
	contentBytes, err := base64.StdEncoding.DecodeString(string(contentEncoded))
	if err != nil {
		return fmt.Errorf("decode file: %w", err)
	}
	content := map[authFieldKey]string{}
	err = json.Unmarshal(contentBytes, &content)
	if err != nil {
		return fmt.Errorf("unmarshal file: %w", err)
	}

	content[key] = value

	contentBytes, err = json.Marshal(content)
	if err != nil {
		return fmt.Errorf("marshal file: %w", err)
	}
	contentEncoded = []byte(base64.StdEncoding.EncodeToString(contentBytes))
	err = os.WriteFile(textFilePath, contentEncoded, 0o600)
	if err != nil {
		return fmt.Errorf("write file: %w", err)
	}
	return nil
}

// Populates the values in the given map according to the auth storage
func GetAuthFieldMap(keyMap map[authFieldKey]string) error {
	return GetAuthFieldMapWithContext(StorageContextCLI, keyMap)
}

// GetAuthFieldMapWithContext populates the values in the given map according to the auth storage for a specific context
func GetAuthFieldMapWithContext(context StorageContext, keyMap map[authFieldKey]string) error {
	for key := range keyMap {
		value, err := GetAuthFieldWithContext(context, key)
		if err != nil {
			return fmt.Errorf("get auth field \"%s\": %w", key, err)
		}
		keyMap[key] = value
	}
	return nil
}

func GetAuthFlow() (AuthFlow, error) {
	return GetAuthFlowWithContext(StorageContextCLI)
}

func GetAuthFlowWithContext(context StorageContext) (AuthFlow, error) {
	value, err := GetAuthFieldWithContext(context, authFlowType)
	return AuthFlow(value), err
}

func GetAuthField(key authFieldKey) (string, error) {
	return GetAuthFieldWithContext(StorageContextCLI, key)
}

// GetAuthFieldWithContext retrieves an auth field for a specific storage context
func GetAuthFieldWithContext(context StorageContext, key authFieldKey) (string, error) {
	activeProfile, err := config.GetProfile()
	if err != nil {
		return "", fmt.Errorf("get profile: %w", err)
	}
	return getAuthFieldWithProfileAndContext(context, activeProfile, key)
}

func getAuthFieldWithProfile(profile string, key authFieldKey) (string, error) {
	return getAuthFieldWithProfileAndContext(StorageContextCLI, profile, key)
}

func getAuthFieldWithProfileAndContext(context StorageContext, profile string, key authFieldKey) (string, error) {
	value, err := getAuthFieldFromKeyringWithContext(context, profile, key)
	if err != nil {
		var errFallback error
		value, errFallback = getAuthFieldFromEncodedTextFileWithContext(context, profile, key)
		if errFallback != nil {
			return "", fmt.Errorf("read from keyring: %w, read from encoded file as fallback: %w", err, errFallback)
		}
	}
	return value, nil
}

func getAuthFieldFromKeyring(activeProfile string, key authFieldKey) (string, error) {
	return getAuthFieldFromKeyringWithContext(StorageContextCLI, activeProfile, key)
}

func getAuthFieldFromKeyringWithContext(context StorageContext, activeProfile string, key authFieldKey) (string, error) {
	keyringServiceName := getKeyringServiceName(context, activeProfile)
	return keyring.Get(keyringServiceName, string(key))
}

func getAuthFieldFromEncodedTextFile(activeProfile string, key authFieldKey) (string, error) {
	return getAuthFieldFromEncodedTextFileWithContext(StorageContextCLI, activeProfile, key)
}

func getAuthFieldFromEncodedTextFileWithContext(context StorageContext, activeProfile string, key authFieldKey) (string, error) {
	err := createEncodedTextFileWithContext(context, activeProfile)
	if err != nil {
		return "", err
	}

	textFileDir := config.GetProfileFolderPath(activeProfile)
	fileName := getTextFileName(context)
	textFilePath := filepath.Join(textFileDir, fileName)

	contentEncoded, err := os.ReadFile(textFilePath)
	if err != nil {
		return "", fmt.Errorf("read file: %w", err)
	}
	contentBytes, err := base64.StdEncoding.DecodeString(string(contentEncoded))
	if err != nil {
		return "", fmt.Errorf("decode file: %w", err)
	}
	var content map[authFieldKey]string
	err = json.Unmarshal(contentBytes, &content)
	if err != nil {
		return "", fmt.Errorf("unmarshal file: %w", err)
	}
	value, ok := content[key]
	if !ok {
		return "", fmt.Errorf("value not found")
	}
	return value, nil
}

// createEncodedTextFileWithContext checks if the encoded text file exist.
// If it doesn't, creates it with the content "{}" encoded.
// If it does, does nothing (and returns nil).
func createEncodedTextFileWithContext(context StorageContext, activeProfile string) error {
	textFileDir := config.GetProfileFolderPath(activeProfile)
	fileName := getTextFileName(context)
	textFilePath := filepath.Join(textFileDir, fileName)

	err := os.MkdirAll(textFileDir, 0o750)
	if err != nil {
		return fmt.Errorf("create file dir: %w", err)
	}
	_, err = os.Stat(textFilePath)
	if !os.IsNotExist(err) {
		return nil
	}

	contentEncoded := base64.StdEncoding.EncodeToString([]byte("{}"))
	err = os.WriteFile(textFilePath, []byte(contentEncoded), 0o600)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}

	return nil
}

// GetProfileEmail returns the email of the user or service account associated with the given profile.
// If the profile is not authenticated or the email can't be obtained, it returns an empty string.
func GetProfileEmail(profile string) string {
	value, err := getAuthFieldWithProfile(profile, authFlowType)
	if err != nil {
		return ""
	}

	var email string
	switch AuthFlow(value) {
	case AUTH_FLOW_USER_TOKEN:
		email, err = getAuthFieldWithProfile(profile, USER_EMAIL)
		if err != nil {
			email = ""
		}
	case AUTH_FLOW_SERVICE_ACCOUNT_TOKEN, AUTH_FLOW_SERVICE_ACCOUNT_KEY:
		email, err = getAuthFieldWithProfile(profile, SERVICE_ACCOUNT_EMAIL)
		if err != nil {
			email = ""
		}
	}
	return email
}

// GetAuthEmail returns the email of the authenticated account.
// If the environment variable STACKIT_ACCESS_TOKEN is set, the email of this token will be returned.
func GetAuthEmail() (string, error) {
	// If STACKIT_ACCESS_TOKEN is set, get the mail from the token
	if accessToken := os.Getenv(envAccessTokenName); accessToken != "" {
		email, err := getEmailFromToken(accessToken)
		if err != nil {
			return "", fmt.Errorf("error getting email from token: %w", err)
		}
		return email, nil
	}

	profile, err := config.GetProfile()
	if err != nil {
		return "", fmt.Errorf("error getting profile: %w", err)
	}
	email := GetProfileEmail(profile)
	if email == "" {
		return "", fmt.Errorf("error getting profile email. email is empty")
	}
	return email, nil
}

func LoginUser(email, accessToken, refreshToken, sessionExpiresAtUnix string) error {
	return LoginUserWithContext(StorageContextCLI, email, accessToken, refreshToken, sessionExpiresAtUnix)
}

// LoginUserWithContext stores user login credentials for a specific storage context
func LoginUserWithContext(context StorageContext, email, accessToken, refreshToken, sessionExpiresAtUnix string) error {
	authFields := map[authFieldKey]string{
		SESSION_EXPIRES_AT_UNIX: sessionExpiresAtUnix,
		ACCESS_TOKEN:            accessToken,
		REFRESH_TOKEN:           refreshToken,
		USER_EMAIL:              email,
	}

	err := SetAuthFieldMapWithContext(context, authFields)
	if err != nil {
		return fmt.Errorf("set auth fields: %w", err)
	}
	return nil
}

func LogoutUser() error {
	return LogoutUserWithContext(StorageContextCLI)
}

// LogoutUserWithContext removes user authentication for a specific storage context
func LogoutUserWithContext(context StorageContext) error {
	for _, key := range loginAuthFieldKeys {
		err := DeleteAuthFieldWithContext(context, key)
		if err != nil {
			return fmt.Errorf("delete auth field \"%s\": %w", key, err)
		}
	}
	return nil
}

func DeleteProfileAuth(profile string) error {
	err := config.ValidateProfile(profile)
	if err != nil {
		return fmt.Errorf("validate profile: %w", err)
	}

	if profile == config.DefaultProfileName {
		return &pkgErrors.DeleteDefaultProfile{DefaultProfile: config.DefaultProfileName}
	}

	for _, key := range authFieldKeys {
		err := deleteAuthFieldWithProfile(profile, key)
		if err != nil {
			return fmt.Errorf("delete auth field \"%s\": %w", key, err)
		}
	}

	return nil
}
