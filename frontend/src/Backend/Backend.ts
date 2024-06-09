import { process } from 'react'
import { MenuEndpoint, MockMenuEndpoint } from './endpoints/Menu.tsx'
import { DishesEndpoint, MockDishesEndpoint } from './endpoints/Dishes.tsx'
import { PantryEndpoint, MockPantryEndpoint } from './endpoints/Pantry.tsx'
import { MockShoppingListEndpoint, ShoppingListEndpoint } from './endpoints/ShoppingList.tsx'
import { MockIngredientUseEndpoint, IngredientUseEndpoint } from './endpoints/IngredientUse.tsx'
import { MockNeedsEndpoint, NeedsEndpoint } from './endpoints/Needs.tsx'

class Backend {
    constructor() {
        if (process.env.REACT_APP_MOCK_BACKEND !== "") {
            this.mock = true
        }
    }

    private mock: boolean = false

    Menu(): MenuEndpoint {
        return this.mock ? new MockMenuEndpoint() : new MenuEndpoint()
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

    IngredientUse(): IngredientUseEndpoint {
        return this.mock ? new MockIngredientUseEndpoint() : new IngredientUseEndpoint()
    }

    Shopping(): ShoppingListEndpoint {
        return this.mock ? new MockShoppingListEndpoint() : new ShoppingListEndpoint()
    }
}

export default Backend