package utils

import (
	sca "github.com/stackitcloud/stackit-sdk-go/services/sca/v1alphaapi"
)

func ApplicationStatusToStr(status sca.CurrentStatus) string {
	switch status {
	case sca.CURRENTSTATUS_CURRENT_STATUS_RUNNING:
		return "Running"
	case sca.CURRENTSTATUS_CURRENT_STATUS_IDLE:
		return "Idle"
	case sca.CURRENTSTATUS_CURRENT_STATUS_FAILED:
		return "Failed"
	case sca.CURRENTSTATUS_CURRENT_STATUS_PROGRESSING:
		return "Progressing"
	case sca.CURRENTSTATUS_CURRENT_STATUS_NONE:
		return "None"
	default:
		return "None"
	}
}

func ApplicationStateToStr(isStopped bool) string {
	if isStopped == false {
		return "Active"
	}
	return "Stopped"
}
