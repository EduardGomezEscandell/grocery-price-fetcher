package e2e

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/database/dbtypes"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/database/mysql"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/logger"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/recipe"
	"github.com/sirupsen/logrus"
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

const mysqlRootPassword = "test-db-password"

func Setup(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	// Stop the service and remove data
	if _, err := Make(ctx, "clean"); err != nil {
		return fmt.Errorf("could not uninstall service: %v", err)
	}

	start := time.Now()

	ch := make(chan error)
	go func() {
		defer close(ch)
		err := func() error {
			if _, err := Make(ctx, "build-docker"); err != nil {
				return fmt.Errorf("could not build container: %v", err)
			}

			if _, err := Make(ctx, "install"); err != nil {
				return fmt.Errorf("could not install service: %v", err)
			}

			if err := writeConfig(ctx, "db_root_password.txt", mysqlRootPassword); err != nil {
				return err
			}

			if err := writeConfig(ctx, "google_client_secret.txt", "made-up-123"); err != nil {
				return err
			}

			if _, err := Make(ctx, "start"); err != nil {
				return fmt.Errorf("could not start service: %v", err)
			}

			return nil
		}()

		if err != nil {
			cancel()
			ch <- err
		}
	}()

	if err := waitServiceReady(ctx, start); err != nil {
		cancel()
		return errors.Join(err, <-ch)
	}

	select {
	case err := <-ch:
		if err != nil {
			return err
		}
	case <-time.After(5 * time.Second):
		return errors.New("timeout waiting for setup")
	}

	if err := insertTestUser(); err != nil {
		return fmt.Errorf("could not insert a test user and user session: %v", err)
	}

	return nil
}

func writeConfig(ctx context.Context, path, text string) error {
	p := filepath.Join("/etc/grocery-price-fetcher", path)
	script := fmt.Sprintf("printf %q > %q", text, p)
	out, err := exec.CommandContext(ctx, "sudo", "bash", "-c", script).CombinedOutput()
	if err != nil {
		return fmt.Errorf("could not set config %q: %v: %s", path, err, out)
	}

	return nil
}

func waitServiceReady(ctx context.Context, since time.Time) error {
	//nolint:gosec // Input arguments are not user input
	cmd := exec.CommandContext(ctx,
		"journalctl", "--no-pager",
		"-fu", "grocery-price-fetcher.service",
		"--output", "cat",
		"--since", since.Format("2006-01-02 15:04:05"))

	r, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("journalctl: %w", err)
	}

	cmd.Stderr = cmd.Stdout

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("journalctl: %w", err)
	}
	defer cmd.Wait() //nolint:errcheck // we don't care about the error

	sc := bufio.NewScanner(r)
	for sc.Scan() {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout: %w", ctx.Err())
		default:
		}

		fmt.Println(sc.Text())
		if strings.Contains(sc.Text(), "Listening on [::]:443") {
			return nil
		}
	}

	if err := sc.Err(); err != nil {
		return fmt.Errorf("error reading journalctl: %v", err)
	}

	return errors.New("unexpected end of journalctl")
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

const TestSessionID = "test-session-ID"
const TestUser = "test-user-ID"

func insertTestUser() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	f, err := os.CreateTemp(os.TempDir(), "db_root_password")
	if err != nil {
		return err
	}
	defer os.RemoveAll(f.Name())

	if _, err := f.WriteString(mysqlRootPassword); err != nil {
		f.Close()
		return err
	}
	f.Close()

	sett := mysql.DefaultSettings()
	sett.PasswordFile = f.Name()

	log := logger.New()
	log.SetOutput(os.Stderr)
	log.SetLevel(int(logrus.TraceLevel))

	db, err := mysql.New(ctx, log, sett)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := db.SetUser(TestUser); err != nil {
		return err
	}

	err = db.SetSession(dbtypes.Session{
		ID:          TestSessionID,
		User:        TestUser,
		AccessToken: "123",
		NotAfter:    time.Now().Add(24 * time.Hour),
	})
	if err != nil {
		return err
	}

	// Insert some sample data for this user
	rID, err := db.SetRecipe(recipe.Recipe{
		User: TestUser,
		Name: "Test Recipe",
		Ingredients: []recipe.Ingredient{{
			ProductID: 1,
			Amount:    5,
		}},
	})
	if err != nil {
		return err
	}

	err = db.SetMenu(dbtypes.Menu{
		User: TestUser,
		Name: "default",
		Days: []dbtypes.Day{{
			Name: "Monday",
			Meals: []dbtypes.Meal{{
				Name: "Lunch",
				Dishes: []dbtypes.Dish{{
					ID:     rID,
					Amount: 2,
				}},
			}},
		}},
	})
	if err != nil {
		return err
	}

	return nil
}
