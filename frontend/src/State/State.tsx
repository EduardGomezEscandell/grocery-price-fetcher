export class Product {
    constructor(
        name: string,
        price: number,
        batch_size: number,
    ) {
        this.name = name
        this.price = price
        this.batch_size = batch_size
    }

    name: string;
    price: number; // Price per batch
    batch_size: number;
}

export class Dish {
    constructor(name: string, amount: number) {
        this.name = name
        this.amount = amount
    }

    name: string;
    amount: number;

    withName(name: string) {
        this.name = name
        return this
    }

    withAmount(amount: number) {
        this.amount = amount
        return this
    }
}

export class Meal {
    constructor(name: string = '') {
        this.name = name
    }

    name: string = '';
    dishes: Array<Dish> = [];
}

export class Day {
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
                let d = new Day()
                d.name = day.name
                d.meals = either(day, 'meals', []).map((meal: any) => {
                    let m = new Meal()
                    m.name = either(meal, 'name', 'Unnamed meal')
                    m.dishes = either(meal, 'dishes', []).map((dish: any) => {
                        return new Dish(dish.name, dish.amount)
                    })
                    return m
                })
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
        const copy = {}
        copy['name'] = this.name
        copy['days'] = this.days.map(day => {
            const d = {}
            d['name'] = day.name
            d['meals'] = day.meals
                .filter(meal => meal.name !== "")
                .map(meal => {
                    const m = {}
                    m['name'] = meal.name
                    m['dishes'] = meal.dishes
                        .filter(dish => dish.name !== "")
                        .map(dish => {
                            return {
                                name: dish.name,
                                amount: dish.amount
                            }
                        })
                    return m
                })
            return d
        })
        return JSON.stringify(copy)
    }
}

export class PantryItem {
    name: string
    amount: number
}

export class Pantry {
    name: string = 'default'
    contents: Array<PantryItem> = []

    static fromJSON(json: any) {
        try {
            let pantry = new Pantry()
            pantry.name = either(json, 'name', 'Default')
            pantry.contents = either(json, 'contents', []).map((content: any): PantryItem => {
                return {
                    name: either(content, 'name', 'Unnamed ingredient'),
                    amount: either(content, 'amount', 0),
                }
            })
            return pantry
        } catch (e) {
            console.error(e)
            return new Pantry()
        }
    }
}

export class ShoppingNeedsItem {
    static fromJSON(json: any): ShoppingNeedsItem {
        const n = new ShoppingNeedsItem(
            either(json, 'name', 'Unnamed ingredient'),
            either(json, 'amount', 0),
        )
        return n
    }

    constructor(name: string, amount: number) {
        this.name = name
        this.amount = amount
    }

    name: string
    amount: number
}

export class ShoppingNeeds {
    static fromJSON(json: any): ShoppingNeeds {
        const need = new ShoppingNeeds()
        need.items = either(json, 'items', []).map((ingredient: any) => ShoppingNeedsItem.fromJSON(ingredient))
        return need
    }

    menu: string = 'default'
    items: Array<ShoppingNeedsItem> = [];
}

export class ShoppingListItem {
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
                name: either(name, 'name', 'Unnamed ingredient'),
                done: either(name, 'done', false),
                units: either(name, 'units', 0),
                packs: either(name, 'packs', 0),
                cost: either(name, 'cost', 0),
            } as ShoppingListItem
        })
        console.log(shoppingList)
        return shoppingList
    }
}

function either<T>(struct: any, key: string, val: T): T {
    return struct[key] || val
}
