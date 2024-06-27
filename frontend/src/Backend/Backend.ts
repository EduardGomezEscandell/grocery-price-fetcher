import { MenuEndpoint, MockMenuEndpoint } from './endpoints/Menu'
import { DishesEndpoint, MockDishesEndpoint } from './endpoints/Dishes'
import { PantryEndpoint, MockPantryEndpoint } from './endpoints/Pantry'
import { MockShoppingListEndpoint, ShoppingListEndpoint } from './endpoints/ShoppingList'
import { MockIngredientUseEndpoint, IngredientUseEndpoint } from './endpoints/IngredientUse'
import { MockNeedsEndpoint, NeedsEndpoint } from './endpoints/Needs'
import RecipeEndpoint, { MockRecipeEndpoint } from './endpoints/Recipe'
import Cache from './cache/Cache'
import ProductsEndpoint, { MockProductsEndpoint } from './endpoints/Products'
import ProviderEndpoint, { MockProvidersEndpoint } from './endpoints/Provider'

class Backend {
    constructor() {
        if (import.meta.env.VITE_APP_MOCK_BACKEND !== "") {
            this.mock = true
        }
    }

    private mock: boolean = false
    cache: Cache = new Cache()

    Provider(): ProviderEndpoint {
        return this.mock ? new MockProvidersEndpoint(this.cache) : new ProviderEndpoint(this.cache)
    }

    Products(namespace: string): ProductsEndpoint {
        return this.mock ? new MockProductsEndpoint(namespace, this.cache) : new ProductsEndpoint(namespace, this.cache)
    }

    Recipe(namespace: string, id: number): RecipeEndpoint {
        return this.mock ? new MockRecipeEndpoint(namespace, id, this.cache) : new RecipeEndpoint(namespace, id, this.cache)
    }

    Menu(which: string): MenuEndpoint {
        return this.mock ? new MockMenuEndpoint(which) : new MenuEndpoint(which)
    }

    Dishes(): DishesEndpoint {
        return this.mock ? new MockDishesEndpoint() : new DishesEndpoint()
    }

    Pantry(which: string): PantryEndpoint {
        return this.mock ? new MockPantryEndpoint(which) : new PantryEndpoint(which)
    }

    Needs(which: string): NeedsEndpoint {
        return this.mock ? new MockNeedsEndpoint(which) : new NeedsEndpoint(which)
    }

    IngredientUse(menu: string, ingredient: string): IngredientUseEndpoint {
        return this.mock ? new MockIngredientUseEndpoint(menu, ingredient) : new IngredientUseEndpoint(menu, ingredient)
    }

    ShoppingList(menu: string, pantry: string): ShoppingListEndpoint {
        return this.mock ? new MockShoppingListEndpoint(menu, pantry) : new ShoppingListEndpoint(menu, pantry)
    }

    ClearCache() {
        this.cache.clear()
    }
}

export default Backend
