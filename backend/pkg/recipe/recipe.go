package recipe

import "github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/product"


// Recipe represents a combination of ingredients that can be used to prepare a dish.
type Recipe struct {
	Name        string
	Ingredients []Ingredient
}

// Ingredient represents a single ingredient that is part of a recipe.
type Ingredient struct {
	ProductID product.ID `json:"product_id"`
	Amount    float32    `json:"amount"`
}
