export class Ingredient {
    constructor(
        name: string,
        price: number,
        batch_size: number,
        have: number,
        need: number
    ) {
        this.name = name
        this.price = price
        this.batch_size = batch_size
        this.have = have
        this.need = need
    }

    name: string;
    price: number; // Price per batch

    batch_size: number;
    have: number;
    need: number;
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

export class ShoppingNeeds {
    static fromJSON(json: any): ShoppingNeeds {
        let shoppingList = new ShoppingNeeds()
        shoppingList.ingredients = json.map((ingredient: any) => {
            return new Ingredient(
                either(ingredient, 'product', 'Unnamed ingredient'),
                either(ingredient, 'price', 0),
                either(ingredient, 'batch_size', 1),
                either(ingredient, 'have', 0),
                either(ingredient, 'need', 0),
            )
        })
        return shoppingList
    }

    ingredients: Array<Ingredient> = [];
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

export class Pantry {
    name: string = 'default'
    contents: Array<Ingredient> = []

    static fromJSON(json: any) {
        try {
            let pantry = new Pantry()
            pantry.name = either(json, 'name', 'Default')
            pantry.contents = either(json, 'contents', []).map((content: any): Ingredient => {
                return new Ingredient(
                    /* name: */ either(content, 'name', 'Unnamed ingredient'),
                    /* price */ 0,
                    /* batch size */ 0,
                    /* have: */ either(content, 'amount', 0),
                    /* need */ 0
                )
            })
            return pantry
        } catch (e) {
            console.error(e)
            return new Pantry()
        }
    }
}

export class ShoppingListItem {
    name: string
    done: boolean
    units: number
    packs: number
    cost: number
}

export class ShoppingList {
    name: string = 'default'
    timeStamp: string = ''
    items: Array<ShoppingListItem> = []

    static fromJSON(json: any): ShoppingList {
        const shoppingList = new ShoppingList()
        shoppingList.name = either(json, 'name', 'default')
        shoppingList.timeStamp = either(json, 'time_stamp', '2000-01-01T00:00:00Z00:00')
        shoppingList.items = either(json, 'items', []).map((name: string) => {
            return {
                name: name,
                done: true,
            } as ShoppingListItem
        })
        console.log(shoppingList)
        return shoppingList
    }
}

function either<T>(struct: any, key: string, val: T): T {
    return struct[key] || val
}

export class State {
    constructor() {
        this.dishes = []
        this.inNeed = new ShoppingNeeds()
    }

    attachMenu(menu: Menu, setMenu: (m: Menu) => void): State {
        this.menu = menu
        this._setMenu = setMenu
        return this
    }

    menu: Menu;
    private _setMenu: (m: Menu) => void;

    setMenu(menu: Menu) {
        menu.days.forEach(day => {
            day.meals.forEach(meal => {
                meal.dishes = meal.dishes
                    .filter(dish => dish.name !== "")
                    .filter(dish => dish.amount > 0)
            })
        })
        this._setMenu(menu)
    }

    dishes: Array<string>;
    inNeed: ShoppingNeeds;
    shoppingList: ShoppingList;
}

