package dbtypes

import (
	"time"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/product"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/recipe"
)

type Session struct {
	ID           string
	User         string
	AccessToken  string
	RefreshToken string
	NotAfter     time.Time
}

// Dish represents a single dish that is part of a meal.
type Dish struct {
	ID     recipe.ID `json:"recipe_id"`
	Amount float32   `json:"amount"`
}

// Meal represents a meal that is part of a day.
type Meal struct {
	Name   string `json:"name"`
	Dishes []Dish `json:"dishes,omitempty"`
}

// Day represents a day of the week.
type Day struct {
	Name  string `json:"name"`
	Meals []Meal `json:"meals,omitempty"`
}

// Menu represents a menu for a week.
type Menu struct {
	User string `json:"user"`
	Name string `json:"name"`
	Days []Day  `json:"days,omitempty"`
}

type Pantry struct {
	User     string              `json:"user"`
	Name     string              `json:"name"`
	Contents []recipe.Ingredient `json:"contents"`
}

type ShoppingList struct {
	User     string       `json:"user"`
	Menu     string       `json:"menu"`
	Pantry   string       `json:"pantry"`
	Contents []product.ID `json:"contents"`
}
