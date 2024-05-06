package e2e_test

import (
	"bufio"
	"context"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
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

	defer requireMake(t, ctx, "stop")
	cmd := exec.CommandContext(ctx, "make", "run")
	cmd.Dir = ".."
	r, w := io.Pipe()
	cmd.Stdout = w
	cmd.Stderr = w

	started := make(chan struct{})
	exited := make(chan struct{})

	go func() {
		sc := bufio.NewScanner(r)
		for sc.Scan() {
			logLine := sc.Text()
			if strings.Contains(logLine, "Listening on ") {
				close(started)
			}
			t.Log(logLine)
		}
		if err := sc.Err(); err != nil {
			t.Logf("Error scanning server output: %v", err)
		}
	}()

	err := cmd.Start()
	require.NoError(t, err)
	defer cmd.Process.Kill() //nolint:errcheck // We don't really care if Kill fails

	go func() {
		err := cmd.Wait()
		if err != nil {
			t.Logf("Server exited with error: %v", err)
		}
		close(exited)
	}()

	select {
	case <-time.After(time.Minute):
		require.Fail(t, "Server did not start serving")
	case <-exited:
		require.Fail(t, "Server exited before serving")
	case <-started:
	}

	f, err := os.Open(payload)
	require.NoError(t, err)
	defer f.Close()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://localhost/api/menu", f)
	require.NoError(t, err)

	client := http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode, "Unexpected status code %s", http.StatusText(resp.StatusCode))

	rB, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	e2e.CompareToGolden(t, golden, string(rB))

	requireMake(t, ctx, "stop")
	<-exited
}

//nolint:revive // Context goes after t
func requireMake(t *testing.T, ctx context.Context, verb string) {
	t.Helper()
	cmd := exec.CommandContext(ctx, "make", verb)
	cmd.Dir = ".."
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "Failed to 'make %s' (%v): %s", verb, err, out)
}
