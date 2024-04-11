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
	// Compute the amount of each recipe needed
	recipes := make(map[*recipe.Recipe]float32)
	for _, day := range menu.Days {
		for _, meal := range day.Meals {
			for _, recipe := range meal.Recipes {
				rpe, ok := db.LookupRecipe(recipe.Name)
				if !ok {
					log.Warningf("%s: %s: Recipe %q is not registered", day.Name, meal.Name, recipe.Name)
					continue
				}
				recipes[rpe] = recipes[rpe] + recipe.Amount
			}
		}
	}

	// Compute the amount of each product needed
	products := make(map[string]float32)
	for rec, amount := range recipes {
		for _, i := range rec.Ingredients {
			_, ok := products[i.Name]
			if !ok {
				products[i.Name] = 0
			}
			products[i.Name] += amount * i.Amount
		}
	}

	// Asseble the output
	table := make([]ProductData, 0, len(products))
	for _, p := range db.Products {
		table = append(table, ProductData{
			Name:   p.Name,
			Amount: products[p.Name],
			Cost:   products[p.Name] * p.Price,
		})
	}

	return table, nil
}
