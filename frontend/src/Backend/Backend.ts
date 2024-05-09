import { process } from 'react'
import { Menu, ShoppingList } from '../State/State.tsx'

class Backend {
    static New(): Backend {
        if (process.env.REACT_APP_MOCK_BACKEND !== "") {
            return new MockBackend()
        }
        return new Backend()
    }

    recipeEndpoint = '/api/recipes'
    menuEndpoint = '/api/menu'

    async GetDishes(): Promise<string[]> {
        return fetch(this.recipeEndpoint, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json'
            }
        }).then(response => response.json())
    }

    async GetMenu(): Promise<Menu[]> {
        return fetch(this.menuEndpoint, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json'
            },
        }).then(response => response.json())
          .then(data => data.map(m => Menu.fromJson(m)))
    }

    async PostMenu(menu: Menu): Promise<ShoppingList> {
        return fetch(this.menuEndpoint, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json',
            },
            body: menu.toJSON()
        }).then(response => response.json())
          .then(data => ShoppingList.fromJSON(data))
    }
}

class MockBackend extends Backend {
    async GetDishes(): Promise<string[]> {
        console.log(`GET from ${this.recipeEndpoint}`)
        return Promise.resolve([
            "Amanida de cigrons", "Amanida de llenties", "Cigrons amb xoriço", "Llenties amb xoriço", "Arròs amb ou", "Arròs amb sofregit", "Cereals amb llet", "Flocs de civada", "Flocs de civada", "Torrada i suc de taronja", "Cafè", "Fajita", "Fruita", "Hamburguesa", "Proteïna amb acompanyament", "Vurguer", "Iogurt", "Iogurt vegà", "Macarrons amb sofregit", "Macarrons amb sofregit vegà", "Bastonets i hummus", "Pica-pica", "Puré de patata", "Sopa de galets", "Torrada d'alvocat", "Ramen", "Truita francesa", "Llom amb acompanyament", "Pollastre amb acompanyament", "Peix amb acompanyament", "Ravioli"
        ])
    }

    async GetMenu(): Promise<Menu[]> {
        console.log(`GET from ${this.menuEndpoint}`)
        return Promise.resolve(
            [
                Menu.fromJson({name: "Menú 1", days: [{ name: "Dilluns", meals: [{ name: "Esmorzar", dishes: [{ name: "Torrada i suc de taronja", amount: 1 }, { name: "Flocs de civada", amount: 1 }] }, { name: "Dinar", dishes: [{ name: "Macarrons amb sofregit", amount: 1 }, { name: "Fruita", amount: 1 }, { name: "Cafè", amount: 1 }] }, { name: "Sopar", dishes: [{ name: "Ravioli", amount: 1 }, { name: "Iogurt", amount: 1 }] }] }, { name: "Dimarts", meals: [{ name: "Esmorzar", dishes: [{ name: "Torrada i suc de taronja", amount: 1 }, { name: "Flocs de civada", amount: 1 }] }, { name: "Dinar", dishes: [{ name: "Llenties amb xoriço", amount: 1 }, { name: "Pollastre amb acompanyament", amount: 1 }, { name: "Fruita", amount: 1 }, { name: "Cafè", amount: 1 }] }, { name: "Sopar", dishes: [{ name: "Macarrons amb sofregit", amount: 1 }, { name: "Iogurt", amount: 1 }] }] }, { name: "Dimecres", meals: [{ name: "Esmorzar", dishes: [{ name: "Torrada i suc de taronja", amount: 1 }] }, { name: "Dinar", dishes: [{ name: "Arròs amb sofregit", amount: 1 }, { name: "Fruita", amount: 1 }, { name: "Cafè", amount: 1 }] }, { name: "Sopar", dishes: [{ name: "Sopa de galets", amount: 2 }, { name: "Proteïna amb acompanyament", amount: 2 }, { name: "Iogurt", amount: 1 }, { name: "Iogurt vegà", amount: 1 }] }] }, { name: "Dijous", meals: [{ name: "Esmorzar", dishes: [{ name: "Torrada i suc de taronja", amount: 1 }, { name: "Flocs de civada", amount: 1 }] }, { name: "Dinar", dishes: [{ name: "Amanida de llenties", amount: 1 }, { name: "Llom amb acompanyament", amount: 1 }, { name: "Fruita", amount: 1 }, { name: "Cafè", amount: 1 }] }, { name: "Sopar", dishes: [{ name: "Arròs amb sofregit", amount: 1 }, { name: "Iogurt", amount: 1 }] }] }, { name: "Divendres", meals: [{ name: "Esmorzar", dishes: [{ name: "Torrada i suc de taronja", amount: 1 }] }, { name: "Dinar", dishes: [{ name: "Amanida de cigrons", amount: 1 }, { name: "Pollastre amb acompanyament", amount: 1 }, { name: "Fruita", amount: 1 }, { name: "Cafè", amount: 1 }] }, { name: "Sopar", dishes: [{ name: "Ravioli", amount: 1 }, { name: "Iogurt", amount: 1 }] }] }, { name: "Dissabte", meals: [{ name: "Esmorzar", dishes: [{ name: "Torrada i suc de taronja", amount: 1 }] }, { name: "Sopar", dishes: [{ name: "Ramen", amount: 2 }, { name: "Iogurt", amount: 1 }, { name: "Iogurt vegà", amount: 1 }] }] }, { name: "Diumenge", meals: [{ name: "Esmorzar", dishes: [{ name: "Cereals amb llet", amount: 2 }] }, { name: "Dinar", dishes: [{ name: "Pica-pica", amount: 2 }, { name: "Puré de patata", amount: 2 }, { name: "Proteïna amb acompanyament", amount: 2 }, { name: "Fruita", amount: 2 }, { name: "Cafè", amount: 1 }] }, { name: "Sopar", dishes: [{ name: "Sopa de galets", amount: 2 }, { name: "Proteïna amb acompanyament", amount: 2 }, { name: "Iogurt", amount: 1 }, { name: "Iogurt vegà", amount: 1 }] }] }]}),
                Menu.fromJson({name: "Dummy menu"})
            ]
        )
    }

    async PostMenu(menu: Menu): Promise<ShoppingList> {
        console.log(`POST to ${this.menuEndpoint}:`)
        console.log(menu.toJSON()) // Ensure toJSON is called without errors
        return Promise.resolve(ShoppingList.fromJSON([
            { batch_size: 1, need: 1.00, price: 0.17, product: "Pastanaga" },
            { batch_size: 1, need: 0.50, price: 0.50, product: "Pebrot verd" },
            { batch_size: 1, need: 0.95, price: 1.10, product: "Pebrot vermell" },
            { batch_size: 4, need: 4.00, price: 1.00, product: "Iogurt" },
        ]))
    }
}

export default Backend