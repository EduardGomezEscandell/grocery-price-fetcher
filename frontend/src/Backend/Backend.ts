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
import { AuthLoginEndpoint, MockAuthLoginEndpoint } from './endpoints/AuthLogin'
import { AuthLogoutEndpoint, MockAuthLogoutEndpoint } from './endpoints/AuthLogout'
import { AuthRefreshEndpoint, MockAuthRefreshEndpoint } from './endpoints/AuthRefresh'

class Backend {
    private mock: boolean = false
    private cache: Cache = new Cache()
    private getAuth: () => string

    constructor(authProvider: () => string | undefined) {
        this.mock = Backend.IsMock()
        this.getAuth = () => authProvider() || ""
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
        return this.mock ? new MockAuthLogoutEndpoint(this.getAuth()) : new AuthLogoutEndpoint(this.getAuth())
    }

    Provider(): ProviderEndpoint {
        return this.mock ? new MockProvidersEndpoint(this.getAuth(), this.cache) : new ProviderEndpoint(this.getAuth(), this.cache)
    }

    Products(): ProductsEndpoint {
        return this.mock ? new MockProductsEndpoint(this.getAuth(), this.cache) : new ProductsEndpoint(this.getAuth(), this.cache)
    }

    Recipe(id: number): RecipeEndpoint {
        return this.mock ? new MockRecipeEndpoint(this.getAuth(), id, this.cache) : new RecipeEndpoint(this.getAuth(), id, this.cache)
    }

    Menu(which: string): MenuEndpoint {
        return this.mock ? new MockMenuEndpoint(this.getAuth(), which) : new MenuEndpoint(this.getAuth(), which)
    }

    Dishes(): DishesEndpoint {
        return this.mock ? new MockDishesEndpoint(this.getAuth()) : new DishesEndpoint(this.getAuth())
    }

    Pantry(which: string): PantryEndpoint {
        return this.mock ? new MockPantryEndpoint(this.getAuth(), which) : new PantryEndpoint(this.getAuth(), which)
    }

    Needs(which: string): NeedsEndpoint {
        return this.mock ? new MockNeedsEndpoint(this.getAuth(), which) : new NeedsEndpoint(this.getAuth(), which)
    }

    IngredientUse(menu: string, ingredient: string): IngredientUseEndpoint {
        return this.mock ? new MockIngredientUseEndpoint(this.getAuth(), menu, ingredient) : new IngredientUseEndpoint(this.getAuth(), menu, ingredient)
    }

    ShoppingList(menu: string, pantry: string): ShoppingListEndpoint {
        return this.mock ? new MockShoppingListEndpoint(this.getAuth(), menu, pantry) : new ShoppingListEndpoint(this.getAuth(), menu, pantry)
    }

    ClearCache() {
        this.cache.clear()
    }
}

export default Backend
