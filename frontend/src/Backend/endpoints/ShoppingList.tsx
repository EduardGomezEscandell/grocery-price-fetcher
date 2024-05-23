import { ShoppingList } from '../../State/State.tsx'

export class ShoppingListEndpoint {
    protected static path: string = '/api/shopping-list'

    static Path(): string {
        return this.path
    }

    async GET(): Promise<ShoppingList[]> {
        return fetch(ShoppingListEndpoint.path, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json',
            }
        })
            .then(response => response.json())
            .then(data => data.map(p => ShoppingList.fromJSON(p)))
    }

    async POST(msg: PostMessage): Promise<void> {
        return fetch(ShoppingListEndpoint.path, {
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
    items: string[]
}

export class MockShoppingListEndpoint extends ShoppingListEndpoint {
    async GET(): Promise<ShoppingList[]> {
        console.log(`GET from ${MockShoppingListEndpoint.path}`)
        return new Promise(resolve => setTimeout(resolve, 100))
            .then(() =>
                [
                    ShoppingList.fromJSON({
                        "name": "default",
                        "time_stamp": "2024-05-24T14:57:36Z",
                        "items": ["Ametlles","Avellanes"],
                    }),
                    ShoppingList.fromJSON({
                        "name": "dummy",
                        "time_stamp": "2024-05-24T14:57:36Z",
                        "items": [],
                    }),
                ]
            )
    }

    async POST(msg: PostMessage): Promise<void> {
        console.log(`POST to ${MockShoppingListEndpoint.path}:`)
        console.log(JSON.stringify(msg)) // Ensure toJSON is called without errors
        return new Promise(resolve => setTimeout(resolve, 100))
    }
}
