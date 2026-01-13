// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 STACKIT GmbH & Co. KG

// Package error provides custom error types for STACKIT Edge Cloud operations.
//
// This package defines structured error types that provide better error handling
// and type checking compared to simple string errors. Each error type can carry
// additional context and implements the standard error interface.
package error

import (
	"fmt"
)

// NoIdentifierError indicates that no identifier was provided when one was required.
type NoIdentifierError struct {
	Operation string // Optional: which operation failed
}

func (e *NoIdentifierError) Error() string {
	if e.Operation != "" {
		return fmt.Sprintf("no identifier provided for %s", e.Operation)
	}
	return "no identifier provided"
}

// InvalidIdentifierError indicates that an unsupported identifier was provided.
type InvalidIdentifierError struct {
	Identifier string // The invalid identifier that was provided
}

func (e *InvalidIdentifierError) Error() string {
	if e.Identifier != "" {
		return fmt.Sprintf("unsupported identifier provided: %s", e.Identifier)
	}
	return "unsupported identifier provided"
}

// InstanceExistsError indicates that a specific instance already exists.
type InstanceExistsError struct {
	DisplayName string // Optional: the display name that was searched for
}

func (e *InstanceExistsError) Error() string {
	if e.DisplayName != "" {
		return fmt.Sprintf("instance already exists: %s", e.DisplayName)
	}
	return "instance already exists"
}

// NoInstanceError indicates that no instance was provided in a context where one was expected.
type NoInstanceError struct {
	Context string // Optional: context where no instance was found (e.g., "in response", "in project")
}

func (e *NoInstanceError) Error() string {
	if e.Context != "" {
		return fmt.Sprintf("no instance provided %s", e.Context)
	}
	return "no instance provided"
}

// NewNoIdentifierError creates a new NoIdentifierError with optional context.
func NewNoIdentifierError(operation string) *NoIdentifierError {
	return &NoIdentifierError{Operation: operation}
}

// NewInvalidIdentifierError creates a new InvalidIdentifierError with the provided identifier.
func NewInvalidIdentifierError(identifier string) *InvalidIdentifierError {
	return &InvalidIdentifierError{
		Identifier: identifier,
	}
}

// NewInstanceExistsError creates a new InstanceExistsError with optional instance details.
func NewInstanceExistsError(displayName string) *InstanceExistsError {
	return &InstanceExistsError{
		DisplayName: displayName,
	}
}

// NewNoInstanceError creates a new NoInstanceError with optional context.
func NewNoInstanceError(context string) *NoInstanceError {
	return &NoInstanceError{Context: context}
}
