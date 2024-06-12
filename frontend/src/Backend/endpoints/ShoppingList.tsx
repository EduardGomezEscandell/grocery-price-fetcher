import { ShoppingList } from '../../State/State'

export class ShoppingListEndpoint {
    protected path: string

    constructor(menu: string, pantry: string) {
        this.path = `/api/shopping-list/${menu}/${pantry}`
    }

    Path(): string {
        return this.path
    }

    async GET(): Promise<ShoppingList> {
        return fetch(this.path, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json',
            }
        })
            .then(r => r.ok ? r : Promise.reject(r))
            .then(r => r.json())
            .then(data => ShoppingList.fromJSON(data))
    }

    async PUT(items: string[]): Promise<void> {
        return fetch(this.path, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json',
            },
            body: JSON.stringify(items)
        })
            .then(r => r.ok ? r : Promise.reject(r))
            .then(() => { })
    }

    async DELETE(): Promise<void> {
        return fetch(this.path, {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json',
            }
        })
            .then(r => r.ok ? r : Promise.reject(r))
            .then(() => { })
    }
}

export class MockShoppingListEndpoint extends ShoppingListEndpoint {
    menu: string
    pantry: string

    constructor(menu: string, pantry: string) {
        super(menu, pantry)
        this.menu = menu
        this.pantry = pantry
    }

    async GET(): Promise<ShoppingList> {
        console.log(`GET from ${this.path}`)
        return new Promise(resolve => setTimeout(resolve, 100))
            .then(() =>

                ShoppingList.fromJSON({
                    "menu": this.menu,
                    "pantry": this.pantry,
                    "items": [
                        { name: "Pastanaga", units: 0, packs: 1, cost: 0.17, done: true },
                        { name: "Pebrot verd", units: 0.50, packs: 1, cost: 0.50, done: false },
                        { name: "Pebrot vermell", units: 18, packs: 3, cost: 1.10, done: true },
                        { name: "Iogurt", units: 4.00, packs: 1, cost: 1.00 },
                        { name: "Poma", units: 0.50, packs: 1, cost: 0.50 },
                        { name: "Plàtan", units: 0.00, packs: 1, cost: 0.50 },
                        { name: "Peres", units: 0.50, packs: 1, cost: 0.50 },
                        { name: "Taronges", units: 0.50, packs: 1, cost: 0.50 },
                        { name: "Maduixes", units: 0.50, packs: 1, cost: 0.50 },
                        { name: "Kiwi", units: 0.50, packs: 1, cost: 0.50 },
                        { name: "Mandarines", units: 0.50, packs: 1, cost: 0.50 },
                        { name: "Pinya", units: 0.50, packs: 1, cost: 0.50 },
                        { name: "Mango", units: 0.50, packs: 1, cost: 0.50 },
                        { name: "Pera", units: 0.50, packs: 1, cost: 0.50 },
                        { name: "Cireres", units: 0.50, packs: 1, cost: 0.50 },
                        { name: "Préssecs", units: 0.50, packs: 1, cost: 0.50 },
                        { name: "Albercocs", units: 0.50, packs: 1, cost: 0.50 },
                        { name: "Nectarines", units: 0.50, packs: 1, cost: 0.50 },
                        { name: "Pressec de vinya", units: 0.50, packs: 1, cost: 0.50 },
                        { name: "Poma àcida", units: 0.50, packs: 1, cost: 0.50 },
                        { name: "Poma verda", units: 0.50, packs: 1, cost: 0.50 },
                        { name: "Ceba", units: 0.50, packs: 1, cost: 0.50 },
                        { name: "Ceba vermella", units: 0.50, packs: 1, cost: 0.50 },
                        { name: "Ceba tendra", units: 0.50, packs: 1, cost: 0.50 },
                        { name: "Ceba de figueres", units: 0.50, packs: 1, cost: 0.50 },
                        { name: "Ceba de calçot", units: 0.50, packs: 1, cost: 0.50 },
                        { name: "Pipes", units: 0.50, packs: 1, cost: 0.50 },
                        { name: "Nous", units: 0.50, packs: 1, cost: 0.50 },
                        { name: "Ametlles", units: 0.50, packs: 1, cost: 0.50 },
                        { name: "Avellanes", units: 0.50, packs: 1, cost: 0.50 },
                        { name: "Pinyons", units: 0.50, packs: 1, cost: 0.50 },
                        { name: "Anacards", units: 0.50, packs: 1, cost: 0.50 },
                        { name: "Cacauets", units: 0.50, packs: 1, cost: 0.50 },
                        { name: "Pistatxos", units: 0.50, packs: 1, cost: 0.50 },
                        { name: "Garrofons", units: 0.50, packs: 1, cost: 0.50 },
                        { name: "Mongetes", units: 0.50, packs: 1, cost: 0.50 },
                        { name: "Mongetes tendres", units: 0.50, packs: 1, cost: 0.50 },
                    ]
                }),
            )
    }

    async PUT(items: string[]): Promise<void> {
        console.log(`PUT to ${this.path}:`)
        console.log(JSON.stringify(items)) // Ensure toJSON is called without errors
        return new Promise(resolve => setTimeout(resolve, 100))
    }

    async DELETE(): Promise<void> {
        console.log(`DELETE to ${this.path}`)
        return new Promise(resolve => setTimeout(resolve, 100))
    }
}
