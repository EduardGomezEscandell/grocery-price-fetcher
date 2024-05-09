import React, { useEffect } from 'react'
import Backend from '../../Backend/Backend'
import { Pantry, ShoppingList, State } from '../../State/State.tsx'

interface Props {
    backend: Backend
    state: State
    onComplete: () => void
}

export default function PantryLoad(pp: Props) {
    useEffect(() => {
        Promise.all([
            pp.backend
                .Menu()
                .POST(pp.state.menu),
            pp.backend
                .Pantry()
                .GET()
                .then((p: Pantry[]) => p.length > 0 ?  p[0] : new Pantry())
        ])
            .then(([shopping, pantry]) => merge(pantry, shopping))
            .then(s => pp.state.shoppingList = s)
            .finally(() => pp.onComplete())
    })

    return (
        <p>Loading...</p>
    )

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