package dbtypes

// Ingredient represents a single ingredient that is part of a recipe.
type Ingredient struct {
	Name   string  `json:"name"`
	Amount float32 `json:"amount"`
}

// Recipe represents a recipe that can be used to prepare a dish.
type Recipe struct {
	Name        string
	Ingredients []Ingredient
}

// Dish represents a single dish that is part of a meal.
type Dish struct {
	Name   string  `json:"name"`
	Amount float32 `json:"amount"`
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
	Name string `json:"name"`
	Days []Day  `json:"days,omitempty"`
}

type Pantry struct {
	Name     string       `json:"name"`
	Contents []Ingredient `json:"contents"`
}

type ShoppingList struct {
	Menu     string   `json:"menu"`
	Pantry   string   `json:"pantry"`
	Contents []string `json:"contents"`
}
