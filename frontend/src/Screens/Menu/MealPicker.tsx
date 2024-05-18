import React, { useEffect, useState } from 'react'
import Select from 'react-select'
import { Dish } from '../../State/State'

interface Props {
    recipes: string[];
    default: Dish;
    onChange: (d: Dish) => void;
}

function MealPicker(pp: Props) {
    const options = pp.recipes.map(recipe => ({ value: recipe, label: recipe }))

    const [name, _setName] = useState(pp.default.name)
    const [amount, _setAmount] = useState(pp.default.amount.toString())

    useEffect(() => {
        _setName(pp.default.name)
        _setAmount(pp.default.amount.toString())
    }, [pp.default])

    const setName = (n: string) => {
        _setName(n)
        pp.onChange(pp.default.withName(n))
    }

    const setAmount = (a: string) => {
        _setAmount(a)
        pp.onChange(pp.default.withAmount(Number(a)))
    }

    return (
        <table>
            <tbody>
                <tr>
                    <td>
                        <Select
                            styles={{ control: (base) => ({ ...base, width: '200px' }) }}
                            onChange={selected => {
                                if (selected == null) {
                                    return setName("")
                                }
                                setName(selected.value)
                            }}
                            value={{ value: name, label: name }}
                            options={options}
                            isClearable
                        />
                    </td>
                    <td>
                        <input
                            type="number"
                            min="0"
                            value={amount}
                            onClick={(e) => {e.target instanceof HTMLInputElement && e.target.select()}}
                            onChange={s => setAmount(s.target.value)}
                            contentEditable={name !== undefined && name !== ""}
                            style={{ width: '50px', height: '25px' }}
                        />
                    </td>
                </tr>
            </tbody>
        </table>
    )
}

export default MealPicker
