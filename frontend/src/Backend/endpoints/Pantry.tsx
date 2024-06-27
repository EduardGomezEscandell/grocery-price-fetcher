import { Pantry } from '../../State/State'

export class PantryEndpoint {
    protected path: string
    constructor(protected which: string) {
        this.path = '/api/pantry/' + which
    }

    Path(): string {
        return this.path
    }

    async GET(): Promise<Pantry> {
        return fetch(this.path, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json',
            },
        })
            .then(r => r.ok ? r : Promise.reject(r))
            .then(r => r.json())
            .then(data => Pantry.fromJSON(data))
    }

    async PUT(msg: Pantry): Promise<void> {
        return fetch(this.path, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json',
            },
            body: JSON.stringify(msg)
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


export class MockPantryEndpoint extends PantryEndpoint {
    which: string
    constructor(which: string) {
        super(which)
        this.which = which
    }

    async GET(): Promise<Pantry> {
        console.log(`GET from ${this.path}`)
        return new Promise(resolve => setTimeout(resolve, 100))
            .then(() => Pantry.fromJSON({
                name: this.which, contents: [
                    { product_id: 1, name: "Albercocs", amount: 7 },
                    { product_id: 2, name: "Pastanaga", amount: 3 },
                    { product_id: 3, name: "Iogurt", amount: 2 },
                ]
            }),
        )
    }

    async PUT(msg: Pantry): Promise<void> {
        console.log(`PUT to ${this.path}:`)
        console.log(JSON.stringify(msg)) // Ensure toJSON is called without errors
        return new Promise(resolve => setTimeout(resolve, 100))
    }

    async DELETE(): Promise<void> {
        console.log(`DELETE from ${this.path}`)
        return new Promise(resolve => setTimeout(resolve, 100))
    }
}
