import Cache from "../cache/Cache"

export default class RecipeEndpoint {
    path: string
    cache: Cache
    private auth: string

    constructor(auth: string, namespace: string, id: number, cache?: Cache) {
        this.auth = auth
        this.path = `/api/recipe/${namespace}/${id.toString()}`
        this.cache = cache || new Cache()
    }

    Path(): string {
        return this.path
    }

    protected async get_uncached(): Promise<Recipe> {
        return fetch(this.path, {
            method: 'GET',
            headers: {
                'Authorization': this.auth,
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

    protected async post_uncached(recipe: Recipe): Promise<number> {
        this.cache.delete(this.path)

        return fetch(this.path, {
            method: 'POST',
            headers: {
                'Authorization': this.auth,
                'Content-Type': 'application/json',
                'Accept': 'application/json',
            },
            body: JSON.stringify(recipe)
        })
            .then(r => r.ok ? r : Promise.reject(r))
            .then(r => r.json())
            .then((data: any) => data.id)
    }

    async POST(recipe: Recipe): Promise<number> {
        return this.post_uncached(recipe)
    }

    protected async delete_uncached(): Promise<void> {
        this.cache.delete(this.path)

        return fetch(this.path, {
            method: 'DELETE',
            headers: {
                'Authorization': this.auth,
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
    recipe_id: number
    constructor(auth: string, namespace: string, recipe_id: number, cache?: Cache) {
        super(auth, namespace, recipe_id, cache)
        this.recipe_id = recipe_id
    }

    protected async get_uncached(): Promise<Recipe> {
        console.log(`GET to ${this.path}:`)

        // Pseudo-random number between 2 and 7 so that the frontend can test
        // recipe loading with different amounts of ingredients
        let pseudoRandom = (Math.cos(this.recipe_id) + 1) * 1000
        pseudoRandom = 2 + (pseudoRandom % 6)

        return new Promise(resolve => setTimeout(resolve, 1000))
            .then(() => Recipe.fromJSON({
                id: this.recipe_id,
                name: "Macarrons amb sofregit",
                ingredients: [
                    { id: 1, name: "Macarrons", unit_price: 1.33, amount: 0.25 },
                    { id: 2, name: "Ceba", unit_price: 0.76, amount: 0.5 },
                    { id: 3, name: "All", unit_price: 0.88, amount: 0.1 },
                    { id: 4, name: "Tomàquet", unit_price: 0.44, amount: 2 },
                    { id: 5, name: "Oli", unit_price: 0.2, amount: 0.1 },
                    { id: 6, name: "Sal", unit_price: 2.1, amount: 0.01 },
                    { id: 7, name: "Pebre", unit_price: 1.57, amount: 0.01 }
                ].slice(0, pseudoRandom)
            }))
    }

    protected async post_uncached(recipe: Recipe): Promise<number> {
        console.log(`PUT to ${this.path}:`)
        console.log(JSON.stringify(recipe))
        return new Promise(resolve => setTimeout(resolve, 100))
            .then(() => this.recipe_id === 0 ? 55 : this.recipe_id)
    }

    protected async delete_uncached(): Promise<void> {
        console.log(`DELETE to ${this.path}:`)
        return new Promise(resolve => setTimeout(resolve, 100))
    }
}

export class Recipe {
    id: number
    name: string
    ingredients: Ingredient[]

    constructor(id: number, name: string, ingredients: Ingredient[]) {
        this.id = id
        this.name = name
        this.ingredients = ingredients
    }

    static fromJSON(obj: any): Recipe {
        return new Recipe(
            obj.id && Number(obj.id) || 0,
            obj.name && String(obj.name) || "Unknown",
            obj.ingredients && obj.ingredients.map((x: any) => Ingredient.fromJSON(x)) || []
        )
    }
}

export class Ingredient {
    id: number
    name: string
    unit_price: number
    amount: number

    constructor(id: number, name: string, unit_price: number, amount: number) {
        this.id = id
        this.name = name
        this.unit_price = unit_price
        this.amount = amount
    }

    static fromJSON(obj: any): Ingredient {
        return new Ingredient(
            obj.id && Number(obj.id) || 0,
            obj.name && String(obj.name) || "Unknown",
            obj.unit_price && Number(obj.unit_price) || 0,
            obj.amount && Number(obj.amount) || 0
        )
    }
}