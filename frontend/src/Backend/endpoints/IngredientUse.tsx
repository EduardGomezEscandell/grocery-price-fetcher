export class IngredientUseEndpoint {
    path: string
    constructor(menu: string, ingredient: string) {
        this.path = `/api/ingredient-use/${menu}/${ingredient}`
    }

    Path(): string {
        return this.path
    }

    async GET(): Promise<IngredientUsage[]> {
        return fetch(this.path, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json',
            },
        })
            .then(r => r.ok ? r : Promise.reject(r))
            .then(r => r.json())
            .then(data => fromJSON(data))
    }
}

function fromJSON(obj: any[]): IngredientUsage[] {
    return obj.map(x => {
        return {
            day: x.day,
            meal: x.meal,
            dish: x.dish,
            amount: x.amount
        } as IngredientUsage
    })
}

export class IngredientUsage {
    day: string;
    meal: string;
    dish: string;
    amount: number;
}

export class MockIngredientUseEndpoint extends IngredientUseEndpoint {
    async GET(): Promise<IngredientUsage[]> {
        console.log(`GET to ${this.path}:`)
        return new Promise(resolve => setTimeout(resolve, 100))
            .then(() => fromJSON([
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
