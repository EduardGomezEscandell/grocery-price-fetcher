import { ShoppingNeeds } from '../../State/State'

export class NeedsEndpoint {
    protected path: string
    private auth: string

    constructor(auth: string, which: string) {
        this.path = '/api/shopping-needs/' + which
        this.auth = auth
    }

    Path(): string {
        return this.path
    }

    async GET(): Promise<ShoppingNeeds> {
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
            .then(data => ShoppingNeeds.fromJSON(data))
    }
}

export class MockNeedsEndpoint extends NeedsEndpoint {
    which: string
    constructor(auth: string, which: string) {
        super(auth, which)
        this.which = which
    }

    async GET(): Promise<ShoppingNeeds> {
        console.log(`GET from ${this.path}`)
        return new Promise(resolve => setTimeout(resolve, 100))
            .then(() => ShoppingNeeds.fromJSON({
                name: this.which,
                items: [
                    { product_id: 2, amount: 1.00, name: "Pastanaga" },
                    { product_id: 17, amount: 0.50, name: "Pebrot verd" },
                    { product_id: 99, amount: 0.95, name: "Pebrot vermell" },
                    { product_id: 3, amount: 4.00, name: "Iogurt" },
                    { product_id: 5, amount: 0.50, name: "Poma" },
                    { product_id: 6, amount: 0.50, name: "Plàtan" },
                    { product_id: 7, amount: 0.50, name: "Peres" },
                    { product_id: 8, amount: 0.50, name: "Taronges" },
                    { product_id: 9, amount: 0.50, name: "Maduixes" },
                    { product_id: 10, amount: 0.50, name: "Kiwi" },
                    { product_id: 11, amount: 0.50, name: "Mandarines" },
                    { product_id: 12, amount: 0.50, name: "Pinya" },
                    { product_id: 13, amount: 0.50, name: "Mango" },
                    { product_id: 14, amount: 0.50, name: "Pera" },
                    { product_id: 15, amount: 0.50, name: "Cireres" },
                    { product_id: 16, amount: 0.50, name: "Préssecs" },
                    { product_id: 1, amount: 0.50, name: "Albercocs" },
                    { product_id: 18, amount: 0.50, name: "Nectarines" },
                    { product_id: 19, amount: 0.50, name: "Pressec de vinya" },
                    { product_id: 20, amount: 0.50, name: "Poma àcida" },
                    { product_id: 21, amount: 0.50, name: "Poma verda" },
                    { product_id: 22, amount: 0.50, name: "Ceba" },
                    { product_id: 23, amount: 0.50, name: "Ceba vermella" },
                    { product_id: 24, amount: 0.50, name: "Ceba tendra" },
                    { product_id: 25, amount: 0.50, name: "Ceba de figueres" },
                    { product_id: 26, amount: 0.50, name: "Ceba de calçot" },
                    { product_id: 27, amount: 0.50, name: "Pipes" },
                    { product_id: 28, amount: 0.50, name: "Nous" },
                    { product_id: 29, amount: 0.50, name: "Ametlles" },
                    { product_id: 30, amount: 0.50, name: "Avellanes" },
                    { product_id: 31, amount: 0.50, name: "Pinyons" },
                    { product_id: 32, amount: 0.50, name: "Anacards" },
                    { product_id: 33, amount: 0.50, name: "Cacauets" },
                    { product_id: 34, amount: 0.50, name: "Pistatxos" },
                    { product_id: 35, amount: 0.50, name: "Garrofons" },
                    { product_id: 36, amount: 0.50, name: "Mongetes" },
                    { product_id: 37, amount: 0.50, name: "Mongetes tendres" },
                ]
            }))
    }
}
