// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 STACKIT GmbH & Co. KG

// Unit tests for package error
package error

import (
	"testing"

	testUtils "github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
)

func TestNoIdentifierError(t *testing.T) {
	type args struct {
		operation string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty",
			args: args{
				operation: "",
			},
			want: "no identifier provided",
		},
		{
			name: "with operation",
			args: args{
				operation: "create",
			},
			want: "no identifier provided for create",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := (&NoIdentifierError{Operation: tt.args.operation}).Error()
			testUtils.AssertValue(t, got, tt.want)
		})
	}
}

func TestInvalidIdentifierError(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty",
			args: args{
				id: "",
			},
			want: "unsupported identifier provided",
		},
		{
			name: "with identifier",
			args: args{
				id: "x-123",
			},
			want: "unsupported identifier provided: x-123",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := (&InvalidIdentifierError{Identifier: tt.args.id}).Error()
			testUtils.AssertValue(t, got, tt.want)
		})
	}
}

func TestInstanceExistsError(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty",
			args: args{name: ""},
			want: "instance already exists"},
		{
			name: "with display name",
			args: args{name: "my-inst"},
			want: "instance already exists: my-inst",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := (&InstanceExistsError{DisplayName: tt.args.name}).Error()
			testUtils.AssertValue(t, got, tt.want)
		})
	}
}

func TestNoInstanceError(t *testing.T) {
	type args struct {
		ctx string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty",
			args: args{
				ctx: "",
			},
			want: "no instance provided",
		},
		{
			name: "with context",
			args: args{
				ctx: "in project",
			},
			want: "no instance provided in project",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := (&NoInstanceError{Context: tt.args.ctx}).Error()
			testUtils.AssertValue(t, got, tt.want)
		})
	}
}

func TestConstructorsReturnExpected(t *testing.T) {
	tests := []struct {
		name string
		got  any
		want any
	}{
		{
			name: "NoIdentifier operation",
			got:  NewNoIdentifierError("op").Operation,
			want: "op",
		},
		{
			name: "InvalidIdentifier identifier",
			got:  NewInvalidIdentifierError("id").Identifier,
			want: "id",
		},
		{
			name: "InstanceExists displayName",
			got:  NewInstanceExistsError("name").DisplayName,
			want: "name",
		},
		{
			name: "NoInstance context",
			got:  NewNoInstanceError("ctx").Context,
			want: "ctx",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantErr, wantIsErr := tt.want.(error)
			gotErr, gotIsErr := tt.got.(error)
			if wantIsErr {
				if !gotIsErr {
					t.Fatalf("expected error but got %T", tt.got)
				}
				testUtils.AssertError(t, gotErr, wantErr)
				return
			}

			testUtils.AssertValue(t, tt.got, tt.want)
		})
	}
}
