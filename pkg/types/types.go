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
	Name   string
	Amount float32
}

// Meal represents a meal that is part of a day.
type Meal struct {
	Name   string
	Dishes []Dish `json:",omitempty"`
}

// Day represents a day of the week.
type Day struct {
	Name  string
	Meals []Meal `json:",omitempty"`
}

// Menu represents a menu for a week.
type Menu struct {
	Name string
	Menu []Day `json:",omitempty"`
}
