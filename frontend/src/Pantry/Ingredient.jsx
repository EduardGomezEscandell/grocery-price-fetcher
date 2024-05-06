import React, { useState } from 'react'

function Ingredient(pp) {
    const item = pp.ingredient

    const [have, setHave] = useState(0)
    const [deficit, setDeficit] = useState(item.need)
    const [cost, setCost] = useState(item.unit_cost * item.need)

    const updateHave = (e) => {
        let have = parseFloat(e.target.value)
        if (isNaN(have)) {
            have = 0
        } else if (have < 0) {
            have = 0
        }
        setHave(have)
        const deficit = clampPositive(item.need - have)
        setDeficit(deficit)
        setCost(deficit * item.unit_cost)
        pp.onChange(have)
    }

    return (
        <tr key={item.name} >
            <td> {item.name}  </td>
            <td> <input
                type="number"
                value={have}
                onChange={updateHave}
                datatype='number'
                style={{width: '40px'}}
            /> </td>
            <td> {item.need} </td>
            <td> {deficit} </td>
            <td> {formatEuro(item.unit_cost)} </td>
            <td> {formatEuro(cost)} </td>
        </tr>
    )
}

function clampPositive(x) {
    return x >= 0 ? x : 0
}

function formatEuro(x) {
    return x.toFixed(2) + ' €'
}

export default Ingredient
