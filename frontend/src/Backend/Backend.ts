import { MenuEndpoint, MockMenuEndpoint } from './endpoints/Menu'
import { DishesEndpoint, MockDishesEndpoint } from './endpoints/Dishes'
import { PantryEndpoint, MockPantryEndpoint } from './endpoints/Pantry'
import { MockShoppingListEndpoint, ShoppingListEndpoint } from './endpoints/ShoppingList'
import { MockIngredientUseEndpoint, IngredientUseEndpoint } from './endpoints/IngredientUse'
import { MockNeedsEndpoint, NeedsEndpoint } from './endpoints/Needs'

class Backend {
    constructor() {
        if (import.meta.env.VITE_APP_MOCK_BACKEND !== "") {
            this.mock = true
        }
    }

    private mock: boolean = false

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
}

export default Backend
