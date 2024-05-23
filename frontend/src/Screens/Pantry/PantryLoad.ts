import Backend from '../../Backend/Backend.ts'
import { Pantry, ShoppingNeeds, State } from '../../State/State.tsx'

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
    state.inNeed = merge(pantry, shopping)
}

function merge(have: Pantry, need: ShoppingNeeds): ShoppingNeeds {
    have.contents.sort((a, b) => a.name.localeCompare(b.name))
    need.ingredients.sort((a, b) => a.name.localeCompare(b.name))

    var i = 0
    var j = 0
    while (i < have.contents.length && j < need.ingredients.length) {
        const p = have.contents[i]
        const s = need.ingredients[j]


        switch (p.name.localeCompare(s.name)) {
            case 0:
                need.ingredients[j].have = p.have
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

    return need
}