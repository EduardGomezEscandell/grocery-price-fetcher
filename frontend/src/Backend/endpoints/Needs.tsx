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
                    { amount: 1.00, name: "Pastanaga" },
                    { amount: 0.50, name: "Pebrot verd" },
                    { amount: 0.95, name: "Pebrot vermell" },
                    { amount: 4.00, name: "Iogurt" },
                    { amount: 0.50, name: "Poma" },
                    { amount: 0.50, name: "Plàtan" },
                    { amount: 0.50, name: "Peres" },
                    { amount: 0.50, name: "Taronges" },
                    { amount: 0.50, name: "Maduixes" },
                    { amount: 0.50, name: "Kiwi" },
                    { amount: 0.50, name: "Mandarines" },
                    { amount: 0.50, name: "Pinya" },
                    { amount: 0.50, name: "Mango" },
                    { amount: 0.50, name: "Pera" },
                    { amount: 0.50, name: "Cireres" },
                    { amount: 0.50, name: "Préssecs" },
                    { amount: 0.50, name: "Albercocs" },
                    { amount: 0.50, name: "Nectarines" },
                    { amount: 0.50, name: "Pressec de vinya" },
                    { amount: 0.50, name: "Poma àcida" },
                    { amount: 0.50, name: "Poma verda" },
                    { amount: 0.50, name: "Ceba" },
                    { amount: 0.50, name: "Ceba vermella" },
                    { amount: 0.50, name: "Ceba tendra" },
                    { amount: 0.50, name: "Ceba de figueres" },
                    { amount: 0.50, name: "Ceba de calçot" },
                    { amount: 0.50, name: "Pipes" },
                    { amount: 0.50, name: "Nous" },
                    { amount: 0.50, name: "Ametlles" },
                    { amount: 0.50, name: "Avellanes" },
                    { amount: 0.50, name: "Pinyons" },
                    { amount: 0.50, name: "Anacards" },
                    { amount: 0.50, name: "Cacauets" },
                    { amount: 0.50, name: "Pistatxos" },
                    { amount: 0.50, name: "Garrofons" },
                    { amount: 0.50, name: "Mongetes" },
                    { amount: 0.50, name: "Mongetes tendres" },
                ]
            }))
    }
}
