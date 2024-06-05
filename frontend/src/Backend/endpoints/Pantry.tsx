import { Pantry } from '../../State/State.tsx'

export class PantryEndpoint {
    protected static path: string = '/api/pantry'

    static Path(): string {
        return this.path
    }

    async GET(): Promise<Pantry[]> {
        return fetch(PantryEndpoint.path, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json',
            }
        })
            .then(response => response.json())
            .then(data => data.map(p => Pantry.fromJSON(p)))
    }

    async POST(msg: PostMessage): Promise<void> {
        return fetch(PantryEndpoint.path, {
            method: 'POST',
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

interface PostMessage {
    name: string
    contents: Array<{
        name: string
        amount: number
    }>
}

export class MockPantryEndpoint extends PantryEndpoint {
    async GET(): Promise<Pantry[]> {
        console.log(`GET from ${MockPantryEndpoint.path}`)
        return new Promise(resolve => setTimeout(resolve, 100))
            .then(() =>
                [
                    Pantry.fromJSON({
                        name: "default", contents: [
                            { name: "Albercocs", amount: 7 },
                            { name: "Pastanaga", amount: 3 },
                            { name: "Iogurt", amount: 2 },
                        ]
                    }),
                    Pantry.fromJSON({ name: "Dummy menu" })
                ]
            )
    }

    async POST(msg: PostMessage): Promise<void> {
        console.log(`POST to ${MockPantryEndpoint.path}:`)
        console.log(JSON.stringify(msg)) // Ensure toJSON is called without errors
        return new Promise(resolve => setTimeout(resolve, 100))
    }
}
