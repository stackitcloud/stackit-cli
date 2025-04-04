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

	"github.com/zalando/go-keyring"
)

// Name of an auth-related field
type authFieldKey string

// Possible values of authentication flows
type AuthFlow string

const (
	keyringService     = "stackit-cli"
	textFileFolderName = "stackit"
	textFileName       = "cli-auth-storage.txt"
	envAccessTokenName = "STACKIT_ACCESS_TOKEN"
)

const (
	SESSION_EXPIRES_AT_UNIX authFieldKey = "session_expires_at_unix"
	ACCESS_TOKEN            authFieldKey = "access_token"
	REFRESH_TOKEN           authFieldKey = "refresh_token"
	SERVICE_ACCOUNT_TOKEN   authFieldKey = "service_account_token"
	SERVICE_ACCOUNT_EMAIL   authFieldKey = "service_account_email"
	USER_EMAIL              authFieldKey = "user_email"
	SERVICE_ACCOUNT_KEY     authFieldKey = "service_account_key"
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

func SetAuthFlow(value AuthFlow) error {
	return SetAuthField(authFlowType, string(value))
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

func SetAuthField(key authFieldKey, value string) error {
	activeProfile, err := config.GetProfile()
	if err != nil {
		return fmt.Errorf("get profile: %w", err)
	}

	return setAuthFieldWithProfile(activeProfile, key, value)
}

func setAuthFieldWithProfile(profile string, key authFieldKey, value string) error {
	err := setAuthFieldInKeyring(profile, key, value)
	if err != nil {
		errFallback := setAuthFieldInEncodedTextFile(profile, key, value)
		if errFallback != nil {
			return fmt.Errorf("write to keyring failed (%w), try writing to encoded text file: %w", err, errFallback)
		}
	}
	return nil
}

func setAuthFieldInKeyring(activeProfile string, key authFieldKey, value string) error {
	if activeProfile != config.DefaultProfileName {
		activeProfileKeyring := filepath.Join(keyringService, activeProfile)
		return keyring.Set(activeProfileKeyring, string(key), value)
	}
	return keyring.Set(keyringService, string(key), value)
}

func DeleteAuthField(key authFieldKey) error {
	activeProfile, err := config.GetProfile()
	if err != nil {
		return fmt.Errorf("get profile: %w", err)
	}
	return deleteAuthFieldWithProfile(activeProfile, key)
}

func deleteAuthFieldWithProfile(profile string, key authFieldKey) error {
	err := deleteAuthFieldInKeyring(profile, key)
	if err != nil {
		// if the key is not found, we can ignore the error
		if !errors.Is(err, keyring.ErrNotFound) {
			errFallback := deleteAuthFieldInEncodedTextFile(profile, key)
			if errFallback != nil {
				return fmt.Errorf("delete from keyring failed (%w), try deleting from encoded text file: %w", err, errFallback)
			}
		}
	}
	return nil
}

func deleteAuthFieldInEncodedTextFile(activeProfile string, key authFieldKey) error {
	err := createEncodedTextFile(activeProfile)
	if err != nil {
		return err
	}

	textFileDir := config.GetProfileFolderPath(activeProfile)
	textFilePath := filepath.Join(textFileDir, textFileName)

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
	keyringServiceLocal := keyringService
	if activeProfile != config.DefaultProfileName {
		keyringServiceLocal = filepath.Join(keyringService, activeProfile)
	}

	return keyring.Delete(keyringServiceLocal, string(key))
}

func setAuthFieldInEncodedTextFile(activeProfile string, key authFieldKey, value string) error {
	err := createEncodedTextFile(activeProfile)
	if err != nil {
		return err
	}
	textFileDir := config.GetProfileFolderPath(activeProfile)
	textFilePath := filepath.Join(textFileDir, textFileName)

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
	for key := range keyMap {
		value, err := GetAuthField(key)
		if err != nil {
			return fmt.Errorf("get auth field \"%s\": %w", key, err)
		}
		keyMap[key] = value
	}
	return nil
}

func GetAuthFlow() (AuthFlow, error) {
	value, err := GetAuthField(authFlowType)
	return AuthFlow(value), err
}

func GetAuthField(key authFieldKey) (string, error) {
	activeProfile, err := config.GetProfile()
	if err != nil {
		return "", fmt.Errorf("get profile: %w", err)
	}
	return getAuthFieldWithProfile(activeProfile, key)
}

func getAuthFieldWithProfile(profile string, key authFieldKey) (string, error) {
	value, err := getAuthFieldFromKeyring(profile, key)
	if err != nil {
		var errFallback error
		value, errFallback = getAuthFieldFromEncodedTextFile(profile, key)
		if errFallback != nil {
			return "", fmt.Errorf("read from keyring: %w, read from encoded file as fallback: %w", err, errFallback)
		}
	}
	return value, nil
}

func getAuthFieldFromKeyring(activeProfile string, key authFieldKey) (string, error) {
	if activeProfile != config.DefaultProfileName {
		activeProfileKeyring := filepath.Join(keyringService, activeProfile)
		return keyring.Get(activeProfileKeyring, string(key))
	}
	return keyring.Get(keyringService, string(key))
}

func getAuthFieldFromEncodedTextFile(activeProfile string, key authFieldKey) (string, error) {
	err := createEncodedTextFile(activeProfile)
	if err != nil {
		return "", err
	}

	textFileDir := config.GetProfileFolderPath(activeProfile)
	textFilePath := filepath.Join(textFileDir, textFileName)

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

// Checks if the encoded text file exist.
// If it doesn't, creates it with the content "{}" encoded.
// If it does, does nothing (and returns nil).
func createEncodedTextFile(activeProfile string) error {
	textFileDir := config.GetProfileFolderPath(activeProfile)
	textFilePath := filepath.Join(textFileDir, textFileName)

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

func LoginUser(email, accessToken, refreshToken, sessionExpiresAtUnix string) error {
	authFields := map[authFieldKey]string{
		SESSION_EXPIRES_AT_UNIX: sessionExpiresAtUnix,
		ACCESS_TOKEN:            accessToken,
		REFRESH_TOKEN:           refreshToken,
		USER_EMAIL:              email,
	}

	err := SetAuthFieldMap(authFields)
	if err != nil {
		return fmt.Errorf("set auth fields: %w", err)
	}
	return nil
}

func LogoutUser() error {
	for _, key := range loginAuthFieldKeys {
		err := DeleteAuthField(key)
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
