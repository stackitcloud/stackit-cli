package cache

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
)

func TestGetObject(t *testing.T) {
	if err := Init(); err != nil {
		t.Fatalf("cache init failed: %s", err)
	}

	tests := []struct {
		description string
		identifier  string
		expectFile  bool
		expectedErr error
	}{
		{
			description: "identifier exists",
			identifier:  "test-cache-get-exists",
			expectFile:  true,
			expectedErr: nil,
		},
		{
			description: "identifier does not exist",
			identifier:  "test-cache-get-not-exists",
			expectFile:  false,
			expectedErr: os.ErrNotExist,
		},
		{
			description: "identifier is invalid",
			identifier:  "in../../valid",
			expectFile:  false,
			expectedErr: ErrorInvalidCacheIdentifier,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			id := tt.identifier + "-" + uuid.NewString()

			// setup
			if tt.expectFile {
				err := os.MkdirAll(cacheFolderPath, 0o750)
				if err != nil {
					t.Fatalf("create cache folder: %s", err.Error())
				}
				path := filepath.Join(cacheFolderPath, id)
				if err := os.WriteFile(path, []byte("dummy"), 0o600); err != nil {
					t.Fatalf("setup: WriteFile (%s) failed", path)
				}
			}
			// test
			file, err := GetObject(id)

			if !errors.Is(err, tt.expectedErr) {
				t.Fatalf("returned error (%q) does not match %q", err.Error(), tt.expectedErr.Error())
			}

			if tt.expectFile {
				if len(file) < 1 {
					t.Fatalf("expected a file but byte array is empty (len %d)", len(file))
				}
			} else {
				if len(file) > 0 {
					t.Fatalf("didn't expect a file, but byte array is not empty (len %d)", len(file))
				}
			}
		})
	}
}
func TestPutObject(t *testing.T) {
	if err := Init(); err != nil {
		t.Fatalf("cache init failed: %s", err)
	}

	tests := []struct {
		description  string
		identifier   string
		existingFile bool
		expectFile   bool
		expectedErr  error
		customPath   string
	}{
		{
			description:  "identifier already exists",
			identifier:   "test-cache-put-exists",
			existingFile: true,
			expectFile:   true,
			expectedErr:  nil,
		},
		{
			description: "identifier does not exist",
			identifier:  "test-cache-put-not-exists",
			expectFile:  true,
			expectedErr: nil,
		},
		{
			description: "identifier is invalid",
			identifier:  "in../../valid",
			expectFile:  false,
			expectedErr: ErrorInvalidCacheIdentifier,
		},
		{
			description: "directory does not yet exist",
			identifier:  "test-cache-put-folder-not-exists",
			expectFile:  true,
			expectedErr: nil,
			customPath:  "/tmp/stackit-cli-test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			id := tt.identifier + "-" + uuid.NewString()
			if tt.customPath != "" {
				cacheFolderPath = tt.customPath
			} else {
				cacheDir, _ := os.UserCacheDir()
				cacheFolderPath = filepath.Join(cacheDir, "stackit")
			}
			path := filepath.Join(cacheFolderPath, id)

			// setup
			if tt.existingFile {
				if err := os.WriteFile(path, []byte("dummy"), 0o600); err != nil {
					t.Fatalf("setup: WriteFile (%s) failed", path)
				}
			}
			// test
			err := PutObject(id, []byte("dummy"))

			if !errors.Is(err, tt.expectedErr) {
				t.Fatalf("returned error (%q) does not match %q", err.Error(), tt.expectedErr.Error())
			}

			if tt.expectFile {
				if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
					t.Fatalf("expected file (%q) to exist", path)
				}
			}
		})
	}
}

func TestDeleteObject(t *testing.T) {
	if err := Init(); err != nil {
		t.Fatalf("cache init failed: %s", err)
	}

	tests := []struct {
		description  string
		identifier   string
		existingFile bool
		expectedErr  error
	}{
		{
			description:  "identifier exists",
			identifier:   "test-cache-delete-exists",
			existingFile: true,
			expectedErr:  nil,
		},
		{
			description:  "identifier does not exist",
			identifier:   "test-cache-delete-not-exists",
			existingFile: false,
			expectedErr:  nil,
		},
		{
			description:  "identifier is invalid",
			identifier:   "in../../valid",
			existingFile: false,
			expectedErr:  ErrorInvalidCacheIdentifier,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			id := tt.identifier + "-" + uuid.NewString()
			path := filepath.Join(cacheFolderPath, id)

			// setup
			if tt.existingFile {
				if err := os.WriteFile(path, []byte("dummy"), 0o600); err != nil {
					t.Fatalf("setup: WriteFile (%s) failed", path)
				}
			}
			// test
			err := DeleteObject(id)

			if !errors.Is(err, tt.expectedErr) {
				t.Fatalf("returned error (%q) does not match %q", err.Error(), tt.expectedErr.Error())
			}

			if tt.existingFile {
				if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
					t.Fatalf("expected file (%q) to not exist", path)
				}
			}
		})
	}
}
