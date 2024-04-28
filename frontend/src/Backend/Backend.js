import {process} from 'react'

function SelectBackend() {
    if (process.env.REACT_APP_MOCK_BACKEND !== "") {
        return new MockBackend()
    }
    return new Backend()
}

class Backend {
    recipeEndpoint = '/api/recipes'
    menuEndpoint = '/api/menu'

    fetchDishes() {
        return fetch(this.recipeEndpoint, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json'
            }
        })
    }

    postMenu(menu) {
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
    fetchDishes() {
        console.log(`GET from ${this.menuEndpoint}`)
        return Promise.resolve({
            json: () => Promise.resolve([
                "Pasta",
                "Pizza quattro formaggi",
                "Tiramisu",
            ])
        })
    }

    postMenu(menu) {
        console.log(`POST to ${this.menuEndpoint}:`)
        console.log(JSON.stringify(menu))
        return Promise.resolve({
            json: () => Promise.resolve([
                {amount: 4.00, unit_cost: 0.49, product: "Fruita"},
                {amount: 1.00, unit_cost: 0.17, product: "Pastanaga"},
                {amount: 0.50, unit_cost: 0.50, product: "Pebrot verd"},
                {amount: 0.95, unit_cost: 1.10, product: "Pebrot vermell"},
            ])
        })
    }
}

export default SelectBackend