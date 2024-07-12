package e2e_test

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"path"
	"testing"

	e2e "github.com/EduardGomezEscandell/grocery-price-fetcher/end-to-end"
	"github.com/stretchr/testify/require"
)

func TestHelloWorld(t *testing.T) {
	t.Parallel()
	resp, err := request(t, http.MethodGet, "api/helloworld", nil)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode, "Unexpected status code %s", resp.Status)
	require.Equal(t, "Hello, world!\n", resp.Body)
}

func TestVersion(t *testing.T) {
	t.Parallel()
	resp, err := request(t, http.MethodGet, "api/version", nil)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode, "Unexpected status code %s", resp.Status)
	require.NotContains(t, resp.Body, "Dev", "Version was not properly set during build")
}

func TestRefreshLogin(t *testing.T) {
	t.Parallel()
	resp, err := request(t, http.MethodPost, "api/auth/refresh", nil, authHeader)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode, "Unexpected status code %s", resp.Status)
	require.NotEmpty(t, resp.Body, "Body should not be empty")
}

func TestRecipes(t *testing.T) {
	t.Parallel()
	resp, err := request(t, http.MethodGet, "api/recipes", nil, authHeader)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode, "Unexpected status code %s", resp.Status)
	require.NotEmpty(t, resp.Body, "Body should not be empty")
}

func TestMenu(t *testing.T) {
	t.Parallel()
	resp, err := request(t, http.MethodGet, "api/menu/default", nil, authHeader)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode, "Unexpected status code %s", resp.Status)
	require.NotEmpty(t, resp.Body, "Body should not be empty")
}

func TestFrontEnd(t *testing.T) {
	t.Parallel()
	resp, err := request(t, http.MethodGet, "", nil, browserHeaders...)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode, "Unexpected status code %s", resp.Status)
	require.NotEmpty(t, resp.Body, "Body should not be empty")
	require.Contains(t, resp.Body, "</html>", "Frontend should contain the HTML tag")
}

func TestFrontEndRouting(t *testing.T) {
	t.Parallel()
	resp, err := request(t, http.MethodGet, "menu", nil, browserHeaders...)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode, "Unexpected status code %s", resp.Status)
	require.NotEmpty(t, resp.Body, "Body should not be empty")
	require.Contains(t, resp.Body, "</html>", "Frontend should contain the HTML tag")
}

var browserHeaders = []kv{
	{
		key:   "Accept",
		value: "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
	},
	{
		key:   "Accept-Encoding",
		value: "gzip, deflate, br, zstd",
	},
	{
		key:   "User-Agent",
		value: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36",
	},
}

var authHeader = kv{
	key:   "Authorization",
	value: fmt.Sprintf("Bearer %s", e2e.TestSessionID),
}

type kv struct {
	key   string
	value string
}

//nolint:unparam // Method is always GET in this test
func request(t *testing.T, method string, endpoint string, body []byte, headers ...kv) (*response, error) {
	t.Helper()

	var buff bytes.Buffer
	_, err := buff.Write(body)
	if err != nil {
		return nil, fmt.Errorf("could not write body into buffer: %v", err)
	}

	url := "https://" + path.Join("localhost", endpoint)
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %v", err)
	}

	for _, h := range headers {
		req.Header.Add(h.key, h.value)
	}

	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				//nolint:gosec // InsecureSkipVerify is used to avoid certificate validation in tests
				InsecureSkipVerify: true,
			},
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not perform request: %v", err)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response body: %v", err)
	}

	t.Logf("Response to %s %s: %s\n%s", method, url, resp.Status, string(b))

	return &response{
		Body:     string(b),
		Response: resp,
	}, nil
}

type response struct {
	Body string
	*http.Response
}
