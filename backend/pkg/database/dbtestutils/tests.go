package dbtestutils

import (
	"io/fs"
	"slices"
	"testing"
	"time"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/database"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/database/dbtypes"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/product"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/providers/blank"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/recipe"
	"github.com/stretchr/testify/require"
)

func UsersTest(t *testing.T, openDB func() database.DB) {
	t.Helper()

	db := openDB()
	defer db.Close()

	user1 := "user1"
	user2 := "user2"

	exists, err := db.LookupUser(user1)
	require.NoError(t, err)
	require.False(t, exists, "User should not exist")

	require.NoError(t, db.SetUser(user1), "Could not set User")
	exists, err = db.LookupUser(user1)
	require.NoError(t, err)
	require.True(t, exists, "User should exist")

	require.NoError(t, db.SetUser(user2), "Could not set User")
	exists, err = db.LookupUser(user2)
	require.NoError(t, err)
	require.True(t, exists, "User should exist")

	require.NoError(t, db.DeleteUser(user1), "Could not delete User")
	exists, err = db.LookupUser(user1)
	require.NoError(t, err)
	require.False(t, exists, "User should not exist")

	t.Log("Closing DB and reopening")
	require.NoError(t, db.Close())

	db = openDB()
	defer db.Close()

	exists, err = db.LookupUser(user1)
	require.NoError(t, err)
	require.False(t, exists, "User should not exist after reopening DB")

	exists, err = db.LookupUser(user2)
	require.NoError(t, err)
	require.False(t, exists, "User should not exist after reopening DB")

	recipeID, err := db.SetRecipe(recipe.Recipe{User: user2, Name: "Recipe"})
	require.NoError(t, err, "Could not set Recipe")
	require.NoError(t, db.SetPantry(dbtypes.Pantry{User: user2, Name: "Pantry"}), "Could not set Pantry")
	require.NoError(t, db.SetMenu(dbtypes.Menu{User: user2, Name: "Menu"}), "Could not set Menu")
	require.NoError(t, db.SetShoppingList(dbtypes.ShoppingList{User: user2, Menu: "Menu", Pantry: "Pantry"}), "Could not set ShoppingList")

	require.NoError(t, db.DeleteUser(user2), "Could not delete User")
	exists, err = db.LookupUser(user2)
	require.NoError(t, err)
	require.False(t, exists, "User should not exist")

	_, err = db.LookupRecipe(user2, recipeID)
	require.ErrorIs(t, err, fs.ErrNotExist, "Recipe should not exist after deleting User")
	_, err = db.LookupPantry(user2, "Pantry")
	require.ErrorIs(t, err, fs.ErrNotExist, "Pantry should not exist after deleting User")
	_, err = db.LookupMenu(user2, "Menu")
	require.ErrorIs(t, err, fs.ErrNotExist, "Menu should not exist after deleting User")
	_, err = db.LookupShoppingList(user2, "Menu", "Pantry")
	require.ErrorIs(t, err, fs.ErrNotExist, "ShoppingList should not exist after deleting User")
}

func SessionsTest(t *testing.T, openDB func() database.DB) {
	t.Helper()

	db := openDB()
	defer db.Close()

	const user = "test-user-123"
	const anotherUser = "another-user-456"

	_, err := db.LookupSession("AAAAAAAAAAAAAAAAAAAAA")
	require.ErrorIs(t, err, fs.ErrNotExist)

	session := dbtypes.Session{
		ID:           "session123",
		User:         user,
		AccessToken:  "token123",
		RefreshToken: "refresh123",
		NotAfter:     time.Now().Add(time.Hour),
	}

	otherSession := dbtypes.Session{
		ID:           "session456",
		User:         user,
		AccessToken:  "token456",
		RefreshToken: "refresh456",
		NotAfter:     time.Now().Add(10 * time.Hour),
	}

	err = db.SetSession(session)
	require.Error(t, err, "Should not be able to set a session for a non-existent user")

	err = db.SetUser(user)
	require.NoError(t, err)

	err = db.SetUser(anotherUser)
	require.NoError(t, err)

	err = db.SetSession(session)
	require.NoError(t, err, "Could not set a new session")

	err = db.SetSession(otherSession)
	require.NoError(t, err, "Could not set a new session")

	s, err := db.LookupSession(session.ID)
	require.NoError(t, err, "Could not find Session just created")
	require.Equal(t, session, s, "Session does not match the one just created")

	session.NotAfter = time.Now().Add(10 * time.Hour)
	err = db.SetSession(session)
	require.NoError(t, err, "Could not re-set a session")

	s, err = db.LookupSession(session.ID)
	require.NoError(t, err, "Could not find Session just re-set")
	require.Equal(t, session, s, "Session does not match the one just re-set")

	session.User = anotherUser
	err = db.SetSession(session)
	require.Error(t, err, "Should not be able to reuse a session for a different user")

	err = db.DeleteSession(session.ID)
	require.NoError(t, err)

	_, err = db.LookupSession(session.ID)
	require.ErrorIs(t, err, fs.ErrNotExist, "Should not find Session after deleting it")

	session.User = user
	session.NotAfter = time.Now().Add(-time.Hour)
	err = db.SetSession(session)
	require.Error(t, err, "Should not be able to set an expired session")

	session.NotAfter = time.Now().Add(5 * time.Second)
	err = db.SetSession(session)
	require.NoError(t, err, "Could not set a new session")

	time.Sleep(10 * time.Second)

	err = db.PurgeSessions()
	require.NoError(t, err)

	_, err = db.LookupSession(session.ID)
	require.ErrorIs(t, err, fs.ErrNotExist, "Should not find expired Session after purging the database")

	_, err = db.LookupSession(otherSession.ID)
	require.NoError(t, err, "Should find non-expired Session after purging the database")
}

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

	const user = "test-user-123"
	const anotherUser = "another-user-456"

	_, err := db.LookupRecipe(user, 123)
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

	err = db.SetUser(user)
	require.NoError(t, err)

	_, err = db.SetProduct(hydrogen)
	require.NoError(t, err)

	_, err = db.SetProduct(oxygen)
	require.NoError(t, err)

	recipe1 := recipe.Recipe{
		User: user,
		Name: "Water",
		Ingredients: []recipe.Ingredient{
			{ProductID: hydrogen.ID, Amount: 2.0},
			{ProductID: oxygen.ID, Amount: 1.0},
		},
	}

	recipe2 := recipe.Recipe{
		User:        user,
		Name:        "Empty",
		Ingredients: make([]recipe.Ingredient, 0),
	}

	rID, err := db.SetRecipe(recipe1)
	require.NoError(t, err, "Could not set Recipe")
	require.NotZero(t, rID, "ID should be non-zero")
	recipe1.ID = rID

	r, err := db.LookupRecipe(user, rID)
	require.NoError(t, err, "Could not find Recipe just created")
	require.Equal(t, recipe1, r, "Recipe does not match the one just created")

	_, err = db.LookupRecipe(anotherUser, rID)
	require.ErrorIs(t, err, fs.ErrNotExist, "Should not find Recipe from another user")

	recipe1.Ingredients[0].Amount = 5.0

	rID, err = db.SetRecipe(recipe1)
	require.NoError(t, err, "Could not override Recipe")
	require.Equal(t, recipe1.ID, rID, "ID should be the same after overriding")

	r, err = db.LookupRecipe(user, recipe1.ID)
	require.NoError(t, err, "Could not find Recipe just overridden")
	require.Equal(t, recipe1, r, "Recipe does not match the one just overridden")

	// Test implicit deletion of ingredients
	recipe1.Ingredients = recipe1.Ingredients[:1]

	_, err = db.SetRecipe(recipe1)
	require.NoError(t, err, "Could not override Recipe")
	r, err = db.LookupRecipe(user, recipe1.ID)
	require.NoError(t, err, "Could not find Recipe just overridden")
	require.Equal(t, recipe1, r, "Recipe does not match the one just overridden")

	rID, err = db.SetRecipe(recipe2)
	require.NoError(t, err, "Could not set Recipe")
	require.NotZero(t, rID, "ID should be non-zero")
	recipe2.ID = rID

	r, err = db.LookupRecipe(user, recipe2.ID)
	require.NoError(t, err, "Could not find empty Recipe just created")
	require.Equal(t, recipe2, r, "Empty menu does not match the one just created")

	recipes, err := db.Recipes(user)
	require.NoError(t, err)
	require.ElementsMatch(t, []recipe.Recipe{recipe1, recipe2}, recipes, "Recipes do not match the ones just created")

	recipes, err = db.Recipes(anotherUser)
	require.NoError(t, err)
	require.Empty(t, recipes, "Should not find Recipes from another user")

	t.Log("Closing DB and reopening")
	require.NoError(t, db.Close())

	db = openDB()
	defer db.Close()

	recipes, err = db.Recipes(user)
	require.NoError(t, err)
	require.ElementsMatch(t, []recipe.Recipe{recipe1, recipe2}, recipes, "Recipes do not match the ones after reopening DB")

	_, err = db.LookupRecipe(user, 999999)
	require.ErrorIs(t, err, fs.ErrNotExist)

	r, err = db.LookupRecipe(user, recipe1.ID)
	require.NoError(t, err, "Could not find Recipe after reopening DB")
	require.Equal(t, recipe1, r, "Recipe does not match the one after reopening DB")

	r, err = db.LookupRecipe(user, recipe2.ID)
	require.NoError(t, err, "Could not find empty Recipe after reopening DB")
	require.Equal(t, recipe2, r, "Empty menu does not match the one after reopening DB")

	err = db.DeleteRecipe(user, recipe1.ID)
	require.NoError(t, err)

	_, err = db.LookupRecipe(user, recipe1.ID)
	require.ErrorIs(t, err, fs.ErrNotExist)

	err = db.DeleteRecipe(anotherUser, recipe2.ID)
	require.NoError(t, err)
	_, err = db.LookupRecipe(user, recipe2.ID)
	require.NoError(t, err, "Should not delete Recipe from another user")

	require.NoError(t, db.Close())
}

func MenuTest(t *testing.T, openDB func() database.DB) {
	t.Helper()

	db := openDB()
	defer db.Close()

	const user = "test-user-123"
	const anotherUser = "another-user-456"

	_, err := db.LookupMenu(user, "AAAAAAAAAAAAAAAAAAAAA")
	require.Error(t, err)

	err = db.SetUser(user)
	require.NoError(t, err)

	recipeID, err := db.SetRecipe(recipe.Recipe{
		User:        user,
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

	const user = "test-user-123"
	const anotherUser = "another-user-456"

	_, err := db.LookupPantry(user, "AAAAAAAAAA")
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

	err = db.SetUser(user)
	require.NoError(t, err)

	_, err = db.SetProduct(hydrogen)
	require.NoError(t, err)

	_, err = db.SetProduct(oxygen)
	require.NoError(t, err)

	pantry1 := dbtypes.Pantry{
		User: user,
		Name: "Pantry #1",
		Contents: []recipe.Ingredient{
			{ProductID: hydrogen.ID, Amount: 2.0},
			{ProductID: oxygen.ID, Amount: 1.0},
		},
	}

	pantry2 := dbtypes.Pantry{
		User:     user,
		Name:     "Pantry #2",
		Contents: []recipe.Ingredient{},
	}

	require.NoError(t, db.SetPantry(pantry1), "Could not set Pantry")
	p, err := db.LookupPantry(user, pantry1.Name)
	require.NoError(t, err, "Could not find Pantry just created")
	require.Equal(t, pantry1, p, "Pantry does not match the one just created")

	_, err = db.LookupPantry(anotherUser, pantry1.Name)
	require.ErrorIs(t, err, fs.ErrNotExist, "Should not find Pantry from another user")

	pantry1.Contents[0].Amount = 5.0

	require.NoError(t, db.SetPantry(pantry1), "Could not override Pantry")
	p, err = db.LookupPantry(user, pantry1.Name)
	require.NoError(t, err, "Could not find Pantry just overridden")
	require.Equal(t, pantry1, p, "Pantry does not match the one just overridden")

	// Test implicit deletion of ingredients
	pantry1.Contents = pantry1.Contents[:1]
	require.NoError(t, db.SetPantry(pantry1), "Could not override Pantry")
	p, err = db.LookupPantry(user, pantry1.Name)
	require.NoError(t, err, "Could not find Pantry just overridden")
	require.Equal(t, pantry1, p, "Pantry does not match the one just overridden")

	require.NoError(t, db.SetPantry(pantry2), "Could not set Pantry")
	p, err = db.LookupPantry(user, pantry2.Name)
	require.NoError(t, err, "Could not find empty Pantry just created")
	require.Equal(t, pantry2, p, "Empty menu does not match the one just created")

	pantries, err := db.Pantries(user)
	require.NoError(t, err)
	require.ElementsMatch(t, []dbtypes.Pantry{pantry1, pantry2}, pantries, "Pantries do not match the ones just created")

	t.Log("Closing DB and reopening")
	require.NoError(t, db.Close())

	db = openDB()
	defer db.Close()

	pantries, err = db.Pantries(user)
	require.NoError(t, err)
	require.ElementsMatch(t, []dbtypes.Pantry{pantry1, pantry2}, pantries, "Pantries do not match the ones after reopening DB")

	pantries, err = db.Pantries(anotherUser)
	require.NoError(t, err)
	require.Empty(t, pantries, "Should not find Pantries from another user")

	_, err = db.LookupPantry(user, "AAAAAAAAAAAAAAAAAAAAA")
	require.ErrorIs(t, err, fs.ErrNotExist)

	p, err = db.LookupPantry(user, pantry1.Name)
	require.NoError(t, err, "Could not find Pantry after reopening DB")
	require.Equal(t, pantry1, p, "Pantry does not match the one after reopening DB")

	p, err = db.LookupPantry(user, pantry2.Name)
	require.NoError(t, err, "Could not find empty Pantry after reopening DB")
	require.Equal(t, pantry2, p, "Empty menu does not match the one after reopening DB")

	err = db.DeletePantry(user, pantry1.Name)
	require.NoError(t, err)

	_, err = db.LookupPantry(user, pantry1.Name)
	require.ErrorIs(t, err, fs.ErrNotExist)

	require.NoError(t, db.Close())
}

func ShoppingListsTest(t *testing.T, openDB func() database.DB) {
	t.Helper()

	db := openDB()
	defer db.Close()

	const user = "test-user-123"
	const anotherUser = "another-user-456"

	_, err := db.LookupShoppingList(user, "AAAAAAAAAA", "AAAAAAAAAAAA")
	require.ErrorIs(t, err, fs.ErrNotExist)

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

	err = db.SetUser(user)
	require.NoError(t, err)

	id, err := db.SetProduct(hydrogen)
	require.NoError(t, err)
	hydrogen.ID = id

	id, err = db.SetProduct(oxygen)
	require.NoError(t, err)
	oxygen.ID = id

	menu1 := dbtypes.Menu{
		User: user,
		Name: "Menu #1",
	}

	menu2 := dbtypes.Menu{
		User: user,
		Name: "Menu #99",
	}

	err = db.SetMenu(menu1)
	require.NoError(t, err)

	err = db.SetMenu(menu2)
	require.NoError(t, err)

	pantry := dbtypes.Pantry{
		User: user,
		Name: "Pantry #1",
	}

	err = db.SetPantry(pantry)
	require.NoError(t, err)

	list1 := dbtypes.ShoppingList{
		User:     user,
		Menu:     menu1.Name,
		Pantry:   pantry.Name,
		Contents: []product.ID{oxygen.ID},
	}

	list2 := dbtypes.ShoppingList{
		User:     user,
		Menu:     menu2.Name,
		Pantry:   pantry.Name,
		Contents: []product.ID{hydrogen.ID},
	}

	require.NoError(t, db.SetShoppingList(list1), "Could not set ShoppingList")
	p, err := db.LookupShoppingList(user, list1.Menu, list1.Pantry)
	require.NoError(t, err, "Could not find ShoppingList just created")
	// Sort the slices to make sure the order is deterministic
	slices.Sort(list1.Contents)
	slices.Sort(p.Contents)
	require.Equal(t, list1, p, "ShoppingList does not match the one just created")

	p, err = db.LookupShoppingList(anotherUser, list1.Menu, list1.Pantry)
	require.ErrorIs(t, err, fs.ErrNotExist, "Should not find ShoppingList from another user")

	list1.Contents[0] = oxygen.ID

	require.NoError(t, db.SetShoppingList(list1), "Could not override ShoppingList")
	p, err = db.LookupShoppingList(user, list1.Menu, list1.Pantry)
	require.NoError(t, err, "Could not find ShoppingList just overridden")
	// Sort the slices to make sure the order is deterministic
	slices.Sort(list1.Contents)
	slices.Sort(p.Contents)
	require.Equal(t, list1, p, "ShoppingList does not match the one just overridden")

	// Test implicit deletion of items
	list1.Contents = list1.Contents[:1]
	require.NoError(t, db.SetShoppingList(list1), "Could not override ShoppingList")
	p, err = db.LookupShoppingList(user, list1.Menu, list1.Pantry)
	require.NoError(t, err, "Could not find ShoppingList just overridden")
	// Sort the slices to make sure the order is deterministic
	slices.Sort(list1.Contents)
	slices.Sort(p.Contents)
	require.Equal(t, list1, p, "ShoppingList does not match the one just overridden")

	require.NoError(t, db.SetShoppingList(list2), "Could not set ShoppingList")
	p, err = db.LookupShoppingList(user, list2.Menu, list2.Pantry)
	require.NoError(t, err, "Could not find empty ShoppingList just created")
	// Sort the slices to make sure the order is deterministic
	slices.Sort(list2.Contents)
	slices.Sort(p.Contents)
	require.Equal(t, list2, p, "Empty menu does not match the one just created")

	t.Log("Closing DB and reopening")
	require.NoError(t, db.Close())

	db = openDB()
	defer db.Close()

	sLists, err := db.ShoppingLists(user)
	require.NoError(t, err)
	require.ElementsMatch(t, []dbtypes.ShoppingList{list1, list2}, sLists, "ShoppingLists do not match the ones after reopening DB")

	sLists, err = db.ShoppingLists(anotherUser)
	require.NoError(t, err)
	require.Empty(t, sLists, "Should not find ShoppingLists from another user")

	_, err = db.LookupShoppingList(user, "FAKE MENU", "FAKE PANTRY")
	require.ErrorIs(t, err, fs.ErrNotExist)

	p, err = db.LookupShoppingList(user, list1.Menu, list1.Pantry)
	require.NoError(t, err, "Could not find ShoppingList after reopening DB")
	// Sort the slices to make sure the order is deterministic
	slices.Sort(list1.Contents)
	slices.Sort(p.Contents)
	require.Equal(t, list1, p, "ShoppingList does not match the one after reopening DB")

	p, err = db.LookupShoppingList(user, list2.Menu, list2.Pantry)
	require.NoError(t, err, "Could not find empty ShoppingList after reopening DB")
	// Sort the slices to make sure the order is deterministic
	slices.Sort(list2.Contents)
	slices.Sort(p.Contents)
	require.Equal(t, list2, p, "Empty menu does not match the one after reopening DB")

	err = db.DeleteShoppingList(user, list1.Menu, list1.Pantry)
	require.NoError(t, err)

	_, err = db.LookupShoppingList(user, list1.Menu, list1.Pantry)
	require.ErrorIs(t, err, fs.ErrNotExist)

	require.NoError(t, db.Close())
}
