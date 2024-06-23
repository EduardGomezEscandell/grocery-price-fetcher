package recipes_test

import (
	"net/http"
	"testing"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/database/dbtypes"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/providers"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/providers/blank"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/services/recipes"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/testutils"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	providers.Register(blank.Provider{})
	m.Run()
}

func TestRecipes(t *testing.T) {
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
				err := db.SetRecipe(dbtypes.Recipe{Name: "Water", Ingredients: []dbtypes.Ingredient{
					{ProductID: 1, Amount: 2},
					{ProductID: 2, Amount: 1},
				}})
				require.NoError(t, err)

				err = db.SetRecipe(dbtypes.Recipe{Name: "Juice", Ingredients: []dbtypes.Ingredient{
					{ProductID: 3, Amount: 2.12},
				}})
				require.NoError(t, err)
			}

			sv := recipes.New(recipes.Settings{}.Defaults(), db)
			require.True(t, sv.Enabled())

			testutils.TestEndpoint(t, testutils.ResponseTestOptions{
				ServePath: sv.Path(),
				ReqPath:   "/api/recipes",
				Endpoint:  sv.Handle,
				Method:    tc.method,
				Body:      "",
				WantCode:  tc.wantCode,
				WantBody:  tc.wantBody,
			})
		})
	}
}
