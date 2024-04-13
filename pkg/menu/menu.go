package menu

import (
	"encoding/json"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/database"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/recipe"
	log "github.com/sirupsen/logrus"
)

type Recipe struct {
	Name   string
	Amount float32
}

type Meal struct {
	Name    string
	Recipes []Recipe `json:",omitempty"`
}

type Day struct {
	Name  string
	Meals []Meal `json:",omitempty"`
}

type Menu struct {
	Days []Day `json:",omitempty"`
}

func (menu *Menu) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &menu.Days)
}

func (menu Menu) MarshalJSON() ([]byte, error) {
	return json.Marshal(menu.Days)
}

type ProductData struct {
	Name   string
	Amount float32 `json:",omitempty"`
	Cost   float32 `json:",omitempty"`
}

func (menu Menu) Compute(db *database.DB, pantry []ProductData) ([]ProductData, error) {
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

	// Subtract the amount of products in the pantry
	for _, p := range pantry {
		_, ok := products[p.Name]
		if !ok {
			continue
		}
		products[p.Name] = max(0, products[p.Name]-p.Amount)
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
