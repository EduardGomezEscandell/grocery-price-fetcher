import { MenuEndpoint, MockMenuEndpoint } from './endpoints/Menu'
import { DishesEndpoint, MockDishesEndpoint } from './endpoints/Dishes'
import { PantryEndpoint, MockPantryEndpoint } from './endpoints/Pantry'
import { MockShoppingListEndpoint, ShoppingListEndpoint } from './endpoints/ShoppingList'
import { MockIngredientUseEndpoint, IngredientUseEndpoint } from './endpoints/IngredientUse'
import { MockNeedsEndpoint, NeedsEndpoint } from './endpoints/Needs'
import RecipeEndpoint, { MockRecipeEndpoint } from './endpoints/Recipe'
import ProductsEndpoint, { MockProductsEndpoint } from './endpoints/Products'
import ProviderEndpoint, { MockProvidersEndpoint } from './endpoints/Provider'
import { AuthLoginEndpoint, MockAuthLoginEndpoint } from './endpoints/AuthLogin'
import { AuthLogoutEndpoint, MockAuthLogoutEndpoint } from './endpoints/AuthLogout'
import { AuthRefreshEndpoint, MockAuthRefreshEndpoint } from './endpoints/AuthRefresh'
import Cache from './cache/Cache'

class Backend {
    private mock: boolean = false
    private cache: Cache = new Cache()

    constructor() {
        this.mock = Backend.IsMock()
    }

    static IsMock(): boolean {
        return import.meta.env.VITE_APP_MOCK_BACKEND !== ""
    }

    AuthLogin(): AuthLoginEndpoint {
        return this.mock ? new MockAuthLoginEndpoint() : new AuthLoginEndpoint()
    }

    AuthRefresh(): AuthRefreshEndpoint {
        return this.mock ? new MockAuthRefreshEndpoint() : new AuthRefreshEndpoint()
    }

    AuthLogout(): AuthLogoutEndpoint {
        return this.mock ? new MockAuthLogoutEndpoint() : new AuthLogoutEndpoint()
    }

    Provider(): ProviderEndpoint {
        return this.mock ? new MockProvidersEndpoint(this.cache) : new ProviderEndpoint(this.cache)
    }

    Products(): ProductsEndpoint {
        return this.mock ? new MockProductsEndpoint(this.cache) : new ProductsEndpoint(this.cache)
    }

    Recipe(id: number): RecipeEndpoint {
        return this.mock ? new MockRecipeEndpoint(id, this.cache) : new RecipeEndpoint(id, this.cache)
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
