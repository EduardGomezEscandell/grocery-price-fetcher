import { Menu } from '../../State/State.tsx'
import { ShoppingList } from '../../State/State.tsx'

export class MenuEndpoint {
    protected static path: string = '/api/menu'

    static Path(): string {
        return this.path
    }

    async GET(): Promise<Menu[]> {
        return fetch(MenuEndpoint.path, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json',
            }
        })
            .then(response => response.json())
            .then(data => data.map(m => Menu.fromJSON(m)))
    }

    async POST(menu: Menu): Promise<ShoppingList> {
        return fetch(MenuEndpoint.path, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json',
            },
            body: menu.toJSON()
        })
            .then(response => response.json())
            .then(data => ShoppingList.fromJSON(data))
    }
}

export class MockMenuEndpoint extends MenuEndpoint {
    async GET(): Promise<Menu[]> {
        console.log(`GET from ${MockMenuEndpoint.path}`)
        return Promise.resolve(
            [
                Menu.fromJSON({ name: "Menú 1", days: [{ name: "Dilluns", meals: [{ name: "Esmorzar", dishes: [{ name: "Torrada i suc de taronja", amount: 1 }, { name: "Flocs de civada", amount: 1 }] }, { name: "Dinar", dishes: [{ name: "Macarrons amb sofregit", amount: 1 }, { name: "Fruita", amount: 1 }, { name: "Cafè", amount: 1 }] }, { name: "Sopar", dishes: [{ name: "Ravioli", amount: 1 }, { name: "Iogurt", amount: 1 }] }] }, { name: "Dimarts", meals: [{ name: "Esmorzar", dishes: [{ name: "Torrada i suc de taronja", amount: 1 }, { name: "Flocs de civada", amount: 1 }] }, { name: "Dinar", dishes: [{ name: "Llenties amb xoriço", amount: 1 }, { name: "Pollastre amb acompanyament", amount: 1 }, { name: "Fruita", amount: 1 }, { name: "Cafè", amount: 1 }] }, { name: "Sopar", dishes: [{ name: "Macarrons amb sofregit", amount: 1 }, { name: "Iogurt", amount: 1 }] }] }, { name: "Dimecres", meals: [{ name: "Esmorzar", dishes: [{ name: "Torrada i suc de taronja", amount: 1 }] }, { name: "Dinar", dishes: [{ name: "Arròs amb sofregit", amount: 1 }, { name: "Fruita", amount: 1 }, { name: "Cafè", amount: 1 }] }, { name: "Sopar", dishes: [{ name: "Sopa de galets", amount: 2 }, { name: "Proteïna amb acompanyament", amount: 2 }, { name: "Iogurt", amount: 1 }, { name: "Iogurt vegà", amount: 1 }] }] }, { name: "Dijous", meals: [{ name: "Esmorzar", dishes: [{ name: "Torrada i suc de taronja", amount: 1 }, { name: "Flocs de civada", amount: 1 }] }, { name: "Dinar", dishes: [{ name: "Amanida de llenties", amount: 1 }, { name: "Llom amb acompanyament", amount: 1 }, { name: "Fruita", amount: 1 }, { name: "Cafè", amount: 1 }] }, { name: "Sopar", dishes: [{ name: "Arròs amb sofregit", amount: 1 }, { name: "Iogurt", amount: 1 }] }] }, { name: "Divendres", meals: [{ name: "Esmorzar", dishes: [{ name: "Torrada i suc de taronja", amount: 1 }] }, { name: "Dinar", dishes: [{ name: "Amanida de cigrons", amount: 1 }, { name: "Pollastre amb acompanyament", amount: 1 }, { name: "Fruita", amount: 1 }, { name: "Cafè", amount: 1 }] }, { name: "Sopar", dishes: [{ name: "Ravioli", amount: 1 }, { name: "Iogurt", amount: 1 }] }] }, { name: "Dissabte", meals: [{ name: "Esmorzar", dishes: [{ name: "Torrada i suc de taronja", amount: 1 }] }, { name: "Sopar", dishes: [{ name: "Ramen", amount: 2 }, { name: "Iogurt", amount: 1 }, { name: "Iogurt vegà", amount: 1 }] }] }, { name: "Diumenge", meals: [{ name: "Esmorzar", dishes: [{ name: "Cereals amb llet", amount: 2 }] }, { name: "Dinar", dishes: [{ name: "Pica-pica", amount: 2 }, { name: "Puré de patata", amount: 2 }, { name: "Proteïna amb acompanyament", amount: 2 }, { name: "Fruita", amount: 2 }, { name: "Cafè", amount: 1 }] }, { name: "Sopar", dishes: [{ name: "Sopa de galets", amount: 2 }, { name: "Proteïna amb acompanyament", amount: 2 }, { name: "Iogurt", amount: 1 }, { name: "Iogurt vegà", amount: 1 }] }] }] }),
                Menu.fromJSON({ name: "Dummy menu" })
            ]
        )
    }

    async POST(menu: Menu): Promise<ShoppingList> {
        console.log(`POST to ${MockMenuEndpoint.path}:`)
        console.log(menu.toJSON()) // Ensure toJSON is called without errors
        return Promise.resolve(ShoppingList.fromJSON([
            { batch_size: 1, need: 1.00, price: 0.17, product: "Pastanaga" },
            { batch_size: 1, need: 0.50, price: 0.50, product: "Pebrot verd" },
            { batch_size: 1, need: 0.95, price: 1.10, product: "Pebrot vermell" },
            { batch_size: 4, need: 4.00, price: 1.00, product: "Iogurt" },
            { batch_size: 1, need: 0.50, price: 0.50, product: "Poma" },
            { batch_size: 1, need: 0.50, price: 0.50, product: "Plàtan" },
            { batch_size: 1, need: 0.50, price: 0.50, product: "Peres" },
            { batch_size: 1, need: 0.50, price: 0.50, product: "Taronges" },
            { batch_size: 1, need: 0.50, price: 0.50, product: "Maduixes" },
            { batch_size: 1, need: 0.50, price: 0.50, product: "Kiwi" },
            { batch_size: 1, need: 0.50, price: 0.50, product: "Mandarines" },
            { batch_size: 1, need: 0.50, price: 0.50, product: "Pinya" },
            { batch_size: 1, need: 0.50, price: 0.50, product: "Mango" },
            { batch_size: 1, need: 0.50, price: 0.50, product: "Pera" },
            { batch_size: 1, need: 0.50, price: 0.50, product: "Cireres" },
            { batch_size: 1, need: 0.50, price: 0.50, product: "Préssecs" },
            { batch_size: 1, need: 0.50, price: 0.50, product: "Albercocs" },
            { batch_size: 1, need: 0.50, price: 0.50, product: "Nectarines" },
            { batch_size: 1, need: 0.50, price: 0.50, product: "Pressec de vinya" },
            { batch_size: 1, need: 0.50, price: 0.50, product: "Poma àcida" },
            { batch_size: 1, need: 0.50, price: 0.50, product: "Poma verda" },
            { batch_size: 1, need: 0.50, price: 0.50, product: "Ceba" },
            { batch_size: 1, need: 0.50, price: 0.50, product: "Ceba vermella" },
            { batch_size: 1, need: 0.50, price: 0.50, product: "Ceba tendra" },
            { batch_size: 1, need: 0.50, price: 0.50, product: "Ceba de figueres" },
            { batch_size: 1, need: 0.50, price: 0.50, product: "Ceba de calçot" },
            { batch_size: 1, need: 0.50, price: 0.50, product: "Pipes" },
            { batch_size: 1, need: 0.50, price: 0.50, product: "Nous" },
            { batch_size: 1, need: 0.50, price: 0.50, product: "Ametlles" },
            { batch_size: 1, need: 0.50, price: 0.50, product: "Avellanes" },
            { batch_size: 1, need: 0.50, price: 0.50, product: "Pinyons" },
            { batch_size: 1, need: 0.50, price: 0.50, product: "Anacards" },
            { batch_size: 1, need: 0.50, price: 0.50, product: "Cacauets" },
            { batch_size: 1, need: 0.50, price: 0.50, product: "Pistatxos" },
            { batch_size: 1, need: 0.50, price: 0.50, product: "Garrofons" },
            { batch_size: 1, need: 0.50, price: 0.50, product: "Mongetes" },
            { batch_size: 1, need: 0.50, price: 0.50, product: "Mongetes tendres" },
        ]))
    }
}
