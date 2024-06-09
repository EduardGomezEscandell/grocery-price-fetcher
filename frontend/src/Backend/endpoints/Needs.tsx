import { ShoppingNeeds } from '../../State/State.tsx'

export class NeedsEndpoint {
    protected path: string

    constructor(which: string) {
        this.path = '/api/needs/' + which
    }

    Path(): string {
        return this.path
    }

    async GET(): Promise<ShoppingNeeds> {
        return fetch(this.path, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json',
            }
        })
            .then(response => response.json())
            .then(data => data.map(p => ShoppingNeeds.fromJSON(p)))
    }
}

export class MockNeedsEndpoint extends NeedsEndpoint {
    which: string
    constructor(which: string) {
        super(which)
        this.which = which
    }

    async GET(): Promise<ShoppingNeeds> {
        console.log(`GET from ${this.path}`)
        return new Promise(resolve => setTimeout(resolve, 100))
            .then(() => ShoppingNeeds.fromJSON({
                name: this.which,
                contents: [
                    { batch_size: 1, amount: 1.00, price: 0.17, name: "Pastanaga" },
                    { batch_size: 1, amount: 0.50, price: 0.50, name: "Pebrot verd" },
                    { batch_size: 1, amount: 0.95, price: 1.10, name: "Pebrot vermell" },
                    { batch_size: 4, amount: 4.00, price: 1.00, name: "Iogurt" },
                    { batch_size: 1, amount: 0.50, price: 0.50, name: "Poma" },
                    { batch_size: 1, amount: 0.50, price: 0.50, name: "Plàtan" },
                    { batch_size: 1, amount: 0.50, price: 0.50, name: "Peres" },
                    { batch_size: 1, amount: 0.50, price: 0.50, name: "Taronges" },
                    { batch_size: 1, amount: 0.50, price: 0.50, name: "Maduixes" },
                    { batch_size: 1, amount: 0.50, price: 0.50, name: "Kiwi" },
                    { batch_size: 1, amount: 0.50, price: 0.50, name: "Mandarines" },
                    { batch_size: 1, amount: 0.50, price: 0.50, name: "Pinya" },
                    { batch_size: 1, amount: 0.50, price: 0.50, name: "Mango" },
                    { batch_size: 1, amount: 0.50, price: 0.50, name: "Pera" },
                    { batch_size: 1, amount: 0.50, price: 0.50, name: "Cireres" },
                    { batch_size: 1, amount: 0.50, price: 0.50, name: "Préssecs" },
                    { batch_size: 1, amount: 0.50, price: 0.50, name: "Albercocs" },
                    { batch_size: 1, amount: 0.50, price: 0.50, name: "Nectarines" },
                    { batch_size: 1, amount: 0.50, price: 0.50, name: "Pressec de vinya" },
                    { batch_size: 1, amount: 0.50, price: 0.50, name: "Poma àcida" },
                    { batch_size: 1, amount: 0.50, price: 0.50, name: "Poma verda" },
                    { batch_size: 1, amount: 0.50, price: 0.50, name: "Ceba" },
                    { batch_size: 1, amount: 0.50, price: 0.50, name: "Ceba vermella" },
                    { batch_size: 1, amount: 0.50, price: 0.50, name: "Ceba tendra" },
                    { batch_size: 1, amount: 0.50, price: 0.50, name: "Ceba de figueres" },
                    { batch_size: 1, amount: 0.50, price: 0.50, name: "Ceba de calçot" },
                    { batch_size: 1, amount: 0.50, price: 0.50, name: "Pipes" },
                    { batch_size: 1, amount: 0.50, price: 0.50, name: "Nous" },
                    { batch_size: 1, amount: 0.50, price: 0.50, name: "Ametlles" },
                    { batch_size: 1, amount: 0.50, price: 0.50, name: "Avellanes" },
                    { batch_size: 1, amount: 0.50, price: 0.50, name: "Pinyons" },
                    { batch_size: 1, amount: 0.50, price: 0.50, name: "Anacards" },
                    { batch_size: 1, amount: 0.50, price: 0.50, name: "Cacauets" },
                    { batch_size: 1, amount: 0.50, price: 0.50, name: "Pistatxos" },
                    { batch_size: 1, amount: 0.50, price: 0.50, name: "Garrofons" },
                    { batch_size: 1, amount: 0.50, price: 0.50, name: "Mongetes" },
                    { batch_size: 1, amount: 0.50, price: 0.50, name: "Mongetes tendres" },
                ]
            }))
    }
}
