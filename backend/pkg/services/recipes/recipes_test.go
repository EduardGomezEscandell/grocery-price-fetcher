package recipes_test

import (
	"net/http"
	"testing"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/providers"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/providers/blank"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/services/recipes"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/testutils"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/types"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	providers.Register(blank.Provider{})
	m.Run()
}

func TestHelloWorld(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		method  string
		emptyDB bool

		wantCode int
		wantBody string
	}{
		"GET":               {method: "GET", wantCode: http.StatusOK, wantBody: "[\"Water\",\"Juice\"]\n"},
		"GET with empty DB": {method: "GET", emptyDB: true, wantCode: http.StatusOK, wantBody: "[]\n"},
		"POST":              {method: "POST", wantCode: http.StatusMethodNotAllowed},
		"PATCH":             {method: "PATCH", wantCode: http.StatusMethodNotAllowed},
		"PUT":               {method: "PUT", wantCode: http.StatusMethodNotAllowed},
		"DELETE":            {method: "DELETE", wantCode: http.StatusMethodNotAllowed},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			db := testutils.Database(t, "")
			if !tc.emptyDB {
				err := db.SetRecipe(types.Recipe{Name: "Water", Ingredients: []types.Ingredient{
					{Name: "Hydrogen", Amount: 2},
					{Name: "Oxygen", Amount: 1},
				}})
				require.NoError(t, err)

				err = db.SetRecipe(types.Recipe{Name: "Juice", Ingredients: []types.Ingredient{
					{Name: "Orange", Amount: 2.12},
				}})
				require.NoError(t, err)
			}

			sv := recipes.New(recipes.Settings{}.Defaults(), db)
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
