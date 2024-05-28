package version_test

import (
	"net/http"
	"testing"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/services/version"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/testutils"
	"github.com/stretchr/testify/require"
)

func TestHelloWorld(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		method string

		wantCode int
		wantBody string
	}{
		"GET":    {method: "GET", wantCode: http.StatusOK, wantBody: `{ "version": "dev" }`},
		"POST":   {method: "POST", wantCode: http.StatusMethodNotAllowed},
		"PATCH":  {method: "PATCH", wantCode: http.StatusMethodNotAllowed},
		"PUT":    {method: "PUT", wantCode: http.StatusMethodNotAllowed},
		"DELETE": {method: "DELETE", wantCode: http.StatusMethodNotAllowed},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			sv := version.New(version.Settings{}.Defaults())
			require.True(t, sv.Enabled())

			testutils.TestEndpoint(t, testutils.ResponseTestOptions{
				Path:     "/api/version",
				Endpoint: sv.Handle,
				Method:   tc.method,
				Body:     "",
				WantCode: tc.wantCode,
				WantBody: tc.wantBody,
			})
		})
	}
}
