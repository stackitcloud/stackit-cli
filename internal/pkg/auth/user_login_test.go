package auth

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

type apiClientMocked struct {
	getFails    bool
	getResponse string
}

func (a *apiClientMocked) Do(_ *http.Request) (*http.Response, error) {
	if a.getFails {
		return &http.Response{
			StatusCode: http.StatusNotFound,
		}, fmt.Errorf("not found")
	}
	return &http.Response{
		Status:     "200 OK",
		StatusCode: http.StatusAccepted,
		Body:       io.NopCloser(strings.NewReader(a.getResponse)),
	}, nil
}
