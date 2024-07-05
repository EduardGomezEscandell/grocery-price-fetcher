package shoppinglist_test

import (
	"net/http"
	"os"
	"testing"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/providers"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/providers/blank"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/services/shoppinglist"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/testutils"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	providers.Register(blank.Provider{})
	m.Run()
}

func TestShoppingEndpoint(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		method string

		wantCode int
		wantBody string
	}{
		"GET":    {method: "GET", wantCode: http.StatusOK, wantBody: "!golden"},
		"PUT":    {method: "PUT", wantCode: http.StatusCreated},
		"DELETE": {method: "DELETE", wantCode: http.StatusNoContent},
		"PATCH":  {method: "PATCH", wantCode: http.StatusMethodNotAllowed},
		"POST":   {method: "POST", wantCode: http.StatusMethodNotAllowed},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			db := testutils.Database(t, testutils.FixturePath(t, "database"))

			sv := shoppinglist.New(shoppinglist.Settings{}.Defaults(), db, testutils.MockAuthGetter())
			require.True(t, sv.Enabled())

			fixture := testutils.FixturePath(t, "message", "body.json")
			out, err := os.ReadFile(fixture)
			if err != nil {
				require.ErrorIs(t, err, os.ErrNotExist)
				out = nil
				t.Logf("No golden file found at %s", fixture)
			}

			testutils.TestEndpoint(t, testutils.ResponseTestOptions{
				ServePath: sv.Path(),
				ReqPath:   "/api/shopping-list/testmenu1/testpantry1",
				Endpoint:  sv.Handle,
				Method:    tc.method,
				Body:      string(out),
				WantCode:  tc.wantCode,
				WantBody:  tc.wantBody,
			})
		})
	}
}
