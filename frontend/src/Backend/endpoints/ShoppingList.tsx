import { ShoppingList } from '../../State/State'

export class ShoppingListEndpoint {
    protected path: string
    private auth: string

    constructor(auth: string, menu: string, pantry: string) {
        this.auth = auth
        this.path = `/api/shopping-list/${menu}/${pantry}`
    }

    Path(): string {
        return this.path
    }

    async GET(): Promise<ShoppingList> {
        return fetch(this.path, {
            method: 'GET',
            headers: {
                'Authorization': this.auth,
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
                'Authorization': this.auth,
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
                'Authorization': this.auth,
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

    constructor(auth: string, menu: string, pantry: string) {
        super(auth, menu, pantry)
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
                        { product_id: 2, name: "Pastanaga", units: 0, packs: 1, cost: 0.17, done: true },
                        { product_id: 17, name: "Pebrot verd", units: 0.50, packs: 1, cost: 0.50, done: false },
                        { product_id: 99, name: "Pebrot vermell", units: 18, packs: 3, cost: 1.10, done: true },
                        { product_id: 3, name: "Iogurt", units: 4.00, packs: 1, cost: 1.00 },
                        { product_id: 5, name: "Poma", units: 0.50, packs: 1, cost: 0.50 },
                        { product_id: 6, name: "Plàtan", units: 0.00, packs: 0, cost: 0.50 },
                        { product_id: 7, name: "Peres", units: 0.50, packs: 1, cost: 0.50 },
                        { product_id: 8, name: "Taronges", units: 0.50, packs: 1, cost: 0.50 },
                        { product_id: 9, name: "Maduixes", units: 0.50, packs: 1, cost: 0.50 },
                        { product_id: 10, name: "Kiwi", units: 0.50, packs: 1, cost: 0.50 },
                        { product_id: 11, name: "Mandarines", units: 0.50, packs: 1, cost: 0.50 },
                        { product_id: 12, name: "Pinya", units: 0.50, packs: 1, cost: 0.50 },
                        { product_id: 13, name: "Mango", units: 0.50, packs: 1, cost: 0.50 },
                        { product_id: 14, name: "Pera", units: 0.50, packs: 1, cost: 0.50 },
                        { product_id: 15, name: "Cireres", units: 0.50, packs: 1, cost: 0.50 },
                        { product_id: 16, name: "Préssecs", units: 0.50, packs: 1, cost: 0.50 },
                        { product_id: 1, name: "Albercocs", units: 0.50, packs: 1, cost: 0.50 },
                        { product_id: 18, name: "Nectarines", units: 0.50, packs: 1, cost: 0.50 },
                        { product_id: 19, name: "Pressec de vinya", units: 0.50, packs: 1, cost: 0.50 },
                        { product_id: 20, name: "Poma àcida", units: 0.50, packs: 1, cost: 0.50 },
                        { product_id: 21, name: "Poma verda", units: 0.50, packs: 1, cost: 0.50 },
                        { product_id: 22, name: "Ceba", units: 0.50, packs: 1, cost: 0.50 },
                        { product_id: 23, name: "Ceba vermella", units: 0.50, packs: 1, cost: 0.50 },
                        { product_id: 24, name: "Ceba tendra", units: 0.50, packs: 1, cost: 0.50 },
                        { product_id: 25, name: "Ceba de figueres", units: 0.50, packs: 1, cost: 0.50 },
                        { product_id: 26, name: "Ceba de calçot", units: 0.50, packs: 1, cost: 0.50 },
                        { product_id: 27, name: "Pipes", units: 0.50, packs: 1, cost: 0.50 },
                        { product_id: 28, name: "Nous", units: 0.50, packs: 1, cost: 0.50 },
                        { product_id: 29, name: "Ametlles", units: 0.50, packs: 1, cost: 0.50 },
                        { product_id: 30, name: "Avellanes", units: 0.50, packs: 1, cost: 0.50 },
                        { product_id: 31, name: "Pinyons", units: 0.50, packs: 1, cost: 0.50 },
                        { product_id: 32, name: "Anacards", units: 0.50, packs: 1, cost: 0.50 },
                        { product_id: 33, name: "Cacauets", units: 0.50, packs: 1, cost: 0.50 },
                        { product_id: 34, name: "Pistatxos", units: 0.50, packs: 1, cost: 0.50 },
                        { product_id: 35, name: "Garrofons", units: 0.50, packs: 1, cost: 0.50 },
                        { product_id: 36, name: "Mongetes", units: 0.50, packs: 1, cost: 0.50 },
                        { product_id: 37, name: "Mongetes tendres", units: 0.50, packs: 1, cost: 0.50 },
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
