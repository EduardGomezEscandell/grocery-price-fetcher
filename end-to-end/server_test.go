package e2e_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"

	e2e "github.com/EduardGomezEscandell/grocery-price-fetcher/end-to-end"
	"github.com/stretchr/testify/require"
)

func TestServer(t *testing.T) {
	const (
		manifest = "testdata/server/manifest.yaml"
		payload  = "testdata/server/payload.json"
		golden   = "testdata/server/result.json"
	)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	out, err := Make(ctx, "build-docker")
	require.NoError(t, err, "Could not build container")
	t.Log(string(out))

	timestamp := time.Now().Format("2006-01-02 15:04:05")

	out, err = Make(ctx, "start")
	require.NoError(t, err, "Could not start service")
	t.Log(string(out))

	t.Cleanup(func() {
		_, _ = Make(context.Background(), "clean")
	})

	const tick = 5 * time.Second
	require.Eventually(t, func() bool {
		ctx, cancel := context.WithTimeout(ctx, tick)
		defer cancel()

		out, err := exec.CommandContext(ctx,
			"journalctl", "--no-pager",
			"-u", "grocery-price-fetcher.service",
			"--since", timestamp).CombinedOutput()
		if err != nil {
			t.Logf("Could not access journalctl: %v: %s", err, out)
			return false
		}
		if bytes.Contains(out, []byte("Server: serving on [::]:3000")) {
			return true
		}
		return false
	}, 20*time.Second, tick, "Server did not start serving")

	f, err := os.Open(payload)
	require.NoError(t, err)
	defer f.Close()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://localhost/api/menu", f)
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode, "Unexpected status code %s", http.StatusText(resp.StatusCode))

	rB, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	e2e.CompareToGolden(t, golden, string(rB))

	_, err = Make(ctx, "stop")
	require.NoError(t, err)
}

func Make(ctx context.Context, verb string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, "make", verb)
	cmd.Dir = ".."
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("make %s: %w: %s", verb, err, out)
	}
	return out, nil
}
