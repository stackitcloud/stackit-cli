package utils

import (
	"reflect"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
)

func TestParseGeofencing(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  map[string][]string
	}{
		{
			name:  "empty input",
			input: nil,
			want:  map[string][]string{},
		},
		{
			name: "single entry",
			input: []string{
				"https://example.com US,CA,MX",
			},
			want: map[string][]string{
				"https://example.com": {"US", "CA", "MX"},
			},
		},
		{
			name: "multiple entries",
			input: []string{
				"https://example.com US,CA,MX",
				"https://another.com DE,FR",
			},
			want: map[string][]string{
				"https://example.com": {"US", "CA", "MX"},
				"https://another.com": {"DE", "FR"},
			},
		},
	}
	printer := print.NewPrinter()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseGeofencing(printer, tt.input)
			if !reflect.DeepEqual(got, &tt.want) {
				t.Errorf("ParseGeofencing() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseOriginRequestHeaders(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  map[string]string
	}{
		{
			name:  "empty input",
			input: nil,
			want:  map[string]string{},
		},
		{
			name: "single entry",
			input: []string{
				"X-Custom-Header: Value1",
			},
			want: map[string]string{
				"X-Custom-Header": "Value1",
			},
		},
		{
			name: "multiple entries",
			input: []string{
				"X-Custom-Header1: Value1",
				"X-Custom-Header2: Value2",
			},
			want: map[string]string{
				"X-Custom-Header1": "Value1",
				"X-Custom-Header2": "Value2",
			},
		},
	}
	printer := print.NewPrinter()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseOriginRequestHeaders(printer, tt.input)
			if !reflect.DeepEqual(got, &tt.want) {
				t.Errorf("ParseOriginRequestHeaders() = %v, want %v", got, tt.want)
			}
		})
	}
}
