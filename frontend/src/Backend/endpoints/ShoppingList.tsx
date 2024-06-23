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

    async PUT(items: number[]): Promise<void> {
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
                        { id: 1, name: "Pastanaga", units: 0, packs: 1, cost: 0.17, done: true },
                        { id: 2, name: "Pebrot verd", units: 0.50, packs: 1, cost: 0.50, done: false },
                        { id: 3, name: "Pebrot vermell", units: 18, packs: 3, cost: 1.10, done: true },
                        { id: 4, name: "Iogurt", units: 4.00, packs: 1, cost: 1.00 },
                        { id: 5, name: "Poma", units: 0.50, packs: 1, cost: 0.50 },
                        { id: 6, name: "Plàtan", units: 0.00, packs: 1, cost: 0.50 },
                        { id: 7, name: "Peres", units: 0.50, packs: 1, cost: 0.50 },
                        { id: 8, name: "Taronges", units: 0.50, packs: 1, cost: 0.50 },
                        { id: 9, name: "Maduixes", units: 0.50, packs: 1, cost: 0.50 },
                        { id: 10, name: "Kiwi", units: 0.50, packs: 1, cost: 0.50 },
                        { id: 11, name: "Mandarines", units: 0.50, packs: 1, cost: 0.50 },
                        { id: 12, name: "Pinya", units: 0.50, packs: 1, cost: 0.50 },
                        { id: 13, name: "Mango", units: 0.50, packs: 1, cost: 0.50 },
                        { id: 14, name: "Pera", units: 0.50, packs: 1, cost: 0.50 },
                        { id: 15, name: "Cireres", units: 0.50, packs: 1, cost: 0.50 },
                        { id: 16, name: "Préssecs", units: 0.50, packs: 1, cost: 0.50 },
                        { id: 17, name: "Albercocs", units: 0.50, packs: 1, cost: 0.50 },
                        { id: 18, name: "Nectarines", units: 0.50, packs: 1, cost: 0.50 },
                        { id: 19, name: "Pressec de vinya", units: 0.50, packs: 1, cost: 0.50 },
                        { id: 20, name: "Poma àcida", units: 0.50, packs: 1, cost: 0.50 },
                        { id: 21, name: "Poma verda", units: 0.50, packs: 1, cost: 0.50 },
                        { id: 22, name: "Ceba", units: 0.50, packs: 1, cost: 0.50 },
                        { id: 23, name: "Ceba vermella", units: 0.50, packs: 1, cost: 0.50 },
                        { id: 24, name: "Ceba tendra", units: 0.50, packs: 1, cost: 0.50 },
                        { id: 25, name: "Ceba de figueres", units: 0.50, packs: 1, cost: 0.50 },
                        { id: 26, name: "Ceba de calçot", units: 0.50, packs: 1, cost: 0.50 },
                        { id: 27, name: "Pipes", units: 0.50, packs: 1, cost: 0.50 },
                        { id: 28, name: "Nous", units: 0.50, packs: 1, cost: 0.50 },
                        { id: 29, name: "Ametlles", units: 0.50, packs: 1, cost: 0.50 },
                        { id: 30, name: "Avellanes", units: 0.50, packs: 1, cost: 0.50 },
                        { id: 31, name: "Pinyons", units: 0.50, packs: 1, cost: 0.50 },
                        { id: 32, name: "Anacards", units: 0.50, packs: 1, cost: 0.50 },
                        { id: 33, name: "Cacauets", units: 0.50, packs: 1, cost: 0.50 },
                        { id: 34, name: "Pistatxos", units: 0.50, packs: 1, cost: 0.50 },
                        { id: 35, name: "Garrofons", units: 0.50, packs: 1, cost: 0.50 },
                        { id: 36, name: "Mongetes", units: 0.50, packs: 1, cost: 0.50 },
                        { id: 37, name: "Mongetes tendres", units: 0.50, packs: 1, cost: 0.50 },
                    ]
                }),
            )
    }

    async PUT(items: number[]): Promise<void> {
        console.log(`PUT to ${this.path}:`)
        console.log(JSON.stringify(items)) // Ensure toJSON is called without errors
        return new Promise(resolve => setTimeout(resolve, 100))
    }

    async DELETE(): Promise<void> {
        console.log(`DELETE to ${this.path}`)
        return new Promise(resolve => setTimeout(resolve, 100))
    }
}
