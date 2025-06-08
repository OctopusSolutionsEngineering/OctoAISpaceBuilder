package router

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type HeaderRoundTripper struct {
	Transport http.RoundTripper
	Headers   map[string]string
}

func (h *HeaderRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	for key, value := range h.Headers {
		req.Header.Set(key, value)
	}
	return h.Transport.RoundTrip(req)
}

func GetHttpClient(octopusUrl string) (*http.Client, error) {
	if !isDirectlyAccessibleOctopusInstance(octopusUrl) {
		fmt.Println("[SPACEBUILDER] Enabled Octopus AI Assistant redirection service")
		return createHttpClient(octopusUrl)
	}

	fmt.Println("[SPACEBUILDER] Did not enable Octopus AI Assistant redirection service")

	return nil, nil
}

// isDirectlyAccessibleOctopusInstance determines if the host should be contacted directly
func isDirectlyAccessibleOctopusInstance(octopusUrl string) bool {
	serviceEnabled, found := os.LookupEnv("REDIRECTION_SERVICE_ENABLED")

	if !found || serviceEnabled != "true" {
		return true
	}

	parsedUrl, err := url.Parse(octopusUrl)

	// Contact the server directly if the URL is invalid
	if err != nil {
		return true
	}

	return strings.HasSuffix(parsedUrl.Hostname(), ".octopus.app") ||
		strings.HasSuffix(parsedUrl.Hostname(), ".testoctopus.com") ||
		parsedUrl.Hostname() == "localhost" ||
		parsedUrl.Hostname() == "127.0.0.1"
}

func createHttpClient(octopusUrl string) (*http.Client, error) {

	serviceApiKey, found := os.LookupEnv("REDIRECTION_SERVICE_API_KEY")

	if !found {
		return nil, errors.New("REDIRECTION_SERVICE_API_KEY is required")
	}

	parsedUrl, err := url.Parse(octopusUrl)

	if err != nil {
		return nil, err
	}

	headers := map[string]string{
		"X_REDIRECTION_UPSTREAM_HOST":   parsedUrl.Hostname(),
		"X_REDIRECTION_SERVICE_API_KEY": serviceApiKey,
	}

	return &http.Client{
		Transport: &HeaderRoundTripper{
			Transport: http.DefaultTransport,
			Headers:   headers,
		},
	}, nil
}
