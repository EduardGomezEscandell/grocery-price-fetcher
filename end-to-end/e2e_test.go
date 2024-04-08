package e2e_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEndToEnd(t *testing.T) {
	const (
		input  = "example.tsv"
		format = "table"
		golden = "result.tsv"
	)

	output := filepath.Join(t.TempDir(), "output.tsv")

	out, err := exec.Command("../bin/compra", "-i", input, "-o", output).CombinedOutput()
	require.NoError(t, err, "Stdout+Stderr: %s", string(out))

	require.FileExists(t, output)
	out, err = os.ReadFile(output)
	require.NoError(t, err, "Could not read output file")
	got := string(out)

	out, err = os.ReadFile(golden)
	require.NoError(t, err, "Could not read golden")
	want := string(out)

	if os.Getenv("UPDATE_GOLDEN") != "" {
		require.NoError(t, os.WriteFile(golden, []byte(got), 0600), "Could not update golden")
	}

	require.Equal(t, want, got, "Generated file does not match golden")
}
