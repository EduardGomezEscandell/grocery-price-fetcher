package dbtestutils

import (
	"io/fs"
	"slices"
	"testing"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/database"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/database/dbtypes"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/product"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/providers/blank"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/recipe"
	"github.com/stretchr/testify/require"
)

func ProductsTest(t *testing.T, openDB func() database.DB) {
	t.Helper()

	db := openDB()
	defer db.Close()

	_, err := db.LookupProduct(1)
	require.ErrorIs(t, err, fs.ErrNotExist)

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

	id, err := db.SetProduct(product1)
	require.NoError(t, err, "Could not set Product")
	require.NotZero(t, id, "ID should be non-zero")

	product1.ID = id

	p, err := db.LookupProduct(product1.ID)
	require.NoError(t, err, "Could not find Product just created")
	require.Equal(t, product1, p, "Product does not match the one just created")

	product1.BatchSize = 20

	id, err = db.SetProduct(product1)
	require.NoError(t, err, "Could not override Product")
	require.Equal(t, product1.ID, id, "ID should be the same after overriding")

	p, err = db.LookupProduct(id)
	require.NoError(t, err, "Could not find Product just overridden")
	require.Equal(t, product1, p, "Product does not match the one just overridden")

	id, err = db.SetProduct(product2)
	require.NoError(t, err, "Could not set Product")
	require.NotZero(t, p, "ID should be non-zero")
	product2.ID = id

	p, err = db.LookupProduct(product2.ID)
	require.NoError(t, err, "Could not find empty Product just created")
	require.Equal(t, product2, p, "Empty menu does not match the one just created")

	t.Log("Closing DB and reopening")
	require.NoError(t, db.Close())

	db = openDB()
	defer db.Close()

	products, err := db.Products()
	require.NoError(t, err)
	require.ElementsMatch(t, []product.Product{product1, product2}, products, "Products do not match the ones just created")

	// Neither of them can be zero, so this is be safe
	newID := product1.ID + product2.ID + 1

	_, err = db.LookupProduct(newID)
	require.ErrorIs(t, err, fs.ErrNotExist)

	p, err = db.LookupProduct(product1.ID)
	require.NoError(t, err, "Could not find Product after reopening DB")
	require.Equal(t, product1, p, "Product does not match the one after reopening DB")

	p, err = db.LookupProduct(product2.ID)
	require.NoError(t, err, "Could not find empty Product after reopening DB")
	require.Equal(t, product2, p, "Empty menu does not match the one after reopening DB")

	require.ElementsMatch(t, []product.Product{product1, product2}, products, "Products do not match the ones after reopening DB")

	err = db.DeleteProduct(product1.ID)
	require.NoError(t, err)

	_, err = db.LookupProduct(product1.ID)
	require.ErrorIs(t, err, fs.ErrNotExist)

	require.NoError(t, db.Close())
}

func RecipesTest(t *testing.T, openDB func() database.DB) {
	t.Helper()

	db := openDB()
	defer db.Close()

	_, ok := db.LookupRecipe("AAAAAAAAAAAAAAAAAAAAA")
	require.False(t, ok)

	recipe1 := recipe.Recipe{
		Name: "Recipe #1",
		Ingredients: []recipe.Ingredient{
			{ProductID: 53, Amount: 1.0},
			{ProductID: 87, Amount: 2.0},
		},
	}

	recipe2 := recipe.Recipe{
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
	require.True(t, ok, "Could not find Recipe just overridden")
	require.Equal(t, recipe1, p, "Recipe does not match the one just overridden")

	// Test implicit deletion of ingredients
	recipe1.Ingredients = recipe1.Ingredients[:1]
	require.NoError(t, db.SetRecipe(recipe1), "Could not override Recipe")
	p, ok = db.LookupRecipe(recipe1.Name)
	require.True(t, ok, "Could not find Recipe just overridden")
	require.Equal(t, recipe1, p, "Recipe does not match the one just overridden")

	require.NoError(t, db.SetRecipe(recipe2), "Could not set Recipe")
	p, ok = db.LookupRecipe(recipe2.Name)
	require.True(t, ok, "Could not find empty Recipe just created")
	require.Equal(t, recipe2, p, "Empty menu does not match the one just created")

	t.Log("Closing DB and reopening")
	require.NoError(t, db.Close())

	db = openDB()
	defer db.Close()

	menus, err := db.Recipes()
	require.NoError(t, err)
	require.ElementsMatch(t, []recipe.Recipe{recipe1, recipe2}, menus, "Recipes do not match the ones just created")

	_, ok = db.LookupRecipe("AAAAAAAAAAAAAAAAAAAAA")
	require.False(t, ok)

	p, ok = db.LookupRecipe(recipe1.Name)
	require.True(t, ok, "Could not find Recipe after reopening DB")
	require.Equal(t, recipe1, p, "Recipe does not match the one after reopening DB")

	p, ok = db.LookupRecipe(recipe2.Name)
	require.True(t, ok, "Could not find empty Recipe after reopening DB")
	require.Equal(t, recipe2, p, "Empty menu does not match the one after reopening DB")

	require.ElementsMatch(t, []recipe.Recipe{recipe1, recipe2}, menus, "Recipes do not match the ones after reopening DB")

	err = db.DeleteRecipe(recipe1.Name)
	require.NoError(t, err)

	_, ok = db.LookupRecipe(recipe1.Name)
	require.False(t, ok)

	require.NoError(t, db.Close())
}

func MenuTest(t *testing.T, openDB func() database.DB) {
	t.Helper()

	db := openDB()
	defer db.Close()

	_, ok := db.LookupMenu("AAAAAAAAAAAAAAAAAAAAA")
	require.False(t, ok)

	myMenu := dbtypes.Menu{
		Name: "myMenu",
		Days: []dbtypes.Day{
			{
				Name: "Segiunda-Feira",
				Meals: []dbtypes.Meal{
					{
						Name: "Café da Manhã",
						Dishes: []dbtypes.Dish{
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
	require.True(t, ok, "Could not find Menu just overridden")
	require.Equal(t, myMenu, m, "Menu does not match the one just overridden")

	// Test implicit deletion of dishes
	myMenu.Days[0].Meals[0].Dishes = nil

	require.NoError(t, db.SetMenu(myMenu), "Could not override Menu")
	m, ok = db.LookupMenu(myMenu.Name)
	require.True(t, ok, "Could not find Menu just overridden")
	require.Equal(t, myMenu, m, "Menu does not match the one just overridden")

	emptyMenu := dbtypes.Menu{
		Name: "Empty Menu",
	}

	require.NoError(t, db.SetMenu(emptyMenu), "Could not set Menu")
	m, ok = db.LookupMenu(emptyMenu.Name)
	require.True(t, ok, "Could not find empty Menu just created")
	require.Equal(t, emptyMenu, m, "Empty menu does not match the one just created")

	t.Log("Closing DB and reopening")
	require.NoError(t, db.Close())

	db = openDB()
	defer db.Close()

	menus, err := db.Menus()
	require.NoError(t, err)
	require.ElementsMatch(t, []dbtypes.Menu{myMenu, emptyMenu}, menus, "Menus do not match the ones just created")

	_, ok = db.LookupMenu("AAAAAAAAAAAAAAAAAAAAA")
	require.False(t, ok)

	m, ok = db.LookupMenu(myMenu.Name)
	require.True(t, ok, "Could not find Menu after reopening DB")
	require.Equal(t, myMenu, m, "Menu does not match the one after reopening DB")

	m, ok = db.LookupMenu(emptyMenu.Name)
	require.True(t, ok, "Could not find empty Menu after reopening DB")
	require.Equal(t, emptyMenu, m, "Empty menu does not match the one after reopening DB")

	require.ElementsMatch(t, []dbtypes.Menu{myMenu, emptyMenu}, menus, "Menus do not match the ones after reopening DB")

	err = db.DeleteMenu(myMenu.Name)
	require.NoError(t, err)

	_, ok = db.LookupMenu(myMenu.Name)
	require.False(t, ok)

	require.NoError(t, db.Close())
}

func PantriesTest(t *testing.T, openDB func() database.DB) {
	t.Helper()

	db := openDB()
	defer db.Close()

	_, ok := db.LookupPantry("AAAAAAAAAAAAAAAAAAAAA")
	require.False(t, ok)

	pantry1 := dbtypes.Pantry{
		Name: "Pantry #1",
		Contents: []recipe.Ingredient{
			{ProductID: 33, Amount: 1.0},
			{ProductID: 66, Amount: 2.0},
		},
	}

	pantry2 := dbtypes.Pantry{
		Name:     "Pantry #2",
		Contents: []recipe.Ingredient{},
	}

	require.NoError(t, db.SetPantry(pantry1), "Could not set Pantry")
	p, ok := db.LookupPantry(pantry1.Name)
	require.True(t, ok, "Could not find Pantry just created")
	require.Equal(t, pantry1, p, "Pantry does not match the one just created")

	pantry1.Contents[0].Amount = 5.0

	require.NoError(t, db.SetPantry(pantry1), "Could not override Pantry")
	p, ok = db.LookupPantry(pantry1.Name)
	require.True(t, ok, "Could not find Pantry just overridden")
	require.Equal(t, pantry1, p, "Pantry does not match the one just overridden")

	// Test implicit deletion of ingredients
	pantry1.Contents = pantry1.Contents[:1]
	require.NoError(t, db.SetPantry(pantry1), "Could not override Pantry")
	p, ok = db.LookupPantry(pantry1.Name)
	require.True(t, ok, "Could not find Pantry just overridden")
	require.Equal(t, pantry1, p, "Pantry does not match the one just overridden")

	require.NoError(t, db.SetPantry(pantry2), "Could not set Pantry")
	p, ok = db.LookupPantry(pantry2.Name)
	require.True(t, ok, "Could not find empty Pantry just created")
	require.Equal(t, pantry2, p, "Empty menu does not match the one just created")

	t.Log("Closing DB and reopening")
	require.NoError(t, db.Close())

	db = openDB()
	defer db.Close()

	menus, err := db.Pantries()
	require.NoError(t, err)
	require.ElementsMatch(t, []dbtypes.Pantry{pantry1, pantry2}, menus, "Pantries do not match the ones just created")

	_, ok = db.LookupPantry("AAAAAAAAAAAAAAAAAAAAA")
	require.False(t, ok)

	p, ok = db.LookupPantry(pantry1.Name)
	require.True(t, ok, "Could not find Pantry after reopening DB")
	require.Equal(t, pantry1, p, "Pantry does not match the one after reopening DB")

	p, ok = db.LookupPantry(pantry2.Name)
	require.True(t, ok, "Could not find empty Pantry after reopening DB")
	require.Equal(t, pantry2, p, "Empty menu does not match the one after reopening DB")

	require.ElementsMatch(t, []dbtypes.Pantry{pantry1, pantry2}, menus, "Pantries do not match the ones after reopening DB")

	err = db.DeletePantry(pantry1.Name)
	require.NoError(t, err)

	_, ok = db.LookupPantry(pantry1.Name)
	require.False(t, ok)

	require.NoError(t, db.Close())
}

func ShoppingListsTest(t *testing.T, openDB func() database.DB) {
	t.Helper()

	db := openDB()
	defer db.Close()

	_, ok := db.LookupShoppingList("AAAAAAAAAAAAAAAAAAAAA", "AAAAAAAAAAAAaaaa")
	require.False(t, ok)

	list1 := dbtypes.ShoppingList{
		Menu:     "Menu #1",
		Pantry:   "Pantry #1",
		Contents: []product.ID{13, 99},
	}

	list2 := dbtypes.ShoppingList{
		Menu:     "Menu #99",
		Pantry:   "Pantry #1",
		Contents: []product.ID{500, 10000},
	}

	require.NoError(t, db.SetShoppingList(list1), "Could not set ShoppingList")
	p, ok := db.LookupShoppingList(list1.Menu, list1.Pantry)
	require.True(t, ok, "Could not find ShoppingList just created")
	// Sort the slices to make sure the order is deterministic
	slices.Sort(list1.Contents)
	slices.Sort(p.Contents)
	require.Equal(t, list1, p, "ShoppingList does not match the one just created")

	list1.Contents[0] = 111111

	require.NoError(t, db.SetShoppingList(list1), "Could not override ShoppingList")
	p, ok = db.LookupShoppingList(list1.Menu, list1.Pantry)
	require.True(t, ok, "Could not find ShoppingList just overridden")
	// Sort the slices to make sure the order is deterministic
	slices.Sort(list1.Contents)
	slices.Sort(p.Contents)
	require.Equal(t, list1, p, "ShoppingList does not match the one just overridden")

	// Test implicit deletion of items
	list1.Contents = list1.Contents[:1]
	require.NoError(t, db.SetShoppingList(list1), "Could not override ShoppingList")
	p, ok = db.LookupShoppingList(list1.Menu, list1.Pantry)
	require.True(t, ok, "Could not find ShoppingList just overridden")
	// Sort the slices to make sure the order is deterministic
	slices.Sort(list1.Contents)
	slices.Sort(p.Contents)
	require.Equal(t, list1, p, "ShoppingList does not match the one just overridden")

	require.NoError(t, db.SetShoppingList(list2), "Could not set ShoppingList")
	p, ok = db.LookupShoppingList(list2.Menu, list2.Pantry)
	require.True(t, ok, "Could not find empty ShoppingList just created")
	// Sort the slices to make sure the order is deterministic
	slices.Sort(list1.Contents)
	slices.Sort(p.Contents)
	require.Equal(t, list2, p, "Empty menu does not match the one just created")

	t.Log("Closing DB and reopening")
	require.NoError(t, db.Close())

	db = openDB()
	defer db.Close()

	menus, err := db.ShoppingLists()
	require.NoError(t, err)
	require.ElementsMatch(t, []dbtypes.ShoppingList{list1, list2}, menus, "ShoppingLists do not match the ones just created")

	_, ok = db.LookupShoppingList("FAKE MENU", "FAKE PANTRY")
	require.False(t, ok)

	p, ok = db.LookupShoppingList(list1.Menu, list1.Pantry)
	require.True(t, ok, "Could not find ShoppingList after reopening DB")
	// Sort the slices to make sure the order is deterministic
	slices.Sort(list1.Contents)
	slices.Sort(p.Contents)
	require.Equal(t, list1, p, "ShoppingList does not match the one after reopening DB")

	p, ok = db.LookupShoppingList(list2.Menu, list2.Pantry)
	require.True(t, ok, "Could not find empty ShoppingList after reopening DB")
	// Sort the slices to make sure the order is deterministic
	slices.Sort(list1.Contents)
	slices.Sort(p.Contents)
	require.Equal(t, list2, p, "Empty menu does not match the one after reopening DB")

	require.ElementsMatch(t, []dbtypes.ShoppingList{list1, list2}, menus, "ShoppingLists do not match the ones after reopening DB")

	err = db.DeleteShoppingList(list1.Menu, list1.Pantry)
	require.NoError(t, err)

	_, ok = db.LookupShoppingList(list1.Menu, list1.Pantry)
	require.False(t, ok)

	require.NoError(t, db.Close())
}
