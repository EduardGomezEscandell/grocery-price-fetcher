package menuneeds

import (
	"slices"
	"strings"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/database"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/database/dbtypes"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/logger"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/product"
)

type RecipeItem struct {
	Product product.Product `json:"product"`
	Amount  float32         `json:"amount"`
}

type Needs struct {
	Menu   *dbtypes.Menu
	Pantry *dbtypes.Pantry
	Items  []RecipeItem
}

func ComputeNeeds(log logger.Logger, db database.DB, m *dbtypes.Menu) Needs {
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
	need := make(map[string]float32)
	for _, rec := range recipes {
		for _, i := range rec.recipe.Ingredients {
			_, ok := need[i.Name]
			if !ok {
				need[i.Name] = 0
			}
			need[i.Name] += rec.amount * i.Amount
		}
	}

	result := Needs{
		Menu:   m,
		Pantry: &dbtypes.Pantry{},
		Items:  make([]RecipeItem, 0, len(need)),
	}

	for name, amount := range need {
		product, ok := db.LookupProduct(name)
		if !ok {
			log.Warningf("Product %q is not registered", name)
			continue
		}

		result.Items = append(result.Items, RecipeItem{
			Product: product,
			Amount:  amount,
		})
	}

	slices.SortFunc(result.Items, func(i, j RecipeItem) int {
		return strings.Compare(i.Product.Name, j.Product.Name)
	})

	return result
}

func (n *Needs) Subtract(p *dbtypes.Pantry) {
	n.Pantry = p

	var i int
	var j int
	for i < len(n.Items) && j < len(p.Contents) {
		need := &n.Items[i]
		stock := p.Contents[j]

		switch strings.Compare(need.Product.Name, stock.Name) {
		case -1:
			i++
		case 1:
			j++
		default:
			need.Amount = max(need.Amount-stock.Amount, 0)
			i++
			j++
		}
	}
}
