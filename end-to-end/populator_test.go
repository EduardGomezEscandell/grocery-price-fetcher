package e2e_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	e2e "github.com/EduardGomezEscandell/grocery-price-fetcher/end-to-end"
	"github.com/stretchr/testify/require"
)

func TestPopulator(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)

	var (
		fixturesPath = filepath.Join(wd, "testdata/populator/input/data")
		composePath  = filepath.Join(wd, "testdata/populator/compose.yaml")
	)

	output := t.TempDir()

	defer dockerCompose(t, composePath, "down", output) //nolint:errcheck // We cannot do anything with the error
	err = dockerCompose(t, composePath, "up", output)
	require.NoError(t, err)

	// Ensure the output is the same as the input
	// Need to use sudo because the files are created by the root user
	out, err := exec.Command("sudo", "diff", "-r", fixturesPath, output).CombinedOutput()
	require.NoError(t, err, string(out))
}

func dockerCompose(t *testing.T, composePath string, verb string, outDir string) error {
	t.Helper()

	cmd := exec.Command("sudo", "-E", "docker", "compose", "--file", composePath, verb)
	cmd.Env = append(
		os.Environ(),
		"OUTPUT_DIR="+outDir)

	w := e2e.NewTestWriter(t)
	defer w.Close()

	cmd.Stdout = w
	cmd.Stderr = w

	return cmd.Run()
}
