package dbtestutils

import (
	"testing"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/database"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/product"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/providers/blank"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/types"
	"github.com/stretchr/testify/require"
)

func ProductsTest(t *testing.T, openDB func() database.DB) {
	t.Helper()

	db := openDB()
	defer db.Close()

	_, ok := db.LookupProduct("AAAAAAAAAAAAAAAAAAAAA")
	require.False(t, ok)

	product1 := product.Product{
		Name:      "Product #1",
		Price:     1.99,
		BatchSize: 11,
		Provider:  blank.Provider{},
	}

	product2 := product.Product{
		Name:      "Product #2",
		Price:     0.64,
		BatchSize: 99,
		Provider:  blank.Provider{},
	}

	require.NoError(t, db.SetProduct(product1), "Could not set Product")
	p, ok := db.LookupProduct(product1.Name)
	require.True(t, ok, "Could not find Product just created")
	require.Equal(t, product1, p, "Product does not match the one just created")

	product1.BatchSize = 20

	require.NoError(t, db.SetProduct(product1), "Could not override Product")
	p, ok = db.LookupProduct(product1.Name)
	require.True(t, ok, "Could not find Product just overridden created")
	require.Equal(t, product1, p, "Product does not match the one just overridden")

	require.NoError(t, db.SetProduct(product2), "Could not set Product")
	p, ok = db.LookupProduct(product2.Name)
	require.True(t, ok, "Could not find empty Product just created")
	require.Equal(t, product2, p, "Empty menu does not match the one just created")

	t.Log("Closing DB and reopening")
	require.NoError(t, db.Close())

	menus := db.Products()
	require.ElementsMatch(t, []product.Product{product1, product2}, menus, "Products do not match the ones just created")

	db = openDB()
	defer db.Close()

	_, ok = db.LookupProduct("AAAAAAAAAAAAAAAAAAAAA")
	require.False(t, ok)

	p, ok = db.LookupProduct(product1.Name)
	require.True(t, ok, "Could not find Product after reopening DB")
	require.Equal(t, product1, p, "Product does not match the one after reopening DB")

	p, ok = db.LookupProduct(product2.Name)
	require.True(t, ok, "Could not find empty Product after reopening DB")
	require.Equal(t, product2, p, "Empty menu does not match the one after reopening DB")

	require.ElementsMatch(t, []product.Product{product1, product2}, menus, "Products do not match the ones after reopening DB")

	require.NoError(t, db.Close())
}

//nolint:dupl // This is a test file, so it's normal to have similar functions
func RecipesTest(t *testing.T, openDB func() database.DB) {
	t.Helper()

	db := openDB()
	defer db.Close()

	_, ok := db.LookupProduct("AAAAAAAAAAAAAAAAAAAAA")
	require.False(t, ok)

	recipe1 := types.Recipe{
		Name: "Recipe #1",
		Ingredients: []types.Ingredient{
			{Name: "Ingredient #1", Amount: 1.0},
			{Name: "Ingredient #2", Amount: 2.0},
		},
	}

	recipe2 := types.Recipe{
		Name:        "Recipe #2",
		Ingredients: nil,
	}

	require.NoError(t, db.SetRecipe(recipe1), "Could not set Recipe")
	p, ok := db.LookupRecipe(recipe1.Name)
	require.True(t, ok, "Could not find Recipe just created")
	require.Equal(t, recipe1, p, "Recipe does not match the one just created")

	recipe1.Ingredients[0].Amount = 5.0

	require.NoError(t, db.SetRecipe(recipe1), "Could not override Recipe")
	p, ok = db.LookupRecipe(recipe1.Name)
	require.True(t, ok, "Could not find Recipe just overridden created")
	require.Equal(t, recipe1, p, "Recipe does not match the one just overridden")

	require.NoError(t, db.SetRecipe(recipe2), "Could not set Recipe")
	p, ok = db.LookupRecipe(recipe2.Name)
	require.True(t, ok, "Could not find empty Recipe just created")
	require.Equal(t, recipe2, p, "Empty menu does not match the one just created")

	t.Log("Closing DB and reopening")
	require.NoError(t, db.Close())

	menus := db.Recipes()
	require.ElementsMatch(t, []types.Recipe{recipe1, recipe2}, menus, "Recipes do not match the ones just created")

	db = openDB()
	defer db.Close()

	_, ok = db.LookupRecipe("AAAAAAAAAAAAAAAAAAAAA")
	require.False(t, ok)

	p, ok = db.LookupRecipe(recipe1.Name)
	require.True(t, ok, "Could not find Recipe after reopening DB")
	require.Equal(t, recipe1, p, "Recipe does not match the one after reopening DB")

	p, ok = db.LookupRecipe(recipe2.Name)
	require.True(t, ok, "Could not find empty Recipe after reopening DB")
	require.Equal(t, recipe2, p, "Empty menu does not match the one after reopening DB")

	require.ElementsMatch(t, []types.Recipe{recipe1, recipe2}, menus, "Recipes do not match the ones after reopening DB")

	require.NoError(t, db.Close())
}

func MenuTest(t *testing.T, openDB func() database.DB) {
	t.Helper()

	db := openDB()
	defer db.Close()

	_, ok := db.LookupMenu("AAAAAAAAAAAAAAAAAAAAA")
	require.False(t, ok)

	myMenu := types.Menu{
		Name: "myMenu",
		Days: []types.Day{
			{
				Name: "Segiunda-Feira",
				Meals: []types.Meal{
					{
						Name: "Café da Manhã",
						Dishes: []types.Dish{
							{
								Name:   "Torrada i suc de taronja",
								Amount: 16,
							},
						},
					},
				},
			},
		},
	}

	require.NoError(t, db.SetMenu(myMenu), "Could not set Menu")
	m, ok := db.LookupMenu(myMenu.Name)
	require.True(t, ok, "Could not find Menu just created")
	require.Equal(t, myMenu, m, "Menu does not match the one just created")

	myMenu.Days[0].Meals[0].Dishes[0].Amount = 20

	require.NoError(t, db.SetMenu(myMenu), "Could not override Menu")
	m, ok = db.LookupMenu(myMenu.Name)
	require.True(t, ok, "Could not find Menu just overridden created")
	require.Equal(t, myMenu, m, "Menu does not match the one just overridden")

	emptyMenu := types.Menu{
		Name: "Empty Menu",
	}

	require.NoError(t, db.SetMenu(emptyMenu), "Could not set Menu")
	m, ok = db.LookupMenu(emptyMenu.Name)
	require.True(t, ok, "Could not find empty Menu just created")
	require.Equal(t, emptyMenu, m, "Empty menu does not match the one just created")

	t.Log("Closing DB and reopening")
	require.NoError(t, db.Close())

	menus := db.Menus()
	require.ElementsMatch(t, []types.Menu{myMenu, emptyMenu}, menus, "Menus do not match the ones just created")

	db = openDB()
	defer db.Close()

	_, ok = db.LookupMenu("AAAAAAAAAAAAAAAAAAAAA")
	require.False(t, ok)

	m, ok = db.LookupMenu(myMenu.Name)
	require.True(t, ok, "Could not find Menu after reopening DB")
	require.Equal(t, myMenu, m, "Menu does not match the one after reopening DB")

	m, ok = db.LookupMenu(emptyMenu.Name)
	require.True(t, ok, "Could not find empty Menu after reopening DB")
	require.Equal(t, emptyMenu, m, "Empty menu does not match the one after reopening DB")

	require.ElementsMatch(t, []types.Menu{myMenu, emptyMenu}, menus, "Menus do not match the ones after reopening DB")

	require.NoError(t, db.Close())
}

//nolint:dupl // This is a test file, so it's normal to have similar functions
func PantriesTest(t *testing.T, openDB func() database.DB) {
	t.Helper()

	db := openDB()
	defer db.Close()

	_, ok := db.LookupProduct("AAAAAAAAAAAAAAAAAAAAA")
	require.False(t, ok)

	pantry1 := types.Pantry{
		Name: "Pantry #1",
		Contents: []types.Ingredient{
			{Name: "Ingredient #1", Amount: 1.0},
			{Name: "Ingredient #2", Amount: 2.0},
		},
	}

	pantry2 := types.Pantry{
		Name:     "Pantry #2",
		Contents: nil,
	}

	require.NoError(t, db.SetPantry(pantry1), "Could not set Pantry")
	p, ok := db.LookupPantry(pantry1.Name)
	require.True(t, ok, "Could not find Pantry just created")
	require.Equal(t, pantry1, p, "Pantry does not match the one just created")

	pantry1.Contents[0].Amount = 5.0

	require.NoError(t, db.SetPantry(pantry1), "Could not override Pantry")
	p, ok = db.LookupPantry(pantry1.Name)
	require.True(t, ok, "Could not find Pantry just overridden created")
	require.Equal(t, pantry1, p, "Pantry does not match the one just overridden")

	require.NoError(t, db.SetPantry(pantry2), "Could not set Pantry")
	p, ok = db.LookupPantry(pantry2.Name)
	require.True(t, ok, "Could not find empty Pantry just created")
	require.Equal(t, pantry2, p, "Empty menu does not match the one just created")

	t.Log("Closing DB and reopening")
	require.NoError(t, db.Close())

	menus := db.Pantries()
	require.ElementsMatch(t, []types.Pantry{pantry1, pantry2}, menus, "Pantries do not match the ones just created")

	db = openDB()
	defer db.Close()

	_, ok = db.LookupPantry("AAAAAAAAAAAAAAAAAAAAA")
	require.False(t, ok)

	p, ok = db.LookupPantry(pantry1.Name)
	require.True(t, ok, "Could not find Pantry after reopening DB")
	require.Equal(t, pantry1, p, "Pantry does not match the one after reopening DB")

	p, ok = db.LookupPantry(pantry2.Name)
	require.True(t, ok, "Could not find empty Pantry after reopening DB")
	require.Equal(t, pantry2, p, "Empty menu does not match the one after reopening DB")

	require.ElementsMatch(t, []types.Pantry{pantry1, pantry2}, menus, "Pantries do not match the ones after reopening DB")

	require.NoError(t, db.Close())
}

func ShoppingListsTest(t *testing.T, openDB func() database.DB) {
	t.Helper()

	db := openDB()
	defer db.Close()

	_, ok := db.LookupProduct("AAAAAAAAAAAAAAAAAAAAA")
	require.False(t, ok)

	list1 := types.ShoppingList{
		Name:      "ShoppingList #1",
		TimeStamp: "2021-09-01T00:00:00Z",
		Items: []string{
			"Item #1",
			"Item #2",
		},
	}

	list2 := types.ShoppingList{
		Name:      "ShoppingList #2",
		TimeStamp: "2027-09-01T00:00:01Z",
		Items:     nil,
	}

	require.NoError(t, db.SetShoppingList(list1), "Could not set ShoppingList")
	p, ok := db.LookupShoppingList(list1.Name)
	require.True(t, ok, "Could not find ShoppingList just created")
	require.Equal(t, list1, p, "ShoppingList does not match the one just created")

	list1.TimeStamp = "2024-01-01T00:00:03Z"

	require.NoError(t, db.SetShoppingList(list1), "Could not override ShoppingList")
	p, ok = db.LookupShoppingList(list1.Name)
	require.True(t, ok, "Could not find ShoppingList just overridden created")
	require.Equal(t, list1, p, "ShoppingList does not match the one just overridden")

	require.NoError(t, db.SetShoppingList(list2), "Could not set ShoppingList")
	p, ok = db.LookupShoppingList(list2.Name)
	require.True(t, ok, "Could not find empty ShoppingList just created")
	require.Equal(t, list2, p, "Empty menu does not match the one just created")

	t.Log("Closing DB and reopening")
	require.NoError(t, db.Close())

	menus := db.ShoppingLists()
	require.ElementsMatch(t, []types.ShoppingList{list1, list2}, menus, "ShoppingLists do not match the ones just created")

	db = openDB()
	defer db.Close()

	_, ok = db.LookupShoppingList("AAAAAAAAAAAAAAAAAAAAA")
	require.False(t, ok)

	p, ok = db.LookupShoppingList(list1.Name)
	require.True(t, ok, "Could not find ShoppingList after reopening DB")
	require.Equal(t, list1, p, "ShoppingList does not match the one after reopening DB")

	p, ok = db.LookupShoppingList(list2.Name)
	require.True(t, ok, "Could not find empty ShoppingList after reopening DB")
	require.Equal(t, list2, p, "Empty menu does not match the one after reopening DB")

	require.ElementsMatch(t, []types.ShoppingList{list1, list2}, menus, "ShoppingLists do not match the ones after reopening DB")

	require.NoError(t, db.Close())
}
