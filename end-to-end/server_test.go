package e2e_test

import (
	"context"
	"io"
	"net/http"
	"testing"
	"time"

	e2e "github.com/EduardGomezEscandell/grocery-price-fetcher/end-to-end"
	"github.com/stretchr/testify/require"
)

func TestHelloWorld(t *testing.T) {
	const (
		golden = "testdata/server/result.txt"
	)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost/api/helloworld", nil)
	require.NoError(t, err)

	client := http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode, "Unexpected status code %s", http.StatusText(resp.StatusCode))

	rB, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	e2e.CompareToGolden(t, golden, string(rB))
}
