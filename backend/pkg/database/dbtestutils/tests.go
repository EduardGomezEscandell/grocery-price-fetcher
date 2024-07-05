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
		Provider:  blank.Provider{},
		Name:      "Product #1",
		Price:     1.99,
		BatchSize: 11,
	}

	product2 := product.Product{
		Provider:  blank.Provider{},
		Name:      "Product #2",
		Price:     0.64,
		BatchSize: 99,
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

	_, err := db.LookupRecipe(1)
	require.ErrorIs(t, err, fs.ErrNotExist)

	hydrogen := product.Product{
		Provider:  blank.Provider{},
		Name:      "Hydrogen",
		ID:        53,
		BatchSize: 1,
	}

	oxygen := product.Product{
		Provider:  blank.Provider{},
		Name:      "Oxygen",
		ID:        87,
		BatchSize: 16,
	}

	_, err = db.SetProduct(hydrogen)
	require.NoError(t, err)

	_, err = db.SetProduct(oxygen)
	require.NoError(t, err)

	recipe1 := recipe.Recipe{
		Name: "Water",
		Ingredients: []recipe.Ingredient{
			{ProductID: hydrogen.ID, Amount: 2.0},
			{ProductID: oxygen.ID, Amount: 1.0},
		},
	}

	recipe2 := recipe.Recipe{
		Name:        "Empty",
		Ingredients: make([]recipe.Ingredient, 0),
	}

	rID, err := db.SetRecipe(recipe1)
	require.NoError(t, err, "Could not set Recipe")
	require.NotZero(t, rID, "ID should be non-zero")
	recipe1.ID = rID

	r, err := db.LookupRecipe(rID)
	require.NoError(t, err, "Could not find Recipe just created")
	require.Equal(t, recipe1, r, "Recipe does not match the one just created")

	recipe1.Ingredients[0].Amount = 5.0

	rID, err = db.SetRecipe(recipe1)
	require.NoError(t, err, "Could not override Recipe")
	require.Equal(t, recipe1.ID, rID, "ID should be the same after overriding")

	r, err = db.LookupRecipe(recipe1.ID)
	require.NoError(t, err, "Could not find Recipe just overridden")
	require.Equal(t, recipe1, r, "Recipe does not match the one just overridden")

	// Test implicit deletion of ingredients
	recipe1.Ingredients = recipe1.Ingredients[:1]

	_, err = db.SetRecipe(recipe1)
	require.NoError(t, err, "Could not override Recipe")
	r, err = db.LookupRecipe(recipe1.ID)
	require.NoError(t, err, "Could not find Recipe just overridden")
	require.Equal(t, recipe1, r, "Recipe does not match the one just overridden")

	rID, err = db.SetRecipe(recipe2)
	require.NoError(t, err, "Could not set Recipe")
	require.NotZero(t, rID, "ID should be non-zero")
	recipe2.ID = rID

	r, err = db.LookupRecipe(recipe2.ID)
	require.NoError(t, err, "Could not find empty Recipe just created")
	require.Equal(t, recipe2, r, "Empty menu does not match the one just created")

	t.Log("Closing DB and reopening")
	require.NoError(t, db.Close())

	db = openDB()
	defer db.Close()

	menus, err := db.Recipes()
	require.NoError(t, err)
	require.ElementsMatch(t, []recipe.Recipe{recipe1, recipe2}, menus, "Recipes do not match the ones just created")

	_, err = db.LookupRecipe(999999)
	require.ErrorIs(t, err, fs.ErrNotExist)

	r, err = db.LookupRecipe(recipe1.ID)
	require.NoError(t, err, "Could not find Recipe after reopening DB")
	require.Equal(t, recipe1, r, "Recipe does not match the one after reopening DB")

	r, err = db.LookupRecipe(recipe2.ID)
	require.NoError(t, err, "Could not find empty Recipe after reopening DB")
	require.Equal(t, recipe2, r, "Empty menu does not match the one after reopening DB")

	require.ElementsMatch(t, []recipe.Recipe{recipe1, recipe2}, menus, "Recipes do not match the ones after reopening DB")

	err = db.DeleteRecipe(recipe1.ID)
	require.NoError(t, err)

	_, err = db.LookupRecipe(recipe1.ID)
	require.ErrorIs(t, err, fs.ErrNotExist)

	require.NoError(t, db.Close())
}

func MenuTest(t *testing.T, openDB func() database.DB) {
	t.Helper()

	db := openDB()
	defer db.Close()

	const user = "default"

	_, err := db.LookupMenu("default", "AAAAAAAAAAAAAAAAAAAAA")
	require.Error(t, err)

	recipeID, err := db.SetRecipe(recipe.Recipe{
		Name:        "Empty recipe",
		Ingredients: []recipe.Ingredient{},
	})
	require.NoError(t, err)

	myMenu := dbtypes.Menu{
		User: user,
		Name: "myMenu",
		Days: []dbtypes.Day{
			{
				Name: "Segunda-Feira",
				Meals: []dbtypes.Meal{
					{
						Name: "Café da Manhã",
						Dishes: []dbtypes.Dish{
							{
								ID:     recipeID,
								Amount: 16,
							},
						},
					},
				},
			},
		},
	}

	require.NoError(t, db.SetMenu(myMenu), "Could not set Menu")
	m, err := db.LookupMenu(user, myMenu.Name)
	require.NoError(t, err, "Could not find Menu just created")
	require.Equal(t, myMenu, m, "Menu does not match the one just created")

	_, err = db.LookupMenu(anotherUser, myMenu.Name)
	require.ErrorIs(t, err, fs.ErrNotExist, "Should not find Menu from another user")

	myMenu.Days[0].Meals[0].Dishes[0].Amount = 20

	require.NoError(t, db.SetMenu(myMenu), "Could not override Menu")
	m, err = db.LookupMenu(user, myMenu.Name)
	require.NoError(t, err, "Could not find Menu just overridden")
	require.Equal(t, myMenu, m, "Menu does not match the one just overridden")

	// Test implicit deletion of dishes
	myMenu.Days[0].Meals[0].Dishes = nil

	require.NoError(t, db.SetMenu(myMenu), "Could not override Menu")
	m, err = db.LookupMenu(user, myMenu.Name)
	require.NoError(t, err, "Could not find Menu just overridden")
	require.Equal(t, myMenu, m, "Menu does not match the one just overridden")

	emptyMenu := dbtypes.Menu{
		User: user,
		Name: "Empty Menu",
	}

	require.NoError(t, db.SetMenu(emptyMenu), "Could not set Menu")
	m, err = db.LookupMenu(user, emptyMenu.Name)
	require.NoError(t, err, "Could not find empty Menu just created")
	require.Equal(t, emptyMenu, m, "Empty menu does not match the one just created")

	t.Log("Closing DB and reopening")
	require.NoError(t, db.Close())

	db = openDB()
	defer db.Close()

	menus, err := db.Menus(user)
	require.NoError(t, err)
	require.ElementsMatch(t, []dbtypes.Menu{myMenu, emptyMenu}, menus, "Menus do not match the ones just created")

	_, err = db.LookupMenu(user, "AAAAAAAAAAAAAAAAAAAAA")
	require.Error(t, err)

	m, err = db.LookupMenu(user, myMenu.Name)
	require.NoError(t, err, "Could not find Menu after reopening DB")
	require.Equal(t, myMenu, m, "Menu does not match the one after reopening DB")

	m, err = db.LookupMenu(user, emptyMenu.Name)
	require.NoError(t, err, "Could not find empty Menu after reopening DB")
	require.Equal(t, emptyMenu, m, "Empty menu does not match the one after reopening DB")

	require.ElementsMatch(t, []dbtypes.Menu{myMenu, emptyMenu}, menus, "Menus do not match the ones after reopening DB")

	err = db.DeleteMenu(user, myMenu.Name)
	require.NoError(t, err)

	_, err = db.LookupMenu(user, myMenu.Name)
	require.Error(t, err)

	require.NoError(t, db.Close())
}

func PantriesTest(t *testing.T, openDB func() database.DB) {
	t.Helper()

	db := openDB()
	defer db.Close()

	_, ok := db.LookupPantry("AAAAAAAAAAAAAAAAAAAAA")
	require.False(t, ok)

	hydrogen := product.Product{
		Provider:  blank.Provider{},
		Name:      "Hydrogen",
		ID:        53,
		BatchSize: 1,
	}

	oxygen := product.Product{
		Provider:  blank.Provider{},
		Name:      "Oxygen",
		ID:        87,
		BatchSize: 16,
	}

	_, err := db.SetProduct(hydrogen)
	require.NoError(t, err)

	_, err = db.SetProduct(oxygen)
	require.NoError(t, err)

	pantry1 := dbtypes.Pantry{
		Name: "Pantry #1",
		Contents: []recipe.Ingredient{
			{ProductID: hydrogen.ID, Amount: 2.0},
			{ProductID: oxygen.ID, Amount: 1.0},
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

	hydrogen := product.Product{
		Provider:  blank.Provider{},
		Name:      "Hydrogen",
		BatchSize: 1,
	}

	oxygen := product.Product{
		Provider:  blank.Provider{},
		Name:      "Oxygen",
		BatchSize: 16,
	}

	id, err := db.SetProduct(hydrogen)
	require.NoError(t, err)
	hydrogen.ID = id

	id, err = db.SetProduct(oxygen)
	require.NoError(t, err)
	oxygen.ID = id

	menu1 := dbtypes.Menu{
		Name: "Menu #1",
	}

	menu2 := dbtypes.Menu{
		Name: "Menu #99",
	}

	err = db.SetMenu(menu1)
	require.NoError(t, err)

	err = db.SetMenu(menu2)
	require.NoError(t, err)

	pantry := dbtypes.Pantry{
		Name: "Pantry #1",
	}

	err = db.SetPantry(pantry)
	require.NoError(t, err)

	list1 := dbtypes.ShoppingList{
		Menu:     menu1.Name,
		Pantry:   pantry.Name,
		Contents: []product.ID{oxygen.ID},
	}

	list2 := dbtypes.ShoppingList{
		Menu:     menu2.Name,
		Pantry:   pantry.Name,
		Contents: []product.ID{hydrogen.ID},
	}

	require.NoError(t, db.SetShoppingList(list1), "Could not set ShoppingList")
	p, ok := db.LookupShoppingList(list1.Menu, list1.Pantry)
	require.True(t, ok, "Could not find ShoppingList just created")
	// Sort the slices to make sure the order is deterministic
	slices.Sort(list1.Contents)
	slices.Sort(p.Contents)
	require.Equal(t, list1, p, "ShoppingList does not match the one just created")

	list1.Contents[0] = oxygen.ID

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
