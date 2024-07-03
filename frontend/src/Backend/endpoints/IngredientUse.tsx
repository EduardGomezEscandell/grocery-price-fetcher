export class IngredientUseEndpoint {
    path: string
    private auth: string

    constructor(auth: string, menu: string, ingredient: string) {
        this.path = `/api/ingredient-use/${menu}/${ingredient}`
        this.auth = auth
    }

    Path(): string {
        return this.path
    }

    async GET(): Promise<IngredientUsage[]> {
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
            .then((data: any[]) => IngredientUsage.fromJSON(data))
    }
}

export class MockIngredientUseEndpoint extends IngredientUseEndpoint {
    async GET(): Promise<IngredientUsage[]> {
        console.log(`GET to ${this.path}:`)
        return new Promise(resolve => setTimeout(resolve, 100))
            .then(() => IngredientUsage.fromJSON([
                { day: "Dilluns", meal: "Esmorzar", dish: "Torrada i suc de taronja", amount: 1 },
                { day: "Divendres", meal: "Esmorzar", dish: "Flocs de civada", amount: 3 },
                { day: "Dissabte", meal: "Dinar", dish: "Macarrons amb sofregit", amount: 1.7 },
                { day: "Dissabte", meal: "Dinar", dish: "Macarrons amb sofregit", amount: 1.7 },
                { day: "Dissabte", meal: "Dinar", dish: "Macarrons amb sofregit", amount: 1.7 },
                { day: "Dissabte", meal: "Dinar", dish: "Macarrons amb sofregit", amount: 1.7 },
                { day: "Dissabte", meal: "Dinar", dish: "Macarrons amb sofregit", amount: 1.7 },
                { day: "Dissabte", meal: "Dinar", dish: "Macarrons amb sofregit", amount: 1.7 },
            ]))
    }
}

export class IngredientUsage {
    constructor(day: string, meal: string, dish: string, amount: number) {
        this.day = day
        this.meal = meal
        this.dish = dish
        this.amount = amount
    }

    day: string;
    meal: string;
    dish: string;
    amount: number;

    static fromJSON(obj: any[]): IngredientUsage[] {
        return obj.map(x => {
            return new IngredientUsage(x.day, x.meal, x.dish, x.amount)
        })
    }
}