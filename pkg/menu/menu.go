package menu

import (
	"encoding/json"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/database"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/recipe"

	log "github.com/sirupsen/logrus"
)

type Day struct {
	Name  string
	Meals []Meal
}

type Meal struct {
	Name    string
	Recipes []struct {
		Name   string
		Amount float32
	}
}

type Menu struct {
	Days []Day
}

func (menu *Menu) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &menu.Days)
}

type ProductData struct {
	Name   string
	Amount float32
	Cost   float32
}

func (menu Menu) Compute(db *database.DB) ([]ProductData, error) {
	counts := make(map[string]float32)
	for _, r := range db.Products {
		counts[r.Name] = 0
	}

	// Calculate the amount of each recipe consumed
	recipes := make(map[*recipe.Recipe]float32)
	for _, day := range menu.Days {
		for _, meal := range day.Meals {
			for _, recipe := range meal.Recipes {
				rpe, ok := db.LookupRecipe(recipe.Name)
				if !ok {
					log.Warningf("%s: %s: Recipe %q is not registered", day.Name, meal.Name, recipe.Name)
				}
				recipes[rpe] = recipes[rpe] + recipe.Amount
			}
		}
	}

	// Calculate the amount of each product needed
	for rec, amount := range recipes {
		for _, i := range rec.Ingredients {
			_, ok := counts[i.Name]
			if !ok {
				log.Warningf("Recipe %q contains product %q which is not registered", rec.Name, i.Name)
				continue
			}
			counts[i.Name] += amount * i.Amount
		}
	}

	// Create the output
	products := make([]ProductData, 0, len(counts))
	for _, p := range db.Products {
		products = append(products, ProductData{
			Name:   p.Name,
			Amount: counts[p.Name],
			Cost:   counts[p.Name] * p.Price,
		})
	}

	return products, nil
}
