import { Dish } from "../../State/State"

export class DishesEndpoint {
    protected static path = '/api/recipes'

    Path(): string {
        return DishesEndpoint.path
    }

    async GET(): Promise<Dish[]> {
        return fetch(DishesEndpoint.path, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json'
            }
        })
            .then(r => r.ok ? r : Promise.reject(r))
            .then(r => r.json())
            .then(parse)
    }
}

function parse(objs: any[]): Dish[] {
    return objs.map(obj => new Dish(obj.id, obj.name, obj.amount))
}

export class MockDishesEndpoint extends DishesEndpoint {
    async GET(): Promise<Dish[]> {
        console.log(`GET from ${MockDishesEndpoint.path}`)
        return new Promise(resolve => setTimeout(resolve, 100))
            .then(() => [
                {"id": 1, "name": "Amanida de cigrons"},
                {"id": 2, "name": "Amanida de llenties"},
                {"id": 3, "name": "Cigrons amb xoriço"},
                {"id": 4, "name": "Llenties amb xoriço"},
                {"id": 5, "name": "Arròs amb ou"},
                {"id": 6, "name": "Arròs amb sofregit"},
                {"id": 7, "name": "Cereals amb llet"},
                {"id": 8, "name": "Flocs de civada"},
                {"id": 9, "name": "Torrada i suc de taronja"},
                {"id": 10, "name": "Cafè"},
                {"id": 11, "name": "Fajita"},
                {"id": 12, "name": "Fruita"},
                {"id": 13, "name": "Hamburguesa"},
                {"id": 14, "name": "Proteïna amb acompanyament"},
                {"id": 15, "name": "Vurguer"},
                {"id": 16, "name": "Iogurt"},
                {"id": 17, "name": "Iogurt vegà"},
                {"id": 18, "name": "Macarrons amb sofregit"},
                {"id": 19, "name": "Macarrons amb sofregit vegà"},
                {"id": 20, "name": "Bastonets i hummus"},
                {"id": 21, "name": "Pica-pica"},
                {"id": 22, "name": "Puré de patata"},
                {"id": 23, "name": "Sopa de galets"},
                {"id": 24, "name": "Torrada d'alvocat"},
                {"id": 25, "name": "Ramen"},
                {"id": 26, "name": "Truita francesa"},
                {"id": 27, "name": "Llom amb acompanyament"},
                {"id": 28, "name": "Pollastre amb acompanyament"},
                {"id": 29, "name": "Peix amb acompanyament"},
                {"id": 30, "name": "Ravioli"}
            ])
            .then(parse)
    }
}