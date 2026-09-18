package utils_test

import (
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/services/sca/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/sca/v1alphaapi"
)

func TestApplicationStatusToStr(t *testing.T) {
	tests := []struct {
		desc     string
		status   v1alphaapi.CurrentStatus
		expected string
	}{
		{desc: "Running", status: v1alphaapi.CURRENTSTATUS_CURRENT_STATUS_RUNNING, expected: "Running"},
		{desc: "Progressing", status: v1alphaapi.CURRENTSTATUS_CURRENT_STATUS_PROGRESSING, expected: "Progressing"},
		{desc: "Idle", status: v1alphaapi.CURRENTSTATUS_CURRENT_STATUS_IDLE, expected: "Idle"},
		{desc: "Failed", status: v1alphaapi.CURRENTSTATUS_CURRENT_STATUS_FAILED, expected: "Failed"},
		{desc: "None", status: v1alphaapi.CURRENTSTATUS_CURRENT_STATUS_NONE, expected: "None"},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			got := utils.ApplicationStatusToStr(tt.status)
			if tt.expected != got {
				t.Errorf("expected converted status to be %q, got %q", tt.expected, got)
			}
		})
	}
}

func TestAppllicationStateToStr(t *testing.T) {
	tests := []struct {
		desc     string
		stopped  bool
		expected string
	}{
		{desc: "Not stopped", stopped: false, expected: "Active"},
		{desc: "Stopped", stopped: true, expected: "Stopped"},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			got := utils.ApplicationStateToStr(tt.stopped)
			if tt.expected != got {
				t.Errorf("expected state to be %q, got %q", tt.expected, got)
			}
		})
	}
}
