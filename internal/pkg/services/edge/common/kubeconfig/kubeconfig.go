// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 STACKIT GmbH & Co. KG

package kubeconfig

import (
	"fmt"
	"maps"
	"math"
	"os"
	"path/filepath"

	"k8s.io/client-go/tools/clientcmd"
)

// Validation constants taken from OpenApi spec.
const (
	expirationSecondsMax = 15552000 // 60 * 60 * 24 * 180 seconds = 180 days
	expirationSecondsMin = 600      // 60 * 10 seconds = 10 minutes
)

// Defaults taken from OpenApi spec.
const (
	ExpirationSecondsDefault = 3600 // 60 * 60 seconds = 1 hour
)

// User input flags for kubeconfig commands
const (
	ExpirationFlag     = "expiration"
	DisableWritingFlag = "disable-writing"
	FilepathFlag       = "filepath"
	OverwriteFlag      = "overwrite"
	SwitchContextFlag  = "switch-context"
)

// Flag usage texts
const (
	ExpirationUsage     = "Expiration time for the kubeconfig, e.g. 5d. By default, the token is valid for 1h."
	FilepathUsage       = "Path to the kubeconfig file. A default is chosen by Kubernetes if not set."
	DisableWritingUsage = "Disable writing the kubeconfig to a file."
	OverwriteUsage      = "Force overwrite the kubeconfig file if it exists."
	SwitchContextUsage  = "Switch to the context in the kubeconfig file to the new context."
)

// Flag shorthands
const (
	ExpirationShorthand     = "e"
	DisableWritingShorthand = ""
	FilepathShorthand       = "f"
	OverwriteShorthand      = ""
	SwitchContextShorthand  = ""
)

func ValidateExpiration(expiration *uint64) error {
	if expiration != nil {
		// We're using utils.ConvertToSeconds to convert the user input string to seconds, which is using
		// math.MaxUint64 internally, if no special limits are set. However: the OpenApi v3 Spec
		// only allows integers (int64). So we could end up in a overflow IF expirationSecondsMax
		// ever is changed beyond the maximum value of int64. This check makes sure this won't happen.
		maxExpiration := uint64(math.Min(expirationSecondsMax, math.MaxInt64))
		if *expiration > maxExpiration {
			return fmt.Errorf("%s is too large (maximum is %d seconds)", ExpirationFlag, maxExpiration)
		}
		// If expiration is ever changed to int64 this check makes sure we never end up with negative expiration times.
		minExpiration := uint64(math.Max(expirationSecondsMin, 0))
		if *expiration < minExpiration {
			return fmt.Errorf("%s is too small (minimum is %d seconds)", ExpirationFlag, minExpiration)
		}
	}
	return nil
}

// EmptyKubeconfigError is returned when the kubeconfig content is empty.
type EmptyKubeconfigError struct{}

// Error returns the error message.
func (e *EmptyKubeconfigError) Error() string {
	return "no data for kubeconfig"
}

// LoadKubeconfigError is returned when loading the kubeconfig fails.
type LoadKubeconfigError struct {
	Err error
}

// Error returns the error message.
func (e *LoadKubeconfigError) Error() string {
	return fmt.Sprintf("load kubeconfig: %v", e.Err)
}

// Unwrap returns the underlying error.
func (e *LoadKubeconfigError) Unwrap() error {
	return e.Err
}

// WriteKubeconfigError is returned when writing the kubeconfig fails.
type WriteKubeconfigError struct {
	Err error
}

// Error returns the error message.
func (e *WriteKubeconfigError) Error() string {
	return fmt.Sprintf("write kubeconfig: %v", e.Err)
}

// Unwrap returns the underlying error.
func (e *WriteKubeconfigError) Unwrap() error {
	return e.Err
}

// InvalidKubeconfigPathError is returned when an invalid kubeconfig path is provided.
type InvalidKubeconfigPathError struct {
	Path string
}

// Error returns the error message.
func (e *InvalidKubeconfigPathError) Error() string {
	return fmt.Sprintf("invalid path: %s", e.Path)
}

// mergeKubeconfig merges new kubeconfig data into a kubeconfig file.
//
// If the destination file does not exist, it will be created. If the file exists,
// the new data (clusters, contexts, and users) is merged into the existing
// configuration, overwriting entries with the same name and replacing the
// current-context if defined in the new data.
//
// The function takes the following parameters:
//   - configPath: The path to the destination file. The file and the directory tree
//     for the file will be created if it does not exist.
//   - data: The new kubeconfig content to merge. Merge is performed based on standard
//     kubeconfig structure.
//   - switchContext: If true, the function will switch to the new context in the
//     kubeconfig file after merging.
//
// It returns a nil error on success. On failure, it returns an error indicating
// if the provided data was empty, malformed, or if there were issues reading from
// or writing to the filesystem.
func mergeKubeconfig(filePath *string, data string, switchContext bool) error {
	if filePath == nil {
		return fmt.Errorf("no kubeconfig file provided to be merged")
	}
	path := *filePath

	// Check if the new kubeconfig data is empty
	if data == "" {
		return &EmptyKubeconfigError{}
	}

	// Load and validate the data into a kubeconfig object
	newConfig, err := clientcmd.Load([]byte(data))
	if err != nil {
		return &LoadKubeconfigError{Err: err}
	}

	// If the destination kubeconfig does not exist, create a new one. IsNotExist will ignore other errors.
	// Other errors are handled separately by the following clientcmd.LoadFromFile clientcmd.LoadFromFile
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return writeKubeconfig(&path, data)
	}

	// If the file exists load and validate the existing kubeconfig into a config object
	existingConfig, err := clientcmd.LoadFromFile(path)
	if err != nil {
		return &LoadKubeconfigError{Err: err}
	}

	// Merge the new kubeconfig data into the existing config object
	maps.Copy(existingConfig.AuthInfos, newConfig.AuthInfos)
	maps.Copy(existingConfig.Clusters, newConfig.Clusters)
	maps.Copy(existingConfig.Contexts, newConfig.Contexts)

	// If no CurrentContext is set or switchContext is true, set the CurrentContext to the CurrentContext of the new kubeconfig
	if newConfig.CurrentContext != "" && (switchContext || existingConfig.CurrentContext == "") {
		existingConfig.CurrentContext = newConfig.CurrentContext
	}

	// Save the merged config to the file, creating missing directories as needed.
	if err := clientcmd.WriteToFile(*existingConfig, path); err != nil {
		return &WriteKubeconfigError{Err: err}
	}

	return nil
}

// writeKubeconfig writes kubeconfig data to a file, overwriting it if it exists.
//
// The function takes the following parameters:
//   - configPath: The path to the destination file. The file and the directory tree
//     for the file will be created if it does not exist.
//   - data: The new kubeconfig content to write to the file.
//
// It returns a nil error on success. On failure, it returns an error indicating
// if the provided data was empty, malformed, or if there were issues reading from
// or writing to the filesystem.
func writeKubeconfig(filePath *string, data string) error {
	if filePath == nil {
		return fmt.Errorf("no kubeconfig file provided to be written")
	}
	path := *filePath

	// Check if the new kubeconfig data is empty
	if data == "" {
		return &EmptyKubeconfigError{}
	}

	// Load and validate the data into a kubeconfig object
	config, err := clientcmd.Load([]byte(data))
	if err != nil {
		return &LoadKubeconfigError{Err: err}
	}

	// Save the merged config to the file, creating missing directories as needed.
	if err := clientcmd.WriteToFile(*config, path); err != nil {
		return &WriteKubeconfigError{Err: err}
	}

	return nil
}

// getDefaultKubeconfigPath returns the default location for the kubeconfig file,
// following standard Kubernetes loading rules.
//
// It returns a string containing the absolute path to the default kubeconfig file.
func getDefaultKubeconfigPath() string {
	return clientcmd.NewDefaultClientConfigLoadingRules().GetDefaultFilename()
}

// Returns the absolute path to the kubeconfig file.
// If a file path is provided, it is validated and, if valid, returned as an absolute path.
// If nil is provided the default kubeconfig path is loaded and returned as an absolute path.
func getKubeconfigPath(filePath *string) (string, error) {
	if filePath == nil {
		return getDefaultKubeconfigPath(), nil
	}

	if isValidFilePath(filePath) {
		return filepath.Abs(*filePath)
	}
	return "", &InvalidKubeconfigPathError{Path: *filePath}
}

// Basic filesystem path validation. Returns true if the provided string is a path. Returns false otherwise.
func isValidFilePath(filePath *string) bool {
	if filePath == nil || *filePath == "" {
		return false
	}

	// Clean the path and check if it's valid
	cleaned := filepath.Clean(*filePath)
	if cleaned == "." || cleaned == string(filepath.Separator) {
		return false
	}

	// Try to get absolute path (this will fail for invalid paths)
	_, err := filepath.Abs(*filePath)
	// If no error, the path is valid (return true). Otherwise, it's invalid (return false).
	return err == nil
}

// Basic filesystem file existence check. Returns true if the file exists. Returns false otherwise.
func isExistingFile(filePath *string) bool {
	// Check if the kubeconfig file exists
	_, errStat := os.Stat(*filePath)
	return !os.IsNotExist(errStat)
}

// ConfirmationCallback is a function that prompts for confirmation with the given message
// and returns true if confirmed, false otherwise
type ConfirmationCallback func(message string) error

// WriteOptions contains options for writing kubeconfig files
type WriteOptions struct {
	Overwrite     bool
	SwitchContext bool
	ConfirmFn     ConfirmationCallback
}

// WithOverwrite sets whether to overwrite existing files instead of merging
func (w WriteOptions) WithOverwrite(overwrite bool) WriteOptions {
	w.Overwrite = overwrite
	return w
}

// WithSwitchContext sets whether to switch to the new context after writing
func (w WriteOptions) WithSwitchContext(switchContext bool) WriteOptions {
	w.SwitchContext = switchContext
	return w
}

// WithConfirmation sets the confirmation callback function
func (w WriteOptions) WithConfirmation(fn ConfirmationCallback) WriteOptions {
	w.ConfirmFn = fn
	return w
}

// NewWriteOptions creates a new WriteOptions with default values
func NewWriteOptions() WriteOptions {
	return WriteOptions{
		Overwrite:     false,
		SwitchContext: false,
		ConfirmFn:     nil,
	}
}

// WriteKubeconfig writes the provided kubeconfig data to a file on the filesystem.
// By default, if the file already exists, it will be merged with the provided data.
// This behavior can be controlled using the provided options.
//
// The function takes the following parameters:
//   - filePath: The path to the destination file. The file and the directory tree for the
//     file will be created if it does not exist. If nil, the default kubeconfig path is used.
//   - kubeconfig: The kubeconfig content to write.
//   - options: Options for controlling the write behavior.
//
// It returns the file path actually used to write to on success.
func WriteKubeconfig(filePath *string, kubeconfig string, options WriteOptions) (*string, error) {
	// Check if the provided filePath is valid or use the default kubeconfig path no filePath is provided
	path, err := getKubeconfigPath(filePath)
	if err != nil {
		return nil, err
	}

	if isExistingFile(&path) {
		// If the file exists
		if !options.Overwrite {
			// If overwrite was not requested the default it to merge
			if options.ConfirmFn != nil {
				// If confirmation callback is provided, prompt the user for confirmation
				prompt := fmt.Sprintf("Update your kubeconfig %q?", path)
				err := options.ConfirmFn(prompt)
				if err != nil {
					// If the user doesn't confirm do not proceed with the merge
					return nil, err
				}
			}
			err := mergeKubeconfig(&path, kubeconfig, options.SwitchContext)
			if err != nil {
				return nil, err
			}
			return &path, err
		}
		// If overwrite was requested overwrite the existing file
		if options.ConfirmFn != nil {
			// If confirmation callback is provided, prompt the user for confirmation
			prompt := fmt.Sprintf("Replace your kubeconfig %q?", path)
			err := options.ConfirmFn(prompt)
			if err != nil {
				// If the user doesn't confirm do not proceed with the overwrite
				return nil, err
			}
			// Fallthrough
		}
	}
	// If the file doesn't exist or in case the user confirmed the overwrite (fallthrough) write the file
	err = writeKubeconfig(&path, kubeconfig)
	if err != nil {
		return nil, err
	}
	return &path, err
}
