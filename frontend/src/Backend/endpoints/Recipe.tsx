import Cache from "../cache/Cache"

export default class RecipeEndpoint {
    path: string
    cache: Cache

    constructor(namespace: string, recipe: string, cache?: Cache) {
        this.path = `/api/recipe/${namespace}/${recipe}`
        this.cache = cache || new Cache()
    }

    Path(): string {
        return this.path
    }

    protected async get_uncached(): Promise<Recipe> {
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

    async GET(): Promise<Recipe> {
        const cached = this.cache.get<Recipe>(this.path)
        if (cached) return cached

        return this
            .get_uncached()
            .then((recipe) => {
                this.cache.set(this.path, recipe)
                return recipe
            })
    }

    protected async post_uncached(recipe: Recipe): Promise<void> {
        this.cache.delete(this.path)

        return fetch(this.path, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json',
            },
            body: JSON.stringify(recipe)
        })
            .then(r => r.ok ? r : Promise.reject(r))
            .then(() => { })
    }

    async POST(recipe: Recipe): Promise<void> {
        return this.post_uncached(recipe)
    }

    protected async delete_uncached(): Promise<void> {
        this.cache.delete(this.path)

        return fetch(this.path, {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json',
            },
        })
            .then(r => r.ok ? r : Promise.reject(r))
            .then(() => { })
    }

    async DELETE(): Promise<void> {
        return this.delete_uncached()
    }
}

export class MockRecipeEndpoint extends RecipeEndpoint {
    recipe: string
    constructor(namespace: string, recipe: string, cache?: Cache) {
        super(namespace, recipe, cache)
        this.recipe = recipe
    }

    protected async get_uncached(): Promise<Recipe> {
        console.log(`GET to ${this.path}:`)
        
        // Pseudo-random number between 2 and 7 so that the frontend can test
        // recipe loading with different amounts of ingredients
        let pseudoRandom = 0
        for (let i = 0; i < this.recipe.length; i++) {
            pseudoRandom += this.recipe.charCodeAt(i)
        }
        pseudoRandom = 2 + (pseudoRandom % 6)

        return new Promise(resolve => setTimeout(resolve, 1000))
            .then(() => Recipe.fromJSON({
                name: this.recipe,
                ingredients: [
                    { name: "Macarrons", unit_price: 1.33, amount: 0.25 },
                    { name: "Ceba", unit_price: 0.76, amount: 0.5 },
                    { name: "All", unit_price: 0.88, amount: 0.1 },
                    { name: "Tomàquet", unit_price: 0.44, amount: 2 },
                    { name: "Oli", unit_price: 0.2, amount: 0.1 },
                    { name: "Sal", unit_price: 2.1, amount: 0.01 },
                    { name: "Pebre", unit_price: 1.57, amount: 0.01 }
                ].slice(0, pseudoRandom)
            }))
    }

    protected async post_uncached(recipe: Recipe): Promise<void> {
        console.log(`PUT to ${this.path}:`)
        console.log(JSON.stringify(recipe))
        return new Promise(resolve => setTimeout(resolve, 100))
    }

    protected async delete_uncached(): Promise<void> {
        console.log(`DELETE to ${this.path}:`)
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
    unit_price: number
    amount: number

    constructor(name: string, unit_price: number, amount: number) {
        this.name = name
        this.unit_price = unit_price
        this.amount = amount
    }

    static fromJSON(obj: any): Ingredient {
        return new Ingredient(
            obj.name && String(obj.name) || "Unknown",
            obj.unit_price && Number(obj.unit_price) || 0,
            obj.amount && Number(obj.amount) || 0
        )
    }
}