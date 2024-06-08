import { process } from 'react'
import { MenuEndpoint, MockMenuEndpoint } from './endpoints/Menu.tsx'
import { DishesEndpoint, MockDishesEndpoint } from './endpoints/Dishes.tsx'
import { PantryEndpoint, MockPantryEndpoint } from './endpoints/Pantry.tsx'
import { MockShoppingListEndpoint, ShoppingListEndpoint } from './endpoints/ShoppingList.tsx'
import { MockIngredientUseEndpoint, IngredientUseEndpoint } from './endpoints/IngredientUse.tsx'

class Backend {
    static New(): Backend {
        if (process.env.REACT_APP_MOCK_BACKEND !== "") {
            return new Backend(true)
        }
        return new Backend(false)
    }

    private constructor(mock: boolean = false) {
        this.menu = mock ? new MockMenuEndpoint() : new MenuEndpoint()
        this.dishes = mock ? new MockDishesEndpoint() : new DishesEndpoint()
        this.pantry = mock ? new MockPantryEndpoint() : new PantryEndpoint()
        this.shopping =  mock ? new MockShoppingListEndpoint() : new ShoppingListEndpoint()
        this.ingredientUse = mock ? new MockIngredientUseEndpoint() : new IngredientUseEndpoint()
    }

    private menu: MenuEndpoint;
    Menu(): MenuEndpoint {
        return this.menu
    }

    private dishes: DishesEndpoint;
    Dishes(): DishesEndpoint {
        return this.dishes
    }

    private pantry: PantryEndpoint;
    Pantry(): PantryEndpoint {
        return this.pantry
    }

    private ingredientUse: IngredientUseEndpoint;
    IngredientUse(): IngredientUseEndpoint {
        return this.ingredientUse
    }

    private shopping: ShoppingListEndpoint;
    Shopping(): ShoppingListEndpoint {
        return this.shopping
    }
}

export default Backend