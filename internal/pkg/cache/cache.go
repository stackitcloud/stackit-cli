package cache

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/stackitcloud/stackit-cli/internal/pkg/auth"
)

var (
	cacheFolderPath    string
	cacheEncryptionKey []byte

	identifierRegex             = regexp.MustCompile("^[a-zA-Z0-9-]+$")
	ErrorInvalidCacheIdentifier = fmt.Errorf("invalid cache identifier")
)

const (
	cacheKeyMaxAge = 90 * 24 * time.Hour
)

func Init() error {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return fmt.Errorf("get user cache dir: %w", err)
	}
	cacheFolderPath = filepath.Join(cacheDir, "stackit")

	// Encryption keys should only be used a limited number of times for aes-gcm.
	// Thus, refresh the key periodically. This will invalidate all cached entries.
	key, _ := auth.GetAuthField(auth.CACHE_ENCRYPTION_KEY)
	age, _ := auth.GetAuthField(auth.CACHE_ENCRYPTION_KEY_AGE)
	cacheEncryptionKey = nil
	var keyAge time.Time
	if age != "" {
		ageSeconds, err := strconv.ParseInt(age, 10, 64)
		if err == nil {
			keyAge = time.Unix(ageSeconds, 0)
		}
	}
	if key != "" && keyAge.Add(cacheKeyMaxAge).After(time.Now()) {
		cacheEncryptionKey, _ = base64.StdEncoding.DecodeString(key)
		// invalid key length
		if len(cacheEncryptionKey) != 32 {
			cacheEncryptionKey = nil
		}
	}
	if len(cacheEncryptionKey) == 0 {
		cacheEncryptionKey = make([]byte, 32)
		_, err := rand.Read(cacheEncryptionKey)
		if err != nil {
			return fmt.Errorf("cache encryption key: %v", err)
		}
		key := base64.StdEncoding.EncodeToString(cacheEncryptionKey)
		err = auth.SetAuthField(auth.CACHE_ENCRYPTION_KEY, key)
		if err != nil {
			return fmt.Errorf("save cache encryption key: %v", err)
		}
		err = auth.SetAuthField(auth.CACHE_ENCRYPTION_KEY_AGE, fmt.Sprint(time.Now().Unix()))
		if err != nil {
			return fmt.Errorf("save cache encryption key age: %v", err)
		}
	}
	return nil
}

func GetObject(identifier string) ([]byte, error) {
	if err := validateCacheFolderPath(); err != nil {
		return nil, err
	}
	if !identifierRegex.MatchString(identifier) {
		return nil, ErrorInvalidCacheIdentifier
	}

	data, err := os.ReadFile(filepath.Join(cacheFolderPath, identifier))
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(cacheEncryptionKey)
	if err != nil {
		return nil, err
	}
	aead, err := cipher.NewGCMWithRandomNonce(block)
	if err != nil {
		return nil, err
	}

	return aead.Open(nil, nil, data, nil)
}

func PutObject(identifier string, data []byte) error {
	if err := validateCacheFolderPath(); err != nil {
		return err
	}
	if !identifierRegex.MatchString(identifier) {
		return ErrorInvalidCacheIdentifier
	}

	err := os.MkdirAll(cacheFolderPath, 0o750)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(cacheEncryptionKey)
	if err != nil {
		return err
	}
	aead, err := cipher.NewGCMWithRandomNonce(block)
	if err != nil {
		return err
	}
	encrypted := aead.Seal(nil, nil, data, nil)

	return os.WriteFile(filepath.Join(cacheFolderPath, identifier), encrypted, 0o600)
}

func DeleteObject(identifier string) error {
	if err := validateCacheFolderPath(); err != nil {
		return err
	}
	if !identifierRegex.MatchString(identifier) {
		return ErrorInvalidCacheIdentifier
	}

	if err := os.Remove(filepath.Join(cacheFolderPath, identifier)); !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}

func validateCacheFolderPath() error {
	if cacheFolderPath == "" {
		return errors.New("cacheFolderPath not set. Forgot to call Init()?")
	}
	return nil
}
