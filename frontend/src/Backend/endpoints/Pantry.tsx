import { Pantry } from '../../State/State.tsx'

export class PantryEndpoint {
    protected path: string = '/api/pantry'
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
            }
        })
            .then(response => response.json())
            .then(data => data.map(p => Pantry.fromJSON(p)))
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
            .then(() => { })
            .catch((error) => { console.error('Error:', error) })
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
                    { name: "Albercocs", amount: 7 },
                    { name: "Pastanaga", amount: 3 },
                    { name: "Iogurt", amount: 2 },
                ]
            }),
        )
    }

    async PUT(msg: Pantry): Promise<void> {
        console.log(`PUT to ${this.path}:`)
        console.log(JSON.stringify(msg)) // Ensure toJSON is called without errors
        return new Promise(resolve => setTimeout(resolve, 100))
    }
}
