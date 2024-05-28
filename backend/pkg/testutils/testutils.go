package testutils

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/database"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/database/jsondb"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/httputils"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/logger"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/services/pricing"
	"github.com/stretchr/testify/require"
)

type ResponseTestOptions struct {
	Path     string
	Endpoint httputils.Handler

	Method string
	Body   string

	WantCode int
	WantBody string
}

func TestEndpoint(t *testing.T, opt ResponseTestOptions) {
	t.Helper()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	addr, stop := HTTPServer(ctx, t, opt.Path, opt.Endpoint)
	defer stop()

	resp := MakeRequest(t, opt.Method, addr, opt.Body)
	defer resp.Body.Close()

	out, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	require.Equal(t, opt.WantCode, resp.StatusCode)
	if opt.WantBody == "" {
		return
	}

	if opt.WantBody == "!golden" {
		CompareToGolden(t, string(out), "http_response.json")
	} else {
		require.Equal(t, opt.WantBody, string(out))
	}
}

func MakeRequest(t *testing.T, method, url string, body string) *http.Response {
	t.Helper()

	var buff bytes.Buffer
	fmt.Fprint(&buff, body)

	req, err := http.NewRequest(method, url, &buff)
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	return resp
}

const PingEndpoint = "/test_utils_api/ping"

func NewLogger(t *testing.T) logger.Logger {
	t.Helper()
	log := logger.New()

	r, w := io.Pipe()
	go func() {
		sc := bufio.NewScanner(r)
		for sc.Scan() {
			t.Log(sc.Text())
		}
	}()

	log.SetOutput(w)
	return log
}

func HTTPServer(ctx context.Context, t *testing.T, p string, handler httputils.Handler) (string, func()) {
	t.Helper()

	ctx, cancel := context.WithCancel(ctx)
	t.Cleanup(cancel)

	server := http.NewServeMux()
	server.HandleFunc(p, httputils.HandleRequest(NewLogger(t), handler))
	server.HandleFunc(PingEndpoint, func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })

	lis, err := (&net.ListenConfig{}).Listen(ctx, "tcp", "localhost:")
	require.NoError(t, err, "failed to listen")

	stop := context.AfterFunc(ctx, func() {
		_ = lis.Close()
	})

	ch := make(chan error)
	go func() {
		//nolint:gosec // this is a test helper
		ch <- http.Serve(lis, server)
	}()

	require.Eventually(t, func() bool {
		url := fmt.Sprintf("http://%s%s", lis.Addr().String(), PingEndpoint)
		//nolint:gosec // this is a test helper
		resp, err := http.Get(url)
		if err != nil {
			return false
		}
		defer resp.Body.Close()
		return resp.StatusCode == http.StatusOK
	}, 10*time.Second, 100*time.Millisecond, "Server never started")
	t.Logf("server started at %s", lis.Addr().String())

	addr := fmt.Sprintf("http://%s%s", lis.Addr().String(), p)

	return addr, func() {
		lis.Close()
		stop()
		err := <-ch
		t.Logf("server stopped: %v", err)
	}
}

func CopyDir(t *testing.T, from, to string) {
	t.Helper()
	//nolint:gosec // this is a test helper
	out, err := exec.Command("rsync", "-r", from+"/", to).CombinedOutput()
	require.NoError(t, err, "failed to copy directory: %s", out)
}

func Database(t *testing.T, from string) database.DB {
	t.Helper()
	dir := t.TempDir()

	if from != "" {
		CopyDir(t, from, dir)
	}

	opts := jsondb.DefaultSettingsPath(dir)

	db, err := jsondb.New(context.Background(), NewLogger(t), opts)
	require.NoError(t, err)

	pricing.OneShot(context.Background(), NewLogger(t), db)

	t.Cleanup(func() {
		require.NoError(t, db.Close())
	})

	return db
}

func FixturePath(t *testing.T, relative ...string) string {
	t.Helper()

	clean := strings.NewReplacer(
		"/", "_",
		"\\", "_",
		":", "_",
		" ", "_",
	).Replace(t.Name())

	path := append([]string{"testdata", clean}, relative...)
	return filepath.Join(path...)
}

func CompareToGolden(t *testing.T, got string, name string) {
	t.Helper()

	goldenPath := FixturePath(t, "golden", name)

	out, err := os.ReadFile(goldenPath)
	require.NoError(t, err, "Could not read golden")

	if os.Getenv("UPDATE_GOLDEN") != "" {
		require.NoError(t, os.WriteFile(goldenPath, []byte(got), 0600), "Could not update golden")
	}

	want := string(out)
	require.Equal(t, want, got, "Generated file does not match golden")
}
