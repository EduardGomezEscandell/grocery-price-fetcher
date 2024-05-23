import Backend from '../../Backend/Backend.ts'
import { Pantry, ShoppingList, State } from '../../State/State.tsx'

export default async function DownloadPantry(backend: Backend, state: State): Promise<void> {
    const [shopping, pantry] = await Promise.all([
        backend
            .Menu()
            .POST(state.menu),
        backend
            .Pantry()
            .GET()
            .then((p: Pantry[]) => p.length > 0 ? p[0] : new Pantry())
    ])
    const s = merge(pantry, shopping)
    state.shoppingList = s
}

function merge(pantry: Pantry, shoppingList: ShoppingList): ShoppingList {
    pantry.contents.sort((a, b) => a.name.localeCompare(b.name))
    shoppingList.ingredients.sort((a, b) => a.name.localeCompare(b.name))

    var i = 0
    var j = 0
    while (i < pantry.contents.length && j < shoppingList.ingredients.length) {
        const p = pantry.contents[i]
        const s = shoppingList.ingredients[j]


        switch (p.name.localeCompare(s.name)) {
            case 0:
                shoppingList.ingredients[j].have = p.have
                i++
                j++
                break
            case -1:
                i++
                break
            case 1:
                j++
                break
        }
    }

    return shoppingList
}