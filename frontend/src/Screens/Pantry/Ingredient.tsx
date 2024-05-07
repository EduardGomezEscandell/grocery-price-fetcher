import React, {useState} from 'react'
import {Ingredient} from '../../State/State.js'

interface Props {
    ingredient: Ingredient;
    onChange: () => void;
}

export default function RenderIngredient(pp: Props): JSX.Element {
    const def = Numbers.positive(pp.ingredient.need - pp.ingredient.have)
    
    const [storage, setStorage] = useState(pp.ingredient.have)
    const [deficit, setDeficit] = useState(def)
    const [cost, setCost] = useState(def * pp.ingredient.price)

    pp.ingredient.have = storage

    return (
        <tr key={pp.ingredient.name}>
            <td> {pp.ingredient.name}  </td>
            <td> <input
                type="number"
                value={storage}
                onChange={(s) => {
                    const newStorage = Numbers.positive(parseFloat(s.target.value))
                    const newDeficit = Numbers.positive(pp.ingredient.need - newStorage)
                    setStorage(newStorage)
                    setDeficit(newDeficit)
                    setCost(newDeficit * pp.ingredient.price)
                    pp.onChange()
                }}
                datatype='number'
                style={{ width: '40px' }}
            /> </td>
            <td> {pp.ingredient.need} </td>
            <td> {deficit} </td>
            <td> {Numbers.asEuro(pp.ingredient.price)} </td>
            <td> {Numbers.asEuro(cost)} </td>
        </tr>
    )
}

export class Numbers {
    static positive(x: number): number {
        return x >= 0 ? x : 0
    }
    
    static asEuro(x: number): string {
        return x.toFixed(2) + ' €'
    }
}

