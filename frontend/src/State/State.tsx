export class Product {
    constructor(
        id: number,
        name: string,
        price: number,
        batch_size: number,
        provider: string,
        provider_code: string,
    ) {
        this.id = id
        this.name = name
        this.price = price
        this.batch_size = batch_size
        this.provider = provider
        this.product_code = provider_code
    }

    static fromJSON(json: any): Product {
        return new Product(
            json.id,
            json.name,
            json.price,
            json.batch_size,
            json.provider,
            json.product_code[0],
        )
    }

    id: number;
    name: string;
    price: number; // Price per batch
    batch_size: number;
    provider: string;
    product_code: string;
}

export class Dish {
    constructor(id: number, name: string, amount: number) {
        this.id = id
        this.name = name
        this.amount = amount
    }

    id: number;
    name: string;
    amount: number;
}

export class Meal {
    constructor(name: string = '') {
        this.name = name
    }

    name: string = '';
    dishes: Array<Dish> = [];
}

export class Day {
    constructor(name: string, meals: Array<Meal> = []) {
        this.name = name
        this.meals = meals
    }

    name: string;
    meals: Array<Meal>;
}

export class Menu {
    days: Array<Day> = [];
    name: string = 'default';

    static fromJSON(json: any): Menu {
        let menu = new Menu()

        try {
            menu.name = either(json, 'name', 'Unnamed menu')
            menu.days = either(json, 'days', []).map((day: any) => {
                let d = new Day(
                    day.name,
                    either(day, 'meals', []).map((meal: any) => {
                        let m = new Meal()
                        m.name = either(meal, 'name', 'Unnamed meal')
                        m.dishes = either(meal, 'dishes', []).map((dish: any) => {
                            return new Dish(dish.recipe_id, dish.name, dish.amount)
                        })
                        return m
                    })
                )
                return d
            })

            // Padding missing meals
            const meals = Array.from(new Set(menu.days.flatMap(day => day.meals)))
            menu.days.forEach(day => {
                meals.forEach(meal => {
                    if (!day.meals.find(m => m.name === meal.name)) {
                        day.meals.push(new Meal(meal.name))
                    }
                })
            })
        } catch (e) {
            console.error(e)
        }

        return menu
    }

    toJSON(): string {
        const copy = {
            name: this.name,
            days: this.days.map(day => {
                return {
                    name: day.name,
                    meals: day.meals
                        .filter(meal => meal.name !== "")
                        .map(meal => {
                            return {
                                name: meal.name,
                                dishes: meal.dishes
                                    .filter(dish => dish.name !== "")
                                    .map(dish => {
                                        return {
                                            recipe_id: dish.id,
                                            amount: dish.amount
                                        }
                                    })
                            }
                        })

                }
            })
        }
        return JSON.stringify(copy)
    }
}

// Ingredient represents a certain amount of a product.
// The amount is expressed in the product's unit.
export class Ingredient {
    constructor(id: number, name: string, amount: number) {
        this.product_id = id
        this.name = name
        this.amount = amount
    }

    static fromJSON(obj: any): Ingredient {
        return {
            product_id: obj.product_id,
            name: obj.name,
            amount: either(obj, 'amount', 0),
        }
    }

    product_id: number
    name: string
    amount: number
}

export class Pantry {
    name: string = 'default'
    contents: Array<Ingredient> = []

    constructor(name: string, contents: Array<Ingredient> = []) {
        this.name = name
        this.contents = contents
    }

    static fromJSON(json: any): Pantry {
        return new Pantry(
            either(json, 'name', 'Default'),
            either(json, 'contents', []).map(Ingredient.fromJSON)
        )
    }
}

export class ShoppingNeeds {
    static fromJSON(json: any): ShoppingNeeds {
        const need = new ShoppingNeeds()
        need.items = either(json, 'items', []).map((ingredient: any) => Ingredient.fromJSON(ingredient))
        return need
    }

    menu: string = 'default'
    items: Array<Ingredient> = [];
}

export interface ShoppingListItem {
    id: number
    name: string
    done: boolean
    units: number
    packs: number
    cost: number
}

export class ShoppingList {
    menu: string = 'default'
    pantry: string = 'default'
    items: Array<ShoppingListItem> = []

    static fromJSON(json: any): ShoppingList {
        const shoppingList = new ShoppingList()
        shoppingList.menu = either(json, 'menu', 'default')
        shoppingList.pantry = either(json, 'pantry', 'default')
        shoppingList.items = either(json, 'items', []).map((name: string) => {
            return {
                id: either(name, 'product_id', 0),
                name: either(name, 'name', 'Unnamed ingredient'),
                done: either(name, 'done', false),
                units: either(name, 'units', 0),
                packs: either(name, 'packs', 0),
                cost: either(name, 'cost', 0),
            }
        })
        return shoppingList
    }
}

function either<T>(struct: any, key: string, val: T): T {
    return struct[key] || val
}
