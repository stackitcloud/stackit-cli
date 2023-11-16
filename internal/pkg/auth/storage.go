package auth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/zalando/go-keyring"
)

// Name of an auth-related field
type authFieldKey string

// Possible values of authentication flows
type AuthFlow string

const (
	keyringService     = "stackit-cli"
	textFileFolderName = ".stackit"
	textFileName       = "cli-auth-storage.txt"
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
	JWKS_CUSTOM_ENDPOINT    authFieldKey = "jwks_custom_endpoint"
)

const (
	authFlowType                    authFieldKey = "auth_flow_type"
	AUTH_FLOW_USER_TOKEN            AuthFlow     = "user_token"
	AUTH_FLOW_SERVICE_ACCOUNT_TOKEN AuthFlow     = "sa_token"
	AUTH_FLOW_SERVICE_ACCOUNT_KEY   AuthFlow     = "sa_key"
)

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
	err := setAuthFieldInKeyring(key, value)
	if err != nil {
		errFallback := setAuthFieldInEncodedTextFile(key, value)
		if errFallback != nil {
			return fmt.Errorf("write to keyring failed (%w), try writing to encoded text file: %w", err, errFallback)
		}
	}
	return nil
}

func setAuthFieldInKeyring(key authFieldKey, value string) error {
	return keyring.Set(keyringService, string(key), value)
}

func setAuthFieldInEncodedTextFile(key authFieldKey, value string) error {
	err := createEncodedTextFile()
	if err != nil {
		return err
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("get home dir: %w", err)
	}
	textFileDir := filepath.Join(homeDir, textFileFolderName)
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
	value, err := getAuthFieldFromKeyring(key)
	if err != nil {
		var errFallback error
		value, errFallback = getAuthFieldFromEncodedTextFile(key)
		if errFallback != nil {
			return "", fmt.Errorf("write to keyring failed (%w), tried write to encoded text file: %w", err, errFallback)
		}
	}
	return value, nil
}

func getAuthFieldFromKeyring(key authFieldKey) (string, error) {
	return keyring.Get(keyringService, string(key))
}

func getAuthFieldFromEncodedTextFile(key authFieldKey) (string, error) {
	err := createEncodedTextFile()
	if err != nil {
		return "", err
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("get home dir: %w", err)
	}
	textFileDir := filepath.Join(homeDir, textFileFolderName)
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
func createEncodedTextFile() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("get home dir: %w", err)
	}
	textFileDir := filepath.Join(homeDir, textFileFolderName)
	textFilePath := filepath.Join(textFileDir, textFileName)

	err = os.MkdirAll(textFileDir, os.ModePerm)
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
