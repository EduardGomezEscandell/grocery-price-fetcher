package menuneeds

import (
	"cmp"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/database"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/database/dbtypes"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/logger"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/utils"
)

func ComputeNeeds(log logger.Logger, db database.DB, m dbtypes.Menu) []dbtypes.Ingredient {
	type recipeAmount struct {
		recipe dbtypes.Recipe
		amount float32
	}

	// Compute the amount of each recipe needed
	recipes := make(map[string]recipeAmount)

	cached := database.NewCachedLookup(db.LookupRecipe)
	for _, day := range m.Days {
		for _, meal := range day.Meals {
			for _, dish := range meal.Dishes {
				rpe, ok := cached.Lookup(dish.Name)
				if !ok {
					log.Warningf("%s: %s: Recipe %q is not registered", day.Name, meal.Name, dish.Name)
					continue
				}
				recipes[rpe.Name] = recipeAmount{
					recipe: rpe,
					amount: recipes[rpe.Name].amount + dish.Amount,
				}
			}
		}
	}

	// Compute the amount of each product needed
	need := make(map[uint32]float32)
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
	out := make([]dbtypes.Ingredient, 0, len(need))
	for k, v := range need {
		out = append(out, dbtypes.Ingredient{
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
func Subtract(need []dbtypes.Ingredient, have []dbtypes.Ingredient) []dbtypes.Ingredient {
	items := make([]dbtypes.Ingredient, 0, len(need))

	utils.Zipper(need, have, func(n, h dbtypes.Ingredient) int { return cmp.Compare(n.ProductID, h.ProductID) },
		func(n dbtypes.Ingredient) {
			// This product is needed but not in the pantry
			items = append(items, n)
		},
		func(n, h dbtypes.Ingredient) {
			// This product is needed and in the pantry
			n.Amount = max(n.Amount-h.Amount, 0)
			items = append(items, n)
		},
		func(h dbtypes.Ingredient) {
			// This product is in the pantry but not needed
		})

	return items
}
