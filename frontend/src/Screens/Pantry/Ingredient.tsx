import React, { useState } from 'react'
import { Ingredient } from '../../State/State.js'

interface Props {
    ingredient: Ingredient;
    onChange: (newHave: number) => void;
    id: string;
}

export default function RenderIngredient(pp: Props): JSX.Element {
    const def = Numbers.positive(pp.ingredient.need - pp.ingredient.have)
    const pks = Math.ceil(def / pp.ingredient.batch_size)

    const [storage, setStorage] = useState(pp.ingredient.have)
    const [deficit, setDeficit] = useState(def)
    const [packs, setMustBuy] = useState(pks)
    const [cost, setCost] = useState(pks * pp.ingredient.price)

    pp.ingredient.have = storage
    const defaultID = pp.id
    const [ID, setID] = useState(defaultID)

    return (
        <tr key={pp.ingredient.name}
            id={ID}
            onMouseEnter={() => setID('mouseover')}
            onMouseLeave={() => setID(defaultID)}
        >
            <td className='Label' key='name'> {pp.ingredient.name}  </td>
            <td className='Select' key='have'> <input
                type="number"
                value={storage}
                onChange={(s) => {
                    const newStorage = Numbers.positive(parseFloat(s.target.value))
                    const newDeficit = Numbers.positive(pp.ingredient.need - newStorage)
                    const newPackCount = Math.ceil(newDeficit / pp.ingredient.batch_size)
                    setStorage(newStorage)
                    setDeficit(newDeficit)
                    setMustBuy(newPackCount)
                    setCost(newPackCount * pp.ingredient.price)
                    pp.onChange(newStorage)
                }}
                
                datatype='number'
                style={{ width: '40px' }}
            /> </td>
            <td className='Number' key='need'> {Numbers.round2(pp.ingredient.need)} </td>
            <td className='Number' key='miss'> {Numbers.round2(deficit)} </td>
            <td className='Number' key='batch-size'> {pp.ingredient.batch_size === 1 ? '' : Numbers.round2(pp.ingredient.batch_size)} </td>
            <td className='Number' key='batch-count'> {Numbers.int(packs)}</td>
            <td className='Number' key='pack-price'> {Numbers.asEuro(pp.ingredient.price)} </td>
            <td className='Number' key='price-total'> {Numbers.asEuro(cost)} </td>
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

    static roundUpTo(x: number, divisor: number): number {
        return Math.ceil(x / divisor) * divisor
    }

    static int(x: number): string {
        return x.toFixed(0)
    }

    static round2(x: number): string {
        return x.toFixed(2)
    }
}

