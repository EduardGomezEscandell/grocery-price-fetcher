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
	e2e.CompareToGolden(t, "testdata/server/result.txt", resp.Body)
}

func TestVersion(t *testing.T) {
	t.Parallel()
	resp, err := request(t, http.MethodGet, "api/version", nil)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode, "Unexpected status code %s", resp.Status)
	require.NotContains(t, resp.Body, "Dev", "Version was not properly set during build")
}

func TestRecipes(t *testing.T) {
	t.Parallel()
	resp, err := request(t, http.MethodGet, "api/recipes", nil)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode, "Unexpected status code %s", resp.Status)
	require.NotEmpty(t, resp.Body, "Body should not be empty")
}

func TestMenu(t *testing.T) {
	t.Parallel()
	resp, err := request(t, http.MethodGet, "api/menu/default", nil)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode, "Unexpected status code %s", resp.Status)
	require.NotEmpty(t, resp.Body, "Body should not be empty")
}

//nolint:unparam // Method is always GET in this test
func request(t *testing.T, method string, endpoint string, body []byte) (*response, error) {
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
