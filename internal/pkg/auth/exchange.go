package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func ExchangeToken(ctx context.Context, idpClient *http.Client, accessToken, resource string) (string, error) {
	tokenEndpoint, err := GetAuthField(IDP_TOKEN_ENDPOINT)
	if err != nil {
		return "", fmt.Errorf("get idp token endpoint: %w", err)
	}

	req, err := buildRequestToExchangeTokens(ctx, tokenEndpoint, accessToken, resource)
	if err != nil {
		return "", fmt.Errorf("build request: %w", err)
	}
	resp, err := idpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("call API: %w", err)
	}
	defer func() {
		tempErr := resp.Body.Close()
		if tempErr != nil {
			err = fmt.Errorf("close response body: %w", tempErr)
		}
	}()

	clusterToken, err := parseTokenExchangeResponse(resp)
	if err != nil {
		return "", fmt.Errorf("parse API response: %w", err)
	}
	return clusterToken, nil
}

func buildRequestToExchangeTokens(ctx context.Context, tokenEndpoint, accessToken, resource string) (*http.Request, error) {
	idpClientID, err := getIDPClientID()
	if err != nil {
		return nil, err
	}

	form := url.Values{}
	form.Set("grant_type", "urn:ietf:params:oauth:grant-type:token-exchange")
	form.Set("client_id", idpClientID)
	form.Set("subject_token_type", "urn:ietf:params:oauth:token-type:access_token")
	form.Set("requested_token_type", "urn:ietf:params:oauth:token-type:id_token")
	form.Set("scope", "openid profile email groups")
	form.Set("subject_token", accessToken)
	form.Set("resource", resource)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		tokenEndpoint,
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return nil, fmt.Errorf("build exchange request: %w", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return req, nil
}

func parseTokenExchangeResponse(resp *http.Response) (accessToken string, err error) {
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read body: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("non-OK %d status: %s", resp.StatusCode, string(respBody))
	}

	respContent := struct {
		AccessToken string `json:"access_token"`
	}{}
	err = json.Unmarshal(respBody, &respContent)
	if err != nil {
		return "", fmt.Errorf("unmarshal body: %w", err)
	}
	if respContent.AccessToken == "" {
		return "", fmt.Errorf("no access token found")
	}
	return respContent.AccessToken, nil
}
