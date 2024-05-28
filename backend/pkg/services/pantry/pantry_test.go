package pantry_test

import (
	"net/http"
	"os"
	"testing"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/providers"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/providers/blank"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/services/pantry"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/testutils"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	providers.Register(blank.Provider{})
	m.Run()
}

func TestPantryEndpoint(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		method string

		wantCode int
		wantBody string
	}{
		"GET":    {method: "GET", wantCode: http.StatusOK, wantBody: "!golden"},
		"POST":   {method: "POST", wantCode: http.StatusCreated},
		"PATCH":  {method: "PATCH", wantCode: http.StatusMethodNotAllowed},
		"PUT":    {method: "PUT", wantCode: http.StatusMethodNotAllowed},
		"DELETE": {method: "DELETE", wantCode: http.StatusMethodNotAllowed},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			db := testutils.Database(t, testutils.FixturePath(t, "database"))

			sv := pantry.New(pantry.Settings{}.Defaults(), db)
			require.True(t, sv.Enabled())

			fixture := testutils.FixturePath(t, "message", "body.json")
			out, err := os.ReadFile(fixture)
			if err != nil {
				require.ErrorIs(t, err, os.ErrNotExist)
				out = nil
				t.Logf("No golden file found at %s", fixture)
			}

			testutils.TestEndpoint(t, testutils.ResponseTestOptions{
				Path:     "/api/pantry",
				Endpoint: sv.Handle,
				Method:   tc.method,
				Body:     string(out),
				WantCode: tc.wantCode,
				WantBody: tc.wantBody,
			})
		})
	}
}
