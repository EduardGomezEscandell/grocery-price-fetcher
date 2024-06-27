package menuneeds

import (
	"cmp"
	"errors"
	"io/fs"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/database"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/database/dbtypes"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/logger"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/product"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/recipe"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/utils"
)

func ComputeNeeds(log logger.Logger, db database.DB, m dbtypes.Menu) []recipe.Ingredient {
	type recipeAmount struct {
		recipe recipe.Recipe
		amount float32
	}

	// Compute the amount of each recipe needed
	recipes := make(map[recipe.ID]recipeAmount)

	cached := database.NewCachedLookup(db.LookupRecipe)
	for _, day := range m.Days {
		for _, meal := range day.Meals {
			for _, dish := range meal.Dishes {
				rpe, err := cached.Lookup(dish.ID)
				if errors.Is(err, fs.ErrNotExist) {
					log.Warningf("%s: %s: Recipe %d is not registered", day.Name, meal.Name, dish.ID)
					continue
				}
				recipes[rpe.ID] = recipeAmount{
					recipe: rpe,
					amount: recipes[rpe.ID].amount + dish.Amount,
				}
			}
		}
	}

	// Compute the amount of each product needed
	need := make(map[product.ID]float32)
	for _, rec := range recipes {
		for _, i := range rec.recipe.Ingredients {
			_, ok := need[i.ProductID]
			if !ok {
				need[i.ProductID] = 0
			}
			need[i.ProductID] += rec.amount * i.Amount
		}
	}

	// Convert the map to a slice
	out := make([]recipe.Ingredient, 0, len(need))
	for k, v := range need {
		out = append(out, recipe.Ingredient{
			ProductID: k,
			Amount:    v,
		})
	}

	return out
}

// Subtract computes the difference between the needed ingredients and the ones in the pantry.
// It returns a list of ingredients that are needed but not in the pantry.
//
// The input slices need and have must be sorted by ProductID. The output slice is also sorted by ProductID.
func Subtract(need []recipe.Ingredient, have []recipe.Ingredient) []recipe.Ingredient {
	items := make([]recipe.Ingredient, 0, len(need))

	utils.Zipper(need, have, func(n, h recipe.Ingredient) int { return cmp.Compare(n.ProductID, h.ProductID) },
		func(n recipe.Ingredient) {
			// This product is needed but not in the pantry
			items = append(items, n)
		},
		func(n, h recipe.Ingredient) {
			// This product is needed and in the pantry
			n.Amount = max(n.Amount-h.Amount, 0)
			items = append(items, n)
		},
		func(h recipe.Ingredient) {
			// This product is in the pantry but not needed
		})

	return items
}
