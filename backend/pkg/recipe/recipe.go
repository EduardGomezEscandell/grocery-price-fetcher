package recipe

import (
	"math/rand"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/product"
)

type ID uint32

const IDSize = 32

// Recipe represents a combination of ingredients that can be used to prepare a dish.
type Recipe struct {
	ID          ID
	Name        string
	Ingredients []Ingredient
}

// Ingredient represents a single ingredient that is part of a recipe.
type Ingredient struct {
	ProductID product.ID `json:"product_id"`
	Amount    float32    `json:"amount"`
}

func NewRandomID() ID {
	//nolint:gosec // This is not used for security purposes
	return ID(rand.Uint32())
}
