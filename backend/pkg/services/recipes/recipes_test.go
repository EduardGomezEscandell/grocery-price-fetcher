package recipes_test

import (
	"net/http"
	"testing"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/providers"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/providers/blank"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/recipe"
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
		"GET":               {method: "GET", wantCode: http.StatusOK, wantBody: "!golden"},
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
				_, err := db.SetRecipe(recipe.Recipe{
					ID:   1,
					User: "test-user-123",
					Name: "Water",
					Ingredients: []recipe.Ingredient{
						{ProductID: 1, Amount: 2},
						{ProductID: 2, Amount: 1},
					}})
				require.NoError(t, err)

				_, err = db.SetRecipe(recipe.Recipe{
					ID:   2,
					User: "test-user-123",
					Name: "Juice",
					Ingredients: []recipe.Ingredient{
						{ProductID: 3, Amount: 2.12},
					}})
				require.NoError(t, err)
			}

			sv := recipes.New(recipes.Settings{}.Defaults(), db, testutils.MockAuthGetter())
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
