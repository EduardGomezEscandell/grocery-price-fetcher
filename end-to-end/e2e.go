package e2e

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func CompareToGolden(t *testing.T, goldenPath string, got string) {
	t.Helper()

	out, err := os.ReadFile(goldenPath)
	require.NoError(t, err, "Could not read golden")

	if os.Getenv("UPDATE_GOLDEN") != "" {
		require.NoError(t, os.WriteFile(goldenPath, []byte(got), 0600), "Could not update golden")
	}

	want := string(out)
	require.Equal(t, want, got, "Generated file does not match golden")
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
	if out, err := Make(ctx, "stop"); err != nil {
		fmt.Fprintf(os.Stderr, "could not stop service: %v. %s\n", err, string(out))
	}

	if out, err := Make(ctx, "uninstall"); err != nil {
		return fmt.Errorf("could not uninstall service: %v. %s", err, string(out))
	}

	return nil
}

type TestWriter struct {
	w io.WriteCloser
	r io.ReadCloser
}

func NewTestWriter(t *testing.T) *TestWriter {
	t.Helper()
	r, w := io.Pipe()

	go func() {
		defer r.Close()

		sc := bufio.NewScanner(r)
		for sc.Scan() {
			t.Log(sc.Text())
		}

		if err := sc.Err(); err != nil {
			t.Error(err)
		}
	}()

	t.Cleanup(func() { w.Close() })

	return &TestWriter{
		w: w,
		r: r,
	}
}

func (w *TestWriter) Write(p []byte) (n int, err error) {
	return w.w.Write(p)
}

func (w *TestWriter) Close() error {
	return w.w.Close()
}
