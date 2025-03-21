package auth

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"

	"golang.org/x/oauth2"

	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
)

const (
	defaultWellKnownConfig = "https://accounts.stackit.cloud/.well-known/openid-configuration"
	defaultCLIClientID     = "stackit-cli-0000-0000-000000000001"

	loginSuccessPath        = "/login-successful"
	stackitLandingPage      = "https://www.stackit.de"
	htmlTemplatesPath       = "templates"
	loginSuccessfulHTMLFile = "login-successful.html"
	logoPath                = "/stackit_nav_logo_light.svg"
	logoSVGFilePath         = "stackit_nav_logo_light.svg"

	// The IDP doesn't support wildcards for the port,
	// so we configure a range of ports from 8000 to 8020
	defaultPort         = 8000
	configuredPortRange = 20
)

//go:embed templates/*
var htmlContent embed.FS

type User struct {
	Email string
}

type apiClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// AuthorizeUser implements the PKCE OAuth2 flow.
func AuthorizeUser(p *print.Printer, isReauthentication bool) error {
	idpWellKnownConfigURL, err := getIDPWellKnownConfigURL()
	if err != nil {
		return fmt.Errorf("get IDP well-known configuration: %w", err)
	}
	if idpWellKnownConfigURL != defaultWellKnownConfig {
		p.Warn("You are using a custom identity provider well-known configuration (%s) for authentication.\n", idpWellKnownConfigURL)
		err := p.PromptForEnter("Press Enter to proceed with the login...")
		if err != nil {
			return err
		}
	}

	p.Debug(print.DebugLevel, "get IDP well-known configuration from %s", idpWellKnownConfigURL)
	httpClient := &http.Client{}
	idpWellKnownConfig, err := parseWellKnownConfiguration(httpClient, idpWellKnownConfigURL)
	if err != nil {
		return fmt.Errorf("parse IDP well-known configuration: %w", err)
	}

	idpClientID, err := getIDPClientID()
	if err != nil {
		return err
	}
	if idpClientID != defaultCLIClientID {
		p.Warn("You are using a custom client ID (%s) for authentication.\n", idpClientID)
		err := p.PromptForEnter("Press Enter to proceed with the login...")
		if err != nil {
			return err
		}
	}

	if isReauthentication {
		err := p.PromptForEnter("Your session has expired, press Enter to login again...")
		if err != nil {
			return err
		}
	}

	var redirectURL string
	var listener net.Listener
	var listenerErr error
	var port int
	for i := range configuredPortRange {
		port = defaultPort + i
		portString := fmt.Sprintf(":%s", strconv.Itoa(port))
		p.Debug(print.DebugLevel, "trying to bind port %d for login redirect", port)
		listener, listenerErr = net.Listen("tcp", portString)
		if listenerErr == nil {
			redirectURL = fmt.Sprintf("http://localhost:%d", port)
			p.Debug(print.DebugLevel, "bound port %d for login redirect", port)
			break
		}
		p.Debug(print.DebugLevel, "unable to bind port %d for login redirect: %s", port, listenerErr)
	}
	if listenerErr != nil {
		return fmt.Errorf("unable to bind port for login redirect, tried from port %d to %d: %w", defaultPort, port, err)
	}

	conf := &oauth2.Config{
		ClientID: idpClientID,
		Endpoint: oauth2.Endpoint{
			AuthURL: idpWellKnownConfig.AuthorizationEndpoint,
		},
		Scopes:      []string{"openid offline_access email"},
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
		p.Debug(print.DebugLevel, "received request from authentication server")
		// Close the server only if there was an error
		// Otherwise, it will redirect to the successful login page
		defer func() {
			if errServer != nil {
				fmt.Println(errServer)
				cleanup(server)
			}
		}()

		// Get the authorization code
		code := r.URL.Query().Get("code")
		errDescription := r.URL.Query().Get("error_description")
		if code == "" {
			errServer = fmt.Errorf("could not find 'code' URL parameter")
			if errDescription != "" {
				errServer = fmt.Errorf("%w: %s", errServer, errDescription)
			}
			return
		}

		p.Debug(print.DebugLevel, "trading authorization code for access and refresh tokens")

		// Trade the authorization code and the code verifier for access and refresh tokens
		accessToken, refreshToken, err := getUserAccessAndRefreshTokens(idpWellKnownConfig, idpClientID, codeVerifier, code, redirectURL)
		if err != nil {
			errServer = fmt.Errorf("retrieve tokens: %w", err)
			return
		}

		p.Debug(print.DebugLevel, "received response from the authentication server")

		sessionExpiresAtUnix, err := getStartingSessionExpiresAtUnix()
		if err != nil {
			errServer = fmt.Errorf("compute session expiration timestamp: %w", err)
			return
		}

		sessionExpiresAtUnixInt, err := strconv.Atoi(sessionExpiresAtUnix)
		if err != nil {
			p.Debug(print.ErrorLevel, "parse session expiration value \"%s\": %s", sessionExpiresAtUnix, err)
		} else {
			sessionExpiresAt := time.Unix(int64(sessionExpiresAtUnixInt), 0)
			p.Debug(print.DebugLevel, "session expires at %s", sessionExpiresAt)
		}

		err = SetAuthFlow(AUTH_FLOW_USER_TOKEN)
		if err != nil {
			errServer = fmt.Errorf("set auth flow type: %w", err)
			return
		}

		email, err := getEmailFromToken(accessToken)
		if err != nil {
			errServer = fmt.Errorf("get email from access token: %w", err)
			return
		}

		p.Debug(print.DebugLevel, "user %s logged in successfully", email)

		err = LoginUser(email, accessToken, refreshToken, sessionExpiresAtUnix)
		if err != nil {
			errServer = fmt.Errorf("set in auth storage: %w", err)
			return
		}

		// Redirect the user to the successful login page
		loginSuccessURL := redirectURL + loginSuccessPath

		p.Debug(print.DebugLevel, "redirecting browser to login successful page")
		http.Redirect(w, r, loginSuccessURL, http.StatusSeeOther)
	})

	mux.HandleFunc(loginSuccessPath, func(w http.ResponseWriter, _ *http.Request) {

		email, err := GetAuthField(USER_EMAIL)
		if err != nil {
			errServer = fmt.Errorf("read user email: %w", err)
		}

		user := User{
			Email: email,
		}

		// ParseFS expects paths using forward slashes, even on Windows
		// See: https://github.com/golang/go/issues/44305#issuecomment-780111748
		htmlTemplate, err := template.ParseFS(htmlContent, path.Join(htmlTemplatesPath, loginSuccessfulHTMLFile))
		if err != nil {
			errServer = fmt.Errorf("parse html file: %w", err)
		}

		err = htmlTemplate.Execute(w, user)
		if err != nil {
			errServer = fmt.Errorf("render page: %w", err)
		}
	})

	mux.HandleFunc(logoPath, func(w http.ResponseWriter, _ *http.Request) {
		defer cleanup(server)

		img, err := htmlContent.ReadFile(path.Join(htmlTemplatesPath, logoSVGFilePath))
		if err != nil {
			errServer = fmt.Errorf("read logo file: %w", err)
		}
		w.Header().Set("Content-Type", "image/svg+xml")
		_, err = w.Write(img)
		if err != nil {
			return
		}
	})

	p.Debug(print.DebugLevel, "opening browser for authentication: %s", authorizationURL)
	p.Debug(print.DebugLevel, "using authentication server on %s", idpWellKnownConfig.Issuer)
	p.Debug(print.DebugLevel, "using client ID %s for authentication ", idpClientID)

	// Open a browser window to the authorizationURL
	err = openBrowser(authorizationURL)
	if err != nil {
		return fmt.Errorf("open browser to URL %s: %w", authorizationURL, err)
	}

	// Start the blocking web server loop
	// It will exit when the handlers get fired and call server.Close()
	p.Debug(print.DebugLevel, "listening for response from authentication server on %s", redirectURL)
	err = server.Serve(listener)
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
func getUserAccessAndRefreshTokens(idpWellKnownConfig *wellKnownConfig, clientID, codeVerifier, authorizationCode, callbackURL string) (accessToken, refreshToken string, err error) {
	// Set form-encoded data for the POST to the access token endpoint
	data := fmt.Sprintf(
		"grant_type=authorization_code&client_id=%s"+
			"&code_verifier=%s"+
			"&code=%s"+
			"&redirect_uri=%s",
		clientID, codeVerifier, authorizationCode, callbackURL)
	payload := strings.NewReader(data)

	// Create the request and execute it
	req, _ := http.NewRequest("POST", idpWellKnownConfig.TokenEndpoint, payload)
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	httpClient := &http.Client{}
	res, err := httpClient.Do(req)
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
		// We need to use the windows way on WSL, otherwise we do not pass query
		// parameters correctly. https://github.com/microsoft/WSL/issues/3832
		if _, ok := os.LookupEnv("WSL_DISTRO_NAME"); !ok {
			err = exec.Command("xdg-open", pageUrl).Start()
			break
		}
		fallthrough
	case "windows":
		err = exec.Command("rundll32.exe", "url.dll,FileProtocolHandler", pageUrl).Start()
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

// parseWellKnownConfiguration gets the well-known OpenID configuration from the provided URL and returns it as a JSON
// the method also stores the IDP token endpoint in the authentication storage
func parseWellKnownConfiguration(httpClient apiClient, wellKnownConfigURL string) (wellKnownConfig *wellKnownConfig, err error) {
	req, _ := http.NewRequest("GET", wellKnownConfigURL, http.NoBody)
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("make the request: %w", err)
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
		return nil, fmt.Errorf("read response body: %w", err)
	}

	err = json.Unmarshal(body, &wellKnownConfig)
	if err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}
	if wellKnownConfig == nil {
		return nil, fmt.Errorf("nil well-known configuration response")
	}
	if wellKnownConfig.Issuer == "" {
		return nil, fmt.Errorf("found no issuer")
	}
	if wellKnownConfig.AuthorizationEndpoint == "" {
		return nil, fmt.Errorf("found no authorization endpoint")
	}
	if wellKnownConfig.TokenEndpoint == "" {
		return nil, fmt.Errorf("found no token endpoint")
	}

	err = SetAuthField(IDP_TOKEN_ENDPOINT, wellKnownConfig.TokenEndpoint)
	if err != nil {
		return nil, fmt.Errorf("set token endpoint in the authentication storage: %w", err)
	}
	return wellKnownConfig, err
}
