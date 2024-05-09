package e2e_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	e2e "github.com/EduardGomezEscandell/grocery-price-fetcher/end-to-end"
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

	e2e.CompareToGolden(t, golden, got)
}
