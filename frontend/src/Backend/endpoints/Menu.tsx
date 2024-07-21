import { Menu } from '../../State/State'

export class MenuEndpoint {
    protected path: string

    constructor(which: string) {
        this.path = `/api/menu/${which}`
    }

    Path(): string {
        return this.path
    }

    async GET(): Promise<Menu> {
        return fetch(this.path, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json',
            }
        })
            .then(r => r.ok ? r : Promise.reject(r))
            .then(r => r.json())
            .then(data => Menu.fromJSON(data))
    }

    async PUT(menu: Menu): Promise<void> {
        return fetch(this.path, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json',
            },
            body: menu.toJSON()
        })
            .then(r => r.ok ? r : Promise.reject(r))
            .then(() => { })
    }
}

export class MockMenuEndpoint extends MenuEndpoint {
    which: string
    constructor(which: string) {
        super(which)
        this.which = which
    }

    async GET(): Promise<Menu> {
        console.log(`GET from ${this.path}`)
        return new Promise(resolve => setTimeout(resolve, 100))
            .then(() => Menu.fromJSON(
                {
                    name: "default", days: [
                        { name: "Dilluns", meals: [{ name: "Esmorzar", dishes: [{ id: 1, name: "Torrada i suc de taronja", amount: 1 }, { id: 3, name: "Flocs de civada", amount: 1 }] }, { name: "Dinar", dishes: [{ id: 6, name: "Macarrons amb sofregit", amount: 1 }, { id: 7, name: "Fruita", amount: 1 }, { id: 8, name: "Cafè", amount: 1 }] }, { name: "Sopar", dishes: [{ id: 15, name: "Ravioli", amount: 1 }, { id: 16, name: "Iogurt", amount: 1 }] }] },
                        { name: "Dimarts", meals: [{ name: "Esmorzar", dishes: [{ id: 1, name: "Torrada i suc de taronja", amount: 1 }, { id: 3, name: "Flocs de civada", amount: 1 }] }, { name: "Dinar", dishes: [{ id: 20, name: "Llenties amb xoriço", amount: 1 }, { id: 9, name: "Pollastre amb acompanyament", amount: 1 }, { id: 7, name: "Fruita", amount: 1 }, { id: 8, name: "Cafè", amount: 1 }] }, { name: "Sopar", dishes: [{ id: 6, name: "Macarrons amb sofregit", amount: 1 }, { id: 16, name: "Iogurt", amount: 1 }] }] },
                        { name: "Dimecres", meals: [{ name: "Esmorzar", dishes: [{ id: 1, name: "Torrada i suc de taronja", amount: 1 }] }, { name: "Dinar", dishes: [{ id: 34, name: "Arròs amb sofregit", amount: 1 }, { id: 7, name: "Fruita", amount: 1 }, { id: 8, name: "Cafè", amount: 1 }] }, { name: "Sopar", dishes: [{ id: 31, name: "Sopa de galets", amount: 2 }, { id: 30, name: "Proteïna amb acompanyament", amount: 2 }, { id: 16, name: "Iogurt", amount: 1 }, { id: 22, name: "Iogurt vegà", amount: 1 }] }] },
                        { name: "Dijous", meals: [{ name: "Esmorzar", dishes: [{ id: 1, name: "Torrada i suc de taronja", amount: 1 }, { id: 3, name: "Flocs de civada", amount: 1 }] }, { name: "Dinar", dishes: [{ id: 40, name: "Amanida de llenties", amount: 1 }, { id: 38, name: "Llom amb acompanyament", amount: 1 }, { id: 7, name: "Fruita", amount: 1 }, { id: 8, name: "Cafè", amount: 1 }] }, { name: "Sopar", dishes: [{ id: 34, name: "Arròs amb sofregit", amount: 1 }, { id: 16, name: "Iogurt", amount: 1 }] }] },
                        { name: "Divendres", meals: [{ name: "Esmorzar", dishes: [{ id: 1, name: "Torrada i suc de taronja", amount: 1 }] }, { name: "Dinar", dishes: [{ id: 100, name: "Amanida de cigrons", amount: 1 }, { id: 9, name: "Pollastre amb acompanyament", amount: 1 }, { id: 7, name: "Fruita", amount: 1 }, { id: 8, name: "Cafè", amount: 1 }] }, { name: "Sopar", dishes: [{ id: 15, name: "Ravioli", amount: 1 }, { id: 16, name: "Iogurt", amount: 1 }] }] },
                        { name: "Dissabte", meals: [{ name: "Esmorzar", dishes: [{ id: 1, name: "Torrada i suc de taronja", amount: 1 }] }, { name: "Dinar" }, { name: "Sopar", dishes: [{ id: 55, name: "Ramen", amount: 2 }, { id: 16, name: "Iogurt", amount: 1 }, { id: 22, name: "Iogurt vegà", amount: 1 }] }] },
                        { name: "Diumenge", meals: [{ name: "Esmorzar", dishes: [{ id: 2, name: "Cereals amb llet", amount: 2 }] }, { name: "Dinar", dishes: [{ id: 56, name: "Pica-pica", amount: 2 }, { id: 60, name: "Puré de patata", amount: 2 }, { id: 30, name: "Proteïna amb acompanyament", amount: 2 }, { id: 7, name: "Fruita", amount: 2 }, { id: 8, name: "Cafè", amount: 1 }] }, { name: "Sopar", dishes: [{ id: 31, name: "Sopa de galets", amount: 2 }, { id: 30, name: "Proteïna amb acompanyament", amount: 2 }, { id: 16, name: "Iogurt", amount: 1 }, { id: 22, name: "Iogurt vegà", amount: 1 }] }] }
                    ]
                }
            ))
    }

    async PUT(menu: Menu): Promise<void> {
        console.log(`PUT to ${this.path}:`)
        console.log(menu.toJSON().substring(0, 30), '...') // Ensure toJSON is called without errors
        return new Promise(resolve => setTimeout(resolve, 100))
    }
}
