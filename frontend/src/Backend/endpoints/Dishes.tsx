export class DishesEndpoint {
    protected static path = '/api/recipes'

    Path(): string {
        return DishesEndpoint.path
    }

    async GET(): Promise<string[]> {
        return fetch(DishesEndpoint.path, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json'
            }
        })
            .then(r => r.ok ? r : Promise.reject(r))
            .then(r => r.json())
    }
}

export class MockDishesEndpoint extends DishesEndpoint {
    async GET(): Promise<string[]> {
        console.log(`GET from ${MockDishesEndpoint.path}`)
        return new Promise(resolve => setTimeout(resolve, 100))
            .then(() => [
                "Amanida de cigrons", "Amanida de llenties", "Cigrons amb xoriço", "Llenties amb xoriço", "Arròs amb ou", "Arròs amb sofregit", "Cereals amb llet", "Flocs de civada", "Torrada i suc de taronja", "Cafè", "Fajita", "Fruita", "Hamburguesa", "Proteïna amb acompanyament", "Vurguer", "Iogurt", "Iogurt vegà", "Macarrons amb sofregit", "Macarrons amb sofregit vegà", "Bastonets i hummus", "Pica-pica", "Puré de patata", "Sopa de galets", "Torrada d'alvocat", "Ramen", "Truita francesa", "Llom amb acompanyament", "Pollastre amb acompanyament", "Peix amb acompanyament", "Ravioli"
            ])
    }
}