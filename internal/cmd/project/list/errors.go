package list

import (
	"errors"
	"net/http"

	oapiError "github.com/stackitcloud/stackit-sdk-go/core/oapierror"
)

func isForbiddenError(err error) bool {
	var oAPIError *oapiError.GenericOpenAPIError
	if ok := errors.As(err, &oAPIError); !ok {
		return false
	}
	if oAPIError.StatusCode != http.StatusForbidden {
		return false
	}
	return true
}
