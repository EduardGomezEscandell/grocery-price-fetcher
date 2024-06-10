import { Menu } from '../../State/State.tsx'

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

    async PUT(menu: Menu): Promise<void> {
        return fetch(MenuEndpoint.path, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json',
            },
            body: menu.toJSON()
        }).then(() => { })
    }
}

export class MockMenuEndpoint extends MenuEndpoint {
    async GET(): Promise<Menu[]> {
        console.log(`GET from ${MockMenuEndpoint.path}`)
        return new Promise(resolve => setTimeout(resolve, 100))
            .then(() =>
                [
                    Menu.fromJSON({ name: "default", days: [{ name: "Dilluns", meals: [{ name: "Esmorzar", dishes: [{ name: "Torrada i suc de taronja", amount: 1 }, { name: "Flocs de civada", amount: 1 }] }, { name: "Dinar", dishes: [{ name: "Macarrons amb sofregit", amount: 1 }, { name: "Fruita", amount: 1 }, { name: "Cafè", amount: 1 }] }, { name: "Sopar", dishes: [{ name: "Ravioli", amount: 1 }, { name: "Iogurt", amount: 1 }] }] }, { name: "Dimarts", meals: [{ name: "Esmorzar", dishes: [{ name: "Torrada i suc de taronja", amount: 1 }, { name: "Flocs de civada", amount: 1 }] }, { name: "Dinar", dishes: [{ name: "Llenties amb xoriço", amount: 1 }, { name: "Pollastre amb acompanyament", amount: 1 }, { name: "Fruita", amount: 1 }, { name: "Cafè", amount: 1 }] }, { name: "Sopar", dishes: [{ name: "Macarrons amb sofregit", amount: 1 }, { name: "Iogurt", amount: 1 }] }] }, { name: "Dimecres", meals: [{ name: "Esmorzar", dishes: [{ name: "Torrada i suc de taronja", amount: 1 }] }, { name: "Dinar", dishes: [{ name: "Arròs amb sofregit", amount: 1 }, { name: "Fruita", amount: 1 }, { name: "Cafè", amount: 1 }] }, { name: "Sopar", dishes: [{ name: "Sopa de galets", amount: 2 }, { name: "Proteïna amb acompanyament", amount: 2 }, { name: "Iogurt", amount: 1 }, { name: "Iogurt vegà", amount: 1 }] }] }, { name: "Dijous", meals: [{ name: "Esmorzar", dishes: [{ name: "Torrada i suc de taronja", amount: 1 }, { name: "Flocs de civada", amount: 1 }] }, { name: "Dinar", dishes: [{ name: "Amanida de llenties", amount: 1 }, { name: "Llom amb acompanyament", amount: 1 }, { name: "Fruita", amount: 1 }, { name: "Cafè", amount: 1 }] }, { name: "Sopar", dishes: [{ name: "Arròs amb sofregit", amount: 1 }, { name: "Iogurt", amount: 1 }] }] }, { name: "Divendres", meals: [{ name: "Esmorzar", dishes: [{ name: "Torrada i suc de taronja", amount: 1 }] }, { name: "Dinar", dishes: [{ name: "Amanida de cigrons", amount: 1 }, { name: "Pollastre amb acompanyament", amount: 1 }, { name: "Fruita", amount: 1 }, { name: "Cafè", amount: 1 }] }, { name: "Sopar", dishes: [{ name: "Ravioli", amount: 1 }, { name: "Iogurt", amount: 1 }] }] }, { name: "Dissabte", meals: [{ name: "Esmorzar", dishes: [{ name: "Torrada i suc de taronja", amount: 1 }] }, { name: "Dinar" }, { name: "Sopar", dishes: [{ name: "Ramen", amount: 2 }, { name: "Iogurt", amount: 1 }, { name: "Iogurt vegà", amount: 1 }] }] }, { name: "Diumenge", meals: [{ name: "Esmorzar", dishes: [{ name: "Cereals amb llet", amount: 2 }] }, { name: "Dinar", dishes: [{ name: "Pica-pica", amount: 2 }, { name: "Puré de patata", amount: 2 }, { name: "Proteïna amb acompanyament", amount: 2 }, { name: "Fruita", amount: 2 }, { name: "Cafè", amount: 1 }] }, { name: "Sopar", dishes: [{ name: "Sopa de galets", amount: 2 }, { name: "Proteïna amb acompanyament", amount: 2 }, { name: "Iogurt", amount: 1 }, { name: "Iogurt vegà", amount: 1 }] }] }] }),
                    Menu.fromJSON({ name: "Dummy menu" })
                ]
            )
    }

    async PUT(menu: Menu): Promise<void> {
        console.log(`PUT to ${MockMenuEndpoint.path}:`)
        console.log(menu.toJSON().substring(0, 30), '...') // Ensure toJSON is called without errors
        return new Promise(resolve => setTimeout(resolve, 100))
    }
}
