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
import { LoginEndpoint, MockLoginEndpoint } from './endpoints/Login'
import { LogoutEndpoint, MockLogoutEndpoint } from './endpoints/Logout'

class Backend {
    private authorization: string

    constructor(auth?: string) {
        if (Backend.IsMock()) {
            this.mock = true
        }

        if (auth) {
            this.authorization = auth
        } else {
            this.authorization = ''
        }
    }

    private mock: boolean = false
    cache: Cache = new Cache()

    static IsMock(): boolean {
        return import.meta.env.VITE_APP_MOCK_BACKEND !== ""
    }

    Login(): LoginEndpoint {
        return this.mock ? new MockLoginEndpoint() : new LoginEndpoint()
    }

    Logout(): LogoutEndpoint {
        return this.mock ? new MockLogoutEndpoint(this.authorization) : new LogoutEndpoint(this.authorization)
    }

    Provider(): ProviderEndpoint {
        return this.mock ? new MockProvidersEndpoint(this.authorization, this.cache) : new ProviderEndpoint(this.authorization, this.cache)
    }

    Products(): ProductsEndpoint {
        return this.mock ? new MockProductsEndpoint(this.authorization, this.cache) : new ProductsEndpoint(this.authorization, this.cache)
    }

    Recipe(id: number): RecipeEndpoint {
        return this.mock ? new MockRecipeEndpoint(this.authorization, id, this.cache) : new RecipeEndpoint(this.authorization, id, this.cache)
    }

    Menu(which: string): MenuEndpoint {
        return this.mock ? new MockMenuEndpoint(this.authorization, which) : new MenuEndpoint(this.authorization, which)
    }

    Dishes(): DishesEndpoint {
        return this.mock ? new MockDishesEndpoint(this.authorization) : new DishesEndpoint(this.authorization)
    }

    Pantry(which: string): PantryEndpoint {
        return this.mock ? new MockPantryEndpoint(this.authorization, which) : new PantryEndpoint(this.authorization, which)
    }

    Needs(which: string): NeedsEndpoint {
        return this.mock ? new MockNeedsEndpoint(this.authorization, which) : new NeedsEndpoint(this.authorization, which)
    }

    IngredientUse(menu: string, ingredient: string): IngredientUseEndpoint {
        return this.mock ? new MockIngredientUseEndpoint(this.authorization, menu, ingredient) : new IngredientUseEndpoint(this.authorization, menu, ingredient)
    }

    ShoppingList(menu: string, pantry: string): ShoppingListEndpoint {
        return this.mock ? new MockShoppingListEndpoint(this.authorization, menu, pantry) : new ShoppingListEndpoint(this.authorization, menu, pantry)
    }

    ClearCache() {
        this.cache.clear()
    }
}

export default Backend
