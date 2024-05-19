package e2e_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"

	e2e "github.com/EduardGomezEscandell/grocery-price-fetcher/end-to-end"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := Setup(ctx); err != nil {
		log.Fatalf("Setup: %v", err)
	}

	e := m.Run()
	defer os.Exit(e)

	if err := Cleanup(ctx); err != nil {
		log.Fatalf("Cleanup: %v", err)
	}
}

func TestServer(t *testing.T) {
	const (
		manifest = "testdata/server/manifest.yaml"
		payload  = "testdata/server/payload.json"
		golden   = "testdata/server/result.json"
	)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

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

func Setup(ctx context.Context) error {
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	if _, err := Make(ctx, "build-docker"); err != nil {
		return fmt.Errorf("could not build container: %v", err)
	}

	if _, err := Make(ctx, "install"); err != nil {
		return fmt.Errorf("could not install service: %v", err)
	}

	if _, err := Make(ctx, "start"); err != nil {
		return fmt.Errorf("could not start service: %v", err)
	}

	const tick = 5 * time.Second
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return errors.New("timed out waiting for server to come online")
		case <-time.After(tick):
		}

		ok, err := func() (bool, error) {
			ctx, cancel := context.WithTimeout(ctx, tick)
			defer cancel()

			out, err := exec.CommandContext(ctx,
				"journalctl", "--no-pager",
				"-u", "grocery-price-fetcher.service",
				"--since", timestamp).CombinedOutput()
			if err != nil {
				return false, fmt.Errorf("could not access journalctl: %v: %s", err, out)
			}
			if !bytes.Contains(out, []byte("Server: serving on [::]:3000")) {
				return false, nil
			}
			return true, nil
		}()
		if err != nil {
			return err
		} else if ok {
			break
		}
	}

	return nil
}

func Cleanup(ctx context.Context) error {
	_, err := Make(context.Background(), "uninstall")
	if err != nil {
		return fmt.Errorf("could not uninstall service: %v", err)
	}
	return nil
}
