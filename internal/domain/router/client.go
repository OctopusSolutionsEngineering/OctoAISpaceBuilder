package router

import (
	"errors"
	"go.uber.org/zap"
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

func GetHttpClient(octopusUrl *url.URL) (*http.Client, *url.URL, error) {
	if !isDirectlyAccessibleOctopusInstance(octopusUrl) {
		zap.L().Info("Enabled Octopus AI Assistant redirection service")
		return createHttpClient(octopusUrl)
	}

	zap.L().Info("Did not enable Octopus AI Assistant redirection service")

	return nil, octopusUrl, nil
}

// isDirectlyAccessibleOctopusInstance determines if the host should be contacted directly
func isDirectlyAccessibleOctopusInstance(octopusUrl *url.URL) bool {
	serviceEnabled, found := os.LookupEnv("REDIRECTION_SERVICE_ENABLED")

	if !found || serviceEnabled != "true" {
		return true
	}

	return strings.HasSuffix(octopusUrl.Hostname(), ".octopus.app") ||
		strings.HasSuffix(octopusUrl.Hostname(), ".testoctopus.com") ||
		octopusUrl.Hostname() == "localhost" ||
		octopusUrl.Hostname() == "127.0.0.1"
}

func createHttpClient(octopusUrl *url.URL) (*http.Client, *url.URL, error) {

	serviceApiKey, found := os.LookupEnv("REDIRECTION_SERVICE_API_KEY")

	if !found {
		return nil, nil, errors.New("REDIRECTION_SERVICE_API_KEY is required")
	}

	redirectionHost, found := os.LookupEnv("REDIRECTION_HOST")

	if !found {
		return nil, nil, errors.New("REDIRECTION_HOST is required")
	}

	redirectionHostUrl, err := url.Parse("https://" + redirectionHost)

	if err != nil {
		return nil, nil, err
	}

	headers := map[string]string{
		"X_REDIRECTION_UPSTREAM_HOST":   octopusUrl.Hostname(),
		"X_REDIRECTION_SERVICE_API_KEY": serviceApiKey,
	}

	return &http.Client{
		Transport: &HeaderRoundTripper{
			Transport: http.DefaultTransport,
			Headers:   headers,
		},
	}, redirectionHostUrl, nil
}
