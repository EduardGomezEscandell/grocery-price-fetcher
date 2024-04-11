package e2e_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCompra(t *testing.T) {
	const (
		input  = "testdata/compra/example.tsv"
		format = "table"
		golden = "testdata/compra/result.tsv"
	)

	output := filepath.Join(t.TempDir(), "output.tsv")

	cmd := exec.Command("../bin/compra", "-i", input, "-o", output, "-v")
	cmd.Env = append(os.Environ(), "LC_NUMERIC=ca_ES.UTF8")

	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "Stdout+Stderr: %s", string(out))
	t.Logf("%s", out)

	require.FileExists(t, output)
	out, err = os.ReadFile(output)
	require.NoError(t, err, "Could not read output file")
	got := string(out)

	compareToGolden(t, golden, got)
}

func TestNeeds(t *testing.T) {
	const (
		inputDB = "testdata/needs/database.json"
		input   = "testdata/needs/menu.json"
		format  = "table"
		golden  = "testdata/needs/result.tsv"
	)

	output := filepath.Join(t.TempDir(), "output.tsv")

	cmd := exec.Command("../bin/needs", "-i", input, "-db", inputDB, "-o", output, "-v", "--skip-empty")
	cmd.Env = append(os.Environ(), "LC_NUMERIC=ca_ES.UTF8")

	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "Stdout+Stderr: %s", string(out))
	t.Logf("%s", out)

	require.FileExists(t, output)
	out, err = os.ReadFile(output)
	require.NoError(t, err, "Could not read output file")
	got := string(out)

	compareToGolden(t, golden, got)
}

func compareToGolden(t *testing.T, goldenPath string, got string) {
	t.Helper()

	out, err := os.ReadFile(goldenPath)
	require.NoError(t, err, "Could not read golden")

	if os.Getenv("UPDATE_GOLDEN") != "" {
		require.NoError(t, os.WriteFile(goldenPath, []byte(got), 0600), "Could not update golden")
	}

	want := string(out)
	require.Equal(t, want, got, "Generated file does not match golden")
}
