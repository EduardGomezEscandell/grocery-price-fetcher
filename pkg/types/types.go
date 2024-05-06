package types

// Ingredient represents a single ingredient that is part of a recipe.
type Ingredient struct {
	Name   string
	Amount float32
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
	Menu []Day  `json:"menu,omitempty"`
}
