export class Ingredient {
    constructor(name: string, price: number, have: number, need: number) {
        this.name = name
        this.price = price
        this.have = have
        this.need = need
    }

    name: string;
    price: number;

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

export class Pantry {
    ingredients: Array<Ingredient>;
}

export class ShoppingList {
    static fromJSON(json: any): ShoppingList {
        let shoppingList = new ShoppingList()
        shoppingList.ingredients = json.map((ingredient: any) => {
            return new Ingredient(
                either(ingredient, 'product', 'Unnamed ingredient'),
                either(ingredient, 'price', 0),
                either(ingredient, 'have', 0),
                either(ingredient, 'need', 0),
            )
        })
        return shoppingList
    }

    ingredients: Array<Ingredient>;
}

export class Menu {
    days: Array<Day>;
    name: string;

    static fromJson(json: any): Menu {
        let menu = new Menu()
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
    
        // Padding small meals
        const maxDishesPerMeal = Math.max(...menu.days.map(day => Math.max(...day.meals.map(meal => meal.dishes.length))))
        menu.days.forEach(day => {
            day.meals.forEach(meal => {
                while (meal.dishes.length < maxDishesPerMeal) {
                    meal.dishes.push(new Dish("", 0))
                }
            })
        })

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

function either<T>(struct: any, key: string, val: T): T {
    return struct[key] || val
}

export class State {
    constructor() {
        this.dishes = []
        this.shoppingList = new ShoppingList()
        this.menu = new Menu()
        this.pantry = new Pantry()
    }

    dishes: Array<string>;
    shoppingList: ShoppingList;
    menu: Menu;
    pantry: Pantry;
}

