package utils

import (
	"reflect"
	"testing"

	"github.com/stackitcloud/stackit-sdk-go/core/utils"
)

func TestTruncate(t *testing.T) {
	type args struct {
		s      *string
		maxLen int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"nil string", args{nil, 3}, ""},
		{"empty string", args{utils.Ptr(""), 10}, ""},
		{"length below maxlength", args{utils.Ptr("foo"), 10}, "foo"},
		{"exactly maxlength", args{utils.Ptr("foo"), 3}, "foo"},
		{"above maxlength", args{utils.Ptr("foobarbaz"), 3}, "foo…"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Truncate(tt.args.s, tt.args.maxLen); got != tt.want {
				t.Errorf("Truncate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSortedKeys(t *testing.T) {
	tests := []struct {
		name  string
		input map[string]int
		want  []string
	}{
		{
			name:  "nil map",
			input: nil,
			want:  []string{},
		},
		{
			name:  "empty map",
			input: map[string]int{},
			want:  []string{},
		},
		{
			name:  "single element",
			input: map[string]int{"b": 1},
			want:  []string{"b"},
		},
		{
			name:  "multiple elements unsorted",
			input: map[string]int{"c": 3, "a": 1, "b": 2},
			want:  []string{"a", "b", "c"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SortedKeys(tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SortedKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJoinStringMap(t *testing.T) {
	tests := []struct {
		name  string
		input map[string]string
		want  string
	}{
		{
			name:  "nil map",
			input: nil,
			want:  "",
		},
		{
			name:  "empty map",
			input: map[string]string{},
			want:  "",
		},
		{
			name:  "single element",
			input: map[string]string{"key1": "value1"},
			want:  "key1=value1",
		},
		{
			name:  "multiple elements",
			input: map[string]string{"key1": "value1", "key2": "value2"},
			want:  "key1=value1, key2=value2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := JoinStringMap(tt.input, "=", ", "); got != tt.want {
				t.Errorf("JoinStringMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
