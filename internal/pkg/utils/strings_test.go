package utils

import (
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
