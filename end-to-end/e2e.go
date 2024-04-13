package e2e

import (
	"os"
	"testing"

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
