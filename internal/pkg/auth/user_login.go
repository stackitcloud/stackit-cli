package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"golang.org/x/oauth2"
)

const (
	authDomain  = "auth.01.idp.eu01.stackit.cloud/oauth"
	clientId    = "stackit-cli-client-id"
	redirectURL = "http://localhost:8000"
)

// AuthorizeUser implements the PKCE OAuth2 flow.
func AuthorizeUser() error {
	conf := &oauth2.Config{
		ClientID: clientId,
		Endpoint: oauth2.Endpoint{
			AuthURL: fmt.Sprintf("https://%s/authorize", authDomain),
		},
		Scopes:      []string{"openid"},
		RedirectURL: redirectURL,
	}

	// Initialize the code verifier
	codeVerifier := oauth2.GenerateVerifier()

	// Construct the authorization URL
	authorizationURL := conf.AuthCodeURL("", oauth2.S256ChallengeOption(codeVerifier))

	// Start a web server to listen on a callback URL
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:              redirectURL,
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}

	// Define a handler that will get the authorization code, call the token endpoint, and close the HTTP server
	var errServer error
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer cleanup(server)

		// Get the authorization code
		code := r.URL.Query().Get("code")
		if code == "" {
			errServer = fmt.Errorf("could not find 'code' URL parameter")
			return
		}

		// Trade the authorization code and the code verifier for access and refresh tokens
		accessToken, refreshToken, err := getUserAccessAndRefreshTokens(authDomain, clientId, codeVerifier, code, redirectURL)
		if err != nil {
			errServer = fmt.Errorf("retrieve tokens: %w", err)
			return
		}

		sessionExpiresAtUnix, err := getStartingSessionExpiresAtUnix()
		if err != nil {
			errServer = fmt.Errorf("compute session expiration timestamp: %w", err)
			return
		}

		err = SetAuthFlow(AUTH_FLOW_USER_TOKEN)
		if err != nil {
			errServer = fmt.Errorf("set auth flow type: %w", err)
			return
		}
		authFields := map[authFieldKey]string{
			SESSION_EXPIRES_AT_UNIX: sessionExpiresAtUnix,
			ACCESS_TOKEN:            accessToken,
			REFRESH_TOKEN:           refreshToken,
		}
		err = SetAuthFieldMap(authFields)
		if err != nil {
			errServer = fmt.Errorf("set in auth storage: %w", err)
		}

		// Return an indication of success to the caller
		_, _ = io.WriteString(w, `
		<html>
			<body>
				<h1>Login successful!</h1>
				<h2>You can close this window and return to the STACKIT CLI.</h2>
			</body>
		</html>`)

		// We can also directly redirect the user to another STACKIT page, or a link to documentation
		// openBrowser("https://www.stackit.de/en/")
	})

	// Parse the redirect URL for the port number
	u, err := url.Parse(redirectURL)
	if err != nil {
		return fmt.Errorf("parse redirect URL: %w", err)
	}

	// Set up a listener on the redirect port
	port := fmt.Sprintf(":%s", u.Port())
	l, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("listen to port %s: %w", port, err)
	}

	// Open a browser window to the authorizationURL
	err = openBrowser(authorizationURL)
	if err != nil {
		return fmt.Errorf("open browser to URL %s: %w", authorizationURL, err)
	}

	// Start the blocking web server loop
	// This will exit when the handler gets fired and calls server.Close()
	err = server.Serve(l)
	if !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("server for PKCE flow closed unexpectedly: %w", err)
	}

	// Check if there was an error in the HTTP server
	if errServer != nil {
		return fmt.Errorf("PKCE flow: %w", errServer)
	}

	return nil
}

// getUserAccessAndRefreshTokens trades the authorization code retrieved from the first OAuth2 leg for an access token and a refresh token
func getUserAccessAndRefreshTokens(authDomain, clientID, codeVerifier, authorizationCode, callbackURL string) (accessToken, refreshToken string, err error) {
	// Set the authUrl and form-encoded data for the POST to the access token endpoint
	authUrl := fmt.Sprintf("https://%s/token", authDomain)
	data := fmt.Sprintf(
		"grant_type=authorization_code&client_id=%s"+
			"&code_verifier=%s"+
			"&code=%s"+
			"&redirect_uri=%s",
		clientID, codeVerifier, authorizationCode, callbackURL)
	payload := strings.NewReader(data)

	// Create the request and execute it
	req, _ := http.NewRequest("POST", authUrl, payload)
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("call access token endpoint: %w", err)
	}

	// Process the response
	defer func() {
		closeErr := res.Body.Close()
		if closeErr != nil {
			err = fmt.Errorf("close response body: %w", closeErr)
		}
	}()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", "", fmt.Errorf("read response body: %w", err)
	}

	// Unmarshal the json into a string map
	responseData := struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}{}
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return "", "", fmt.Errorf("unmarshal response: %w", err)
	}
	if responseData.AccessToken == "" {
		return "", "", fmt.Errorf("found no access token")
	}
	if responseData.RefreshToken == "" {
		return "", "", fmt.Errorf("found no refresh token")
	}

	return responseData.AccessToken, responseData.RefreshToken, nil
}

// Cleanup closes the HTTP server
func cleanup(server *http.Server) {
	// We run this as a goroutine so that this function falls through and
	// the socket to the browser gets flushed/closed before the server goes away
	go func() {
		_ = server.Close()
	}()
}

func openBrowser(pageUrl string) error {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", pageUrl).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", pageUrl).Start()
	case "darwin":
		err = exec.Command("open", pageUrl).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		return err
	}
	return nil
}
