export default class RecipeEndpoint {
    path: string
    cache: Map<string, Recipe>

    constructor(namespace: string, recipe: string) {
        this.path = `/api/recipe/${namespace}/${recipe}/`
        this.cache = new Map<string, Recipe>()
    }

    Path(): string {
        return this.path
    }

    async GET(): Promise<Recipe> {
        return fetch(this.path, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json',
            },
        })
            .then(r => r.ok ? r : Promise.reject(r))
            .then(r => r.json())
            .then((data: any[]) => Recipe.fromJSON(data))
    }

    async PUT(recipe: Recipe): Promise<void> {
        return fetch(this.path, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json',
            },
            body: JSON.stringify(recipe)
        })
            .then(r => r.ok ? r : Promise.reject(r))
            .then(() => { })
    }
}

export class MockRecipeEndpoint extends RecipeEndpoint {
    recipe: string
    constructor(namespace: string, recipe: string) {
        super(namespace, recipe)
        this.recipe = recipe
    }

    async GET(): Promise<Recipe> {
        console.log(`GET to ${this.path}:`)
        return new Promise(resolve => setTimeout(resolve, 100))
            .then(() => Recipe.fromJSON({
                name: this.recipe,
                ingredients: [
                    { name: "Macarrons", unitPrice: 1.33, amount: 0.25 },
                    { name: "Ceba", unitPrice: 0.76, amount: 0.5 },
                    { name: "All", unitPrice: 0.88, amount: 0.1 },
                    { name: "Tomàquet", unitPrice: 0.44, amount: 2 },
                    { name: "Oli", unitPrice: 0.2, amount: 0.1 },
                    { name: "Sal", unitPrice: 2.1, amount: 0.01 },
                    { name: "Pebre", unitPrice: 1.57, amount: 0.01 }
                ]
            }))
    }

    async PUT(recipe: Recipe): Promise<void> {
        console.log(`PUT to ${this.path}:`)
        console.log(JSON.stringify(recipe))
        return new Promise(resolve => setTimeout(resolve, 100))
    }
}

export class Recipe {
    name: string
    ingredients: Ingredient[]

    constructor(name: string, ingredients: Ingredient[]) {
        this.name = name
        this.ingredients = ingredients
    }

    static fromJSON(obj: any): Recipe {
        return new Recipe(
            obj.name && String(obj.name) || "Unknown",
            obj.ingredients && obj.ingredients.map((x: any) => Ingredient.fromJSON(x)) || []
        )
    }
}

export class Ingredient {
    name: string
    unitPrice: number
    amount: number

    constructor(name: string, unitPrice: number, amount: number) {
        this.name = name
        this.unitPrice = unitPrice
        this.amount = amount
    }

    static fromJSON(obj: any): Ingredient {
        return new Ingredient(
            obj.name && String(obj.name) || "Unknown",
            obj.unitPrice && Number(obj.unitPrice) || 0,
            obj.amount && Number(obj.amount) || 0
        )
    }
}