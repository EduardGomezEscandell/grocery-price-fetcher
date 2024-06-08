export class IngredientUseEndpoint {
    protected static path: string = '/api/ingredient-use'

    static Path(): string {
        return this.path
    }

    async POST(req: reqBody): Promise<IngredientUsage[]> {
        return fetch(IngredientUseEndpoint.path, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json',
            },
            body: toJSON(req)
        })
            .then(response => response.json())
            .then(data => fromJSON(data))
    }
}

interface reqBody {
    MenuName: string;
    IngredientName: string;
}

function toJSON(r: reqBody): string {
    return JSON.stringify({
        menu_name: r.MenuName,
        ingredient_name: r.IngredientName
    } as any)
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
    async POST(req: reqBody): Promise<IngredientUsage[]> {
        console.log(`POST to ${MockIngredientUseEndpoint.path}:`)
        console.log(toJSON(req).substring(0, 30), '...') // Ensure toJSON is called without errors
        return new Promise(resolve => setTimeout(resolve, 100))
            .then(() => [
                { day: "Dilluns", meal: "Esmorzar", dish: "Torrada i suc de taronja", amount: 1 } as IngredientUsage,
                { day: "Divendres", meal: "Esmorzar", dish: "Flocs de civada", amount: 3 } as IngredientUsage,
                { day: "Dissabte", meal: "Dinar", dish: "Macarrons amb sofregit", amount: 1.7 } as IngredientUsage,
                { day: "Dissabte", meal: "Dinar", dish: "Macarrons amb sofregit", amount: 1.7 } as IngredientUsage,
                { day: "Dissabte", meal: "Dinar", dish: "Macarrons amb sofregit", amount: 1.7 } as IngredientUsage,
                { day: "Dissabte", meal: "Dinar", dish: "Macarrons amb sofregit", amount: 1.7 } as IngredientUsage,
                { day: "Dissabte", meal: "Dinar", dish: "Macarrons amb sofregit", amount: 1.7 } as IngredientUsage,
                { day: "Dissabte", meal: "Dinar", dish: "Macarrons amb sofregit", amount: 1.7 } as IngredientUsage,
            ])
    }
}
