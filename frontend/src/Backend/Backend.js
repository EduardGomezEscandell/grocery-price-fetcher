import { process } from 'react'

function SelectBackend() {
    if (process.env.REACT_APP_MOCK_BACKEND !== "") {
        return new MockBackend()
    }
    return new Backend()
}

class Backend {
    recipeEndpoint = '/api/recipes'
    menuEndpoint = '/api/menu'

    GetDishes() {
        return fetch(this.recipeEndpoint, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json'
            }
        })
    }

    GetMenu() {
        return fetch(this.menuEndpoint, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json'
            },
        })
    }

    PostMenu(menu) {
        return fetch(this.menuEndpoint, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(menu)
        })
    }
}

class MockBackend extends Backend {
    GetDishes() {
        console.log(`GET from ${this.recipeEndpoint}`)
        return Promise.resolve({
            json: () => Promise.resolve([
                "Amanida de cigrons", "Amanida de llenties", "Cigrons amb xoriço", "Llenties amb xoriço", "Arròs amb ou", "Arròs amb sofregit", "Cereals amb llet", "Flocs de civada", "Flocs de civada", "Torrada i suc de taronja", "Cafè", "Fajita", "Fruita", "Hamburguesa", "Proteïna amb acompanyament", "Vurguer", "Iogurt", "Iogurt vegà", "Macarrons amb sofregit", "Macarrons amb sofregit vegà", "Bastonets i hummus", "Pica-pica", "Puré de patata", "Sopa de galets", "Torrada d'alvocat", "Ramen", "Truita francesa", "Llom amb acompanyament", "Pollastre amb acompanyament", "Peix amb acompanyament", "Ravioli"
            ])
        })
    }

    GetMenu() {
        console.log(`GET from ${this.menuEndpoint}`)
        return Promise.resolve({
            json: () => Promise.resolve([
                { name: "Menú 1", menu: [{ name: "Dilluns", meals: [{ name: "Esmorzar", dishes: [{ name: "Torrada i suc de taronja", amount: 1 }, { name: "Flocs de civada", amount: 1 }] }, { name: "Dinar", dishes: [{ name: "Macarrons amb sofregit", amount: 1 }, { name: "Fruita", amount: 1 }, { name: "Cafè", amount: 1 }] }, { name: "Sopar", dishes: [{ name: "Ravioli", amount: 1 }, { name: "Iogurt", amount: 1 }] }] }, { name: "Dimarts", meals: [{ name: "Esmorzar", dishes: [{ name: "Torrada i suc de taronja", amount: 1 }, { name: "Flocs de civada", amount: 1 }] }, { name: "Dinar", dishes: [{ name: "Llenties amb xoriço", amount: 1 }, { name: "Pollastre amb acompanyament", amount: 1 }, { name: "Fruita", amount: 1 }, { name: "Cafè", amount: 1 }] }, { name: "Sopar", dishes: [{ name: "Macarrons amb sofregit", amount: 1 }, { name: "Iogurt", amount: 1 }] }] }, { name: "Dimecres", meals: [{ name: "Esmorzar", dishes: [{ name: "Torrada i suc de taronja", amount: 1 }] }, { name: "Dinar", dishes: [{ name: "Arròs amb sofregit", amount: 1 }, { name: "Fruita", amount: 1 }, { name: "Cafè", amount: 1 }] }, { name: "Sopar", dishes: [{ name: "Sopa de galets", amount: 2 }, { name: "Proteïna amb acompanyament", amount: 2 }, { name: "Iogurt", amount: 1 }, { name: "Iogurt vegà", amount: 1 }] }] }, { name: "Dijous", meals: [{ name: "Esmorzar", dishes: [{ name: "Torrada i suc de taronja", amount: 1 }, { name: "Flocs de civada", amount: 1 }] }, { name: "Dinar", dishes: [{ name: "Amanida de llenties", amount: 1 }, { name: "Llom amb acompanyament", amount: 1 }, { name: "Fruita", amount: 1 }, { name: "Cafè", amount: 1 }] }, { name: "Sopar", dishes: [{ name: "Arròs amb sofregit", amount: 1 }, { name: "Iogurt", amount: 1 }] }] }, { name: "Divendres", meals: [{ name: "Esmorzar", dishes: [{ name: "Torrada i suc de taronja", amount: 1 }] }, { name: "Dinar", dishes: [{ name: "Amanida de cigrons", amount: 1 }, { name: "Pollastre amb acompanyament", amount: 1 }, { name: "Fruita", amount: 1 }, { name: "Cafè", amount: 1 }] }, { name: "Sopar", dishes: [{ name: "Ravioli", amount: 1 }, { name: "Iogurt", amount: 1 }] }] }, { name: "Dissabte", meals: [{ name: "Esmorzar", dishes: [{ name: "Torrada i suc de taronja", amount: 1 }] }, { name: "Sopar", dishes: [{ name: "Ramen", amount: 2 }, { name: "Iogurt", amount: 1 }, { name: "Iogurt vegà", amount: 1 }] }] }, { name: "Diumenge", meals: [{ name: "Esmorzar", dishes: [{ name: "Cereals amb llet", amount: 2 }] }, { name: "Dinar", dishes: [{ name: "Pica-pica", amount: 2 }, { name: "Puré de patata", amount: 2 }, { name: "Proteïna amb acompanyament", amount: 2 }, { name: "Fruita", amount: 2 }, { name: "Cafè", amount: 1 }] }, { name: "Sopar", dishes: [{ name: "Sopa de galets", amount: 2 }, { name: "Proteïna amb acompanyament", amount: 2 }, { name: "Iogurt", amount: 1 }, { name: "Iogurt vegà", amount: 1 }] }] }] }
            ])
        })
    }

    PostMenu(menu) {
        console.log(`POST to ${this.menuEndpoint}:`)
        console.log(JSON.stringify(menu))
        return Promise.resolve({
            json: () => Promise.resolve([
                { amount: 4.00, unit_cost: 0.49, product: "Fruita" },
                { amount: 1.00, unit_cost: 0.17, product: "Pastanaga" },
                { amount: 0.50, unit_cost: 0.50, product: "Pebrot verd" },
                { amount: 0.95, unit_cost: 1.10, product: "Pebrot vermell" },
            ])
        })
    }
}

export default SelectBackend