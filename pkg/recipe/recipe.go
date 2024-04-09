package recipe

type Ingredient struct {
	Name   string
	Amount float32
}

type Recipe struct {
	Name        string
	Ingredients []Ingredient
}
