package utils

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"

	sca "github.com/stackitcloud/stackit-sdk-go/services/sca/v1alphaapi"
)

func GetDefaultPayload() *sca.CreateApplicationPayload {
	return &sca.CreateApplicationPayload{
		Containers: []sca.Container{
			{
				Name:   "container-1",
				Image:  "",
				Cpu:    sca.PtrInt32(1000),
				Memory: sca.PtrInt32(1024),
			},
		},
		Scaling: sca.Scaling{
			Type: sca.SCALINGTYPE_SCALING_TYPE_MANUAL,
			ManualScaling: &sca.ManualScaling{
				Instances: 1,
			},
		},
		Network: sca.Network{
			PublicIngress: true,
			Port:          sca.PtrInt32(8080),
		},
	}
}

func ApplicationStatusToStr(status *sca.CurrentStatus) string {
	if status == nil {
		return "None"
	}
	switch *status {
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
	if !isStopped {
		return "Active"
	}
	return "Stopped"
}

func HumanReadableScalingType(st sca.ScalingType) string {
	switch st {
	case sca.SCALINGTYPE_SCALING_TYPE_MANUAL:
		return "manual"
	case sca.SCALINGTYPE_SCALING_TYPE_AUTO:
		return "auto"
	default:
		return "-"
	}
}

func EnvironmentVariablesFromMap(m map[string]string) []sca.EnvVar {
	var envVars []sca.EnvVar
	if m != nil {
		envVars = make([]sca.EnvVar, 0, len(m))
		for k, v := range m {
			envVars = append(envVars, sca.EnvVar{
				Key:    k,
				Value:  v,
				Origin: sca.ENVVARTYPE_ENV_FROM_SOURCE_TYPE_MANUAL.Ptr(),
			})
		}
	}
	return envVars
}

func GetApplicationName(ctx context.Context, apiClient sca.DefaultAPI, projectId, environmentID, applicationID string) (string, error) {
	resp, err := apiClient.GetApplication(ctx, projectId, environmentID, applicationID).Execute()
	if err != nil {
		return "", fmt.Errorf("get application: %w", err)
	}
	return resp.DisplayName, nil
}

func GetEnvironmentName(ctx context.Context, apiClient sca.DefaultAPI, projectId, environmentID string) (string, error) {
	resp, err := apiClient.GetEnvironment(ctx, projectId, environmentID).Execute()
	if err != nil {
		return "", fmt.Errorf("get environment: %w", err)
	}
	return resp.DisplayName, nil
}

func ValidatePort(port int32, flag string) error {
	if port <= 1024 || port > 65535 {
		return &errors.FlagValidationError{
			Flag:    flag,
			Details: "must be a valid non-privileged port (from 1025 to 65535)",
		}
	}
	return nil
}

func ValidateInstances(instances int32, flag string) error {
	if instances < 0 || instances > 10 {
		return &errors.FlagValidationError{
			Flag:    flag,
			Details: "must be an integer between 0 and 10",
		}
	}
	return nil
}

func ValidateCPU(cpu int32, flag string) error {
	if cpu <= 0 || cpu%1000 != 0 {
		return &errors.FlagValidationError{
			Flag:    flag,
			Details: "must be a valid value divisible by 1000",
		}
	}
	return nil
}
