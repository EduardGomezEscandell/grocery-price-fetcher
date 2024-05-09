import React, { useState } from 'react'
import { Ingredient } from '../../State/State.js'

interface Props {
    style: React.CSSProperties;
    ingredient: Ingredient;
    onChange: () => void;
}

export default function RenderIngredient(pp: Props): JSX.Element {
    const def = Numbers.positive(pp.ingredient.need - pp.ingredient.have)
    const pks = Math.ceil(def / pp.ingredient.batch_size)

    const [storage, setStorage] = useState(pp.ingredient.have)
    const [deficit, setDeficit] = useState(def)
    const [packs, setMustBuy] = useState(pks)
    const [cost, setCost] = useState(pks * pp.ingredient.price)

    pp.ingredient.have = storage

    const titleStyle: React.CSSProperties = {
        width: '200px',
        textAlign: 'left',
        paddingLeft: '20px'
    }

    const numberStyle: React.CSSProperties = {
        width: '100px',
        textAlign: 'right',
        paddingRight: '20px'
    }

    const defaultBackground = pp.style.background
    const [rowStyle, setRowStyle] = useState(pp.style)

    return (
        <tr key={pp.ingredient.name}
            style={rowStyle}
            onMouseEnter={() => {
                setRowStyle({ ...rowStyle, background: 'lightblue' })
            }}

            onMouseLeave={() => {
                setRowStyle({ ...rowStyle, background: defaultBackground })
            }}
        >
            <td style={titleStyle}> {pp.ingredient.name}  </td>
            <td style={numberStyle}> <input
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
                    pp.onChange()
                }}
                
                datatype='number'
                style={{ width: '40px' }}
            /> </td>
            <td style={numberStyle}> {Numbers.round2(pp.ingredient.need)} </td>
            <td style={numberStyle}> {Numbers.round2(deficit)} </td>
            <td style={numberStyle}> {pp.ingredient.batch_size === 1 ? '' : Numbers.round2(pp.ingredient.batch_size)} </td>
            <td style={numberStyle}> {Numbers.int(packs)}</td>
            <td style={numberStyle}> {Numbers.asEuro(pp.ingredient.price)} </td>
            <td style={numberStyle}> {Numbers.asEuro(cost)} </td>
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

