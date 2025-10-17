package wait

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-sdk-go/core/oapierror"
	"github.com/stackitcloud/stackit-sdk-go/core/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/intake"
)

// apiClientMocked is a mock of the API client
type apiClientMocked struct {
	getRunnerFails       bool
	getIntakeFails       bool
	getUserFails         bool
	getErrorCode         int
	returnRunner         bool
	returnIntake         bool
	returnUser           bool
	intakeRunnerResponse *intake.IntakeRunnerResponse
	intakeResponse       *intake.IntakeResponse
	intakeUserResponse   *intake.IntakeUserResponse
}

func (a *apiClientMocked) GetIntakeRunnerExecute(_ context.Context, _, _, _ string) (*intake.IntakeRunnerResponse, error) {
	if a.getRunnerFails {
		return nil, &oapierror.GenericOpenAPIError{
			StatusCode: a.getErrorCode,
		}
	}
	if !a.returnRunner {
		return nil, nil
	}
	return a.intakeRunnerResponse, nil
}

func (a *apiClientMocked) GetIntakeExecute(_ context.Context, _, _, _ string) (*intake.IntakeResponse, error) {
	if a.getIntakeFails {
		return nil, &oapierror.GenericOpenAPIError{
			StatusCode: a.getErrorCode,
		}
	}
	if !a.returnIntake {
		return nil, nil
	}
	return a.intakeResponse, nil
}

func (a *apiClientMocked) GetIntakeUserExecute(_ context.Context, _, _, _, _ string) (*intake.IntakeUserResponse, error) {
	if a.getUserFails {
		return nil, &oapierror.GenericOpenAPIError{
			StatusCode: a.getErrorCode,
		}
	}
	if !a.returnUser {
		return nil, nil
	}
	return a.intakeUserResponse, nil
}

var (
	PROJECT_ID       = uuid.NewString()
	REGION           = "eu01"
	INTAKE_RUNNER_ID = uuid.NewString()
	INTAKE_ID        = uuid.NewString()
	INTAKE_USER_ID   = uuid.NewString()
)

func TestCreateOrUpdateIntakeRunnerWaitHandler(t *testing.T) {
	tests := []struct {
		desc                 string
		getFails             bool
		getErrorCode         int
		wantErr              bool
		wantResp             bool
		returnRunner         bool
		intakeRunnerResponse *intake.IntakeRunnerResponse
	}{
		{
			desc:         "succeeded",
			getFails:     false,
			wantErr:      false,
			wantResp:     true,
			returnRunner: true,
			intakeRunnerResponse: &intake.IntakeRunnerResponse{
				Id:    utils.Ptr(INTAKE_RUNNER_ID),
				State: utils.Ptr(intake.INTAKERUNNERRESPONSESTATE_ACTIVE),
			},
		},
		{
			desc:         "get fails",
			getFails:     true,
			getErrorCode: http.StatusInternalServerError,
			wantErr:      true,
			wantResp:     false,
			returnRunner: false,
		},
		{
			desc:         "timeout",
			getFails:     false,
			wantErr:      true,
			wantResp:     false,
			returnRunner: true,
			intakeRunnerResponse: &intake.IntakeRunnerResponse{
				Id:    utils.Ptr(INTAKE_RUNNER_ID),
				State: utils.Ptr(intake.INTAKERUNNERRESPONSESTATE_RECONCILING),
			},
		},
		{
			desc:         "nil response",
			getFails:     false,
			wantErr:      true,
			wantResp:     false,
			returnRunner: false,
		},
		{
			desc:         "nil id in response",
			getFails:     false,
			wantErr:      true,
			wantResp:     false,
			returnRunner: true,
			intakeRunnerResponse: &intake.IntakeRunnerResponse{
				State: utils.Ptr(intake.INTAKERUNNERRESPONSESTATE_RECONCILING),
			},
		},
		{
			desc:         "nil state in response",
			getFails:     false,
			wantErr:      true,
			wantResp:     false,
			returnRunner: true,
			intakeRunnerResponse: &intake.IntakeRunnerResponse{
				Id: utils.Ptr(INTAKE_RUNNER_ID),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			apiClient := &apiClientMocked{
				getRunnerFails:       tt.getFails,
				getErrorCode:         tt.getErrorCode,
				returnRunner:         tt.returnRunner,
				intakeRunnerResponse: tt.intakeRunnerResponse,
			}

			var wantResp *intake.IntakeRunnerResponse
			if tt.wantResp {
				wantResp = tt.intakeRunnerResponse
			}

			handler := CreateOrUpdateIntakeRunnerWaitHandler(context.Background(), apiClient, PROJECT_ID, REGION, INTAKE_RUNNER_ID)
			got, err := handler.SetTimeout(10 * time.Millisecond).WaitWithContext(context.Background())

			if (err != nil) != tt.wantErr {
				t.Fatalf("handler error = %v, wantErr %v", err, tt.wantErr)
			}
			if !cmp.Equal(got, wantResp) {
				t.Fatalf("handler got = %v, want %v", got, wantResp)
			}
		})
	}
}

func TestDeleteIntakeRunnerWaitHandler(t *testing.T) {
	tests := []struct {
		desc         string
		getFails     bool
		getErrorCode int
		wantErr      bool
		returnRunner bool
	}{
		{
			desc:         "succeeded",
			getFails:     true,
			getErrorCode: http.StatusNotFound,
			wantErr:      false,
			returnRunner: false,
		},
		{
			desc:         "get fails",
			getFails:     true,
			getErrorCode: http.StatusInternalServerError,
			wantErr:      true,
			returnRunner: false,
		},
		{
			desc:         "timeout",
			getFails:     false,
			wantErr:      true,
			returnRunner: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			apiClient := &apiClientMocked{
				getRunnerFails: tt.getFails,
				getErrorCode:   tt.getErrorCode,
				returnRunner:   tt.returnRunner,
				intakeRunnerResponse: &intake.IntakeRunnerResponse{ // This is only used in the timeout case
					Id: utils.Ptr(INTAKE_RUNNER_ID),
				},
			}
			handler := DeleteIntakeRunnerWaitHandler(context.Background(), apiClient, PROJECT_ID, REGION, INTAKE_RUNNER_ID)
			_, err := handler.SetTimeout(10 * time.Millisecond).WaitWithContext(context.Background())

			if (err != nil) != tt.wantErr {
				t.Fatalf("handler error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateOrUpdateIntakeWaitHandler(t *testing.T) {
	tests := []struct {
		desc           string
		getFails       bool
		getErrorCode   int
		wantErr        bool
		wantResp       bool
		returnIntake   bool
		intakeResponse *intake.IntakeResponse
	}{
		{
			desc:         "succeeded",
			getFails:     false,
			wantErr:      false,
			wantResp:     true,
			returnIntake: true,
			intakeResponse: &intake.IntakeResponse{
				Id:    utils.Ptr(INTAKE_ID),
				State: utils.Ptr(intake.INTAKERESPONSESTATE_ACTIVE),
			},
		},
		{
			desc:         "failed state",
			getFails:     false,
			wantErr:      true,
			wantResp:     true,
			returnIntake: true,
			intakeResponse: &intake.IntakeResponse{
				Id:    utils.Ptr(INTAKE_ID),
				State: utils.Ptr(intake.INTAKERESPONSESTATE_FAILED),
			},
		},
		{
			desc:         "get fails",
			getFails:     true,
			getErrorCode: http.StatusInternalServerError,
			wantErr:      true,
			wantResp:     false,
			returnIntake: false,
		},
		{
			desc:         "timeout",
			getFails:     false,
			wantErr:      true,
			wantResp:     false,
			returnIntake: true,
			intakeResponse: &intake.IntakeResponse{
				Id:    utils.Ptr(INTAKE_ID),
				State: utils.Ptr(intake.INTAKERESPONSESTATE_RECONCILING),
			},
		},
		{
			desc:         "nil response",
			getFails:     false,
			wantErr:      true,
			wantResp:     false,
			returnIntake: false,
		},
		{
			desc:         "nil id in response",
			getFails:     false,
			wantErr:      true,
			wantResp:     false,
			returnIntake: true,
			intakeResponse: &intake.IntakeResponse{
				State: utils.Ptr(intake.INTAKERESPONSESTATE_RECONCILING),
			},
		},
		{
			desc:         "nil state in response",
			getFails:     false,
			wantErr:      true,
			wantResp:     false,
			returnIntake: true,
			intakeResponse: &intake.IntakeResponse{
				Id: utils.Ptr(INTAKE_ID),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			apiClient := &apiClientMocked{
				getIntakeFails: tt.getFails,
				getErrorCode:   tt.getErrorCode,
				returnIntake:   tt.returnIntake,
				intakeResponse: tt.intakeResponse,
			}

			var wantResp *intake.IntakeResponse
			if tt.wantResp {
				wantResp = tt.intakeResponse
			}

			handler := CreateOrUpdateIntakeWaitHandler(context.Background(), apiClient, PROJECT_ID, REGION, INTAKE_ID)
			got, err := handler.SetTimeout(10 * time.Millisecond).WaitWithContext(context.Background())

			if (err != nil) != tt.wantErr {
				t.Fatalf("handler error = %v, wantErr %v", err, tt.wantErr)
			}
			if !cmp.Equal(got, wantResp) {
				t.Fatalf("handler got = %v, want %v", got, wantResp)
			}
		})
	}
}

func TestDeleteIntakeWaitHandler(t *testing.T) {
	tests := []struct {
		desc         string
		getFails     bool
		getErrorCode int
		wantErr      bool
		returnIntake bool
	}{
		{
			desc:         "succeeded",
			getFails:     true,
			getErrorCode: http.StatusNotFound,
			wantErr:      false,
			returnIntake: false,
		},
		{
			desc:         "get fails",
			getFails:     true,
			getErrorCode: http.StatusInternalServerError,
			wantErr:      true,
			returnIntake: false,
		},
		{
			desc:         "timeout",
			getFails:     false,
			wantErr:      true,
			returnIntake: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			apiClient := &apiClientMocked{
				getIntakeFails: tt.getFails,
				getErrorCode:   tt.getErrorCode,
				returnIntake:   tt.returnIntake,
				intakeResponse: &intake.IntakeResponse{
					Id: utils.Ptr(INTAKE_ID),
				},
			}
			handler := DeleteIntakeWaitHandler(context.Background(), apiClient, PROJECT_ID, REGION, INTAKE_ID)
			_, err := handler.SetTimeout(10 * time.Millisecond).WaitWithContext(context.Background())

			if (err != nil) != tt.wantErr {
				t.Fatalf("handler error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateOrUpdateIntakeUserWaitHandler(t *testing.T) {
	tests := []struct {
		desc               string
		getFails           bool
		getErrorCode       int
		wantErr            bool
		wantResp           bool
		returnUser         bool
		intakeUserResponse *intake.IntakeUserResponse
	}{
		{
			desc:       "succeeded",
			getFails:   false,
			wantErr:    false,
			wantResp:   true,
			returnUser: true,
			intakeUserResponse: &intake.IntakeUserResponse{
				Id:    utils.Ptr(INTAKE_USER_ID),
				State: utils.Ptr(intake.INTAKEUSERRESPONSESTATE_ACTIVE),
			},
		},
		{
			desc:         "get fails",
			getFails:     true,
			getErrorCode: http.StatusInternalServerError,
			wantErr:      true,
			wantResp:     false,
			returnUser:   false,
		},
		{
			desc:       "timeout",
			getFails:   false,
			wantErr:    true,
			wantResp:   false,
			returnUser: true,
			intakeUserResponse: &intake.IntakeUserResponse{
				Id:    utils.Ptr(INTAKE_USER_ID),
				State: utils.Ptr(intake.INTAKEUSERRESPONSESTATE_RECONCILING),
			},
		},
		{
			desc:       "nil response",
			getFails:   false,
			wantErr:    true,
			wantResp:   false,
			returnUser: false,
		},
		{
			desc:       "nil id in response",
			getFails:   false,
			wantErr:    true,
			wantResp:   false,
			returnUser: true,
			intakeUserResponse: &intake.IntakeUserResponse{
				State: utils.Ptr(intake.INTAKEUSERRESPONSESTATE_RECONCILING),
			},
		},
		{
			desc:       "nil state in response",
			getFails:   false,
			wantErr:    true,
			wantResp:   false,
			returnUser: true,
			intakeUserResponse: &intake.IntakeUserResponse{
				Id: utils.Ptr(INTAKE_USER_ID),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			apiClient := &apiClientMocked{
				getUserFails:       tt.getFails,
				getErrorCode:       tt.getErrorCode,
				returnUser:         tt.returnUser,
				intakeUserResponse: tt.intakeUserResponse,
			}

			var wantResp *intake.IntakeUserResponse
			if tt.wantResp {
				wantResp = tt.intakeUserResponse
			}

			handler := CreateOrUpdateIntakeUserWaitHandler(context.Background(), apiClient, PROJECT_ID, REGION, INTAKE_ID, INTAKE_USER_ID)
			got, err := handler.SetTimeout(10 * time.Millisecond).WaitWithContext(context.Background())

			if (err != nil) != tt.wantErr {
				t.Fatalf("handler error = %v, wantErr %v", err, tt.wantErr)
			}
			if !cmp.Equal(got, wantResp) {
				t.Fatalf("handler got = %v, want %v", got, wantResp)
			}
		})
	}
}

func TestDeleteIntakeUserWaitHandler(t *testing.T) {
	tests := []struct {
		desc         string
		getFails     bool
		getErrorCode int
		wantErr      bool
		returnUser   bool
	}{
		{
			desc:         "succeeded",
			getFails:     true,
			getErrorCode: http.StatusNotFound,
			wantErr:      false,
			returnUser:   false,
		},
		{
			desc:         "get fails",
			getFails:     true,
			getErrorCode: http.StatusInternalServerError,
			wantErr:      true,
			returnUser:   false,
		},
		{
			desc:       "timeout",
			getFails:   false,
			wantErr:    true,
			returnUser: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			apiClient := &apiClientMocked{
				getUserFails: tt.getFails,
				getErrorCode: tt.getErrorCode,
				returnUser:   tt.returnUser,
				intakeUserResponse: &intake.IntakeUserResponse{
					Id: utils.Ptr(INTAKE_USER_ID),
				},
			}
			handler := DeleteIntakeUserWaitHandler(context.Background(), apiClient, PROJECT_ID, REGION, INTAKE_ID, INTAKE_USER_ID)
			_, err := handler.SetTimeout(10 * time.Millisecond).WaitWithContext(context.Background())

			if (err != nil) != tt.wantErr {
				t.Fatalf("handler error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
