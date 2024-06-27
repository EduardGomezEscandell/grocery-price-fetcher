import React, { useEffect, useState } from 'react'
import Select from 'react-select'
import { Dish } from '../../State/State'
import './DishPicker.css'

interface Props {
    dishes: Dish[];
    default: Dish;
    onChange: (d: Dish) => void;
    onRemove: () => void;
    key: string;
}

export default function DishPicker(pp: Props) {
    const options = pp.dishes.map(recipe => ({ value: recipe, label: recipe.name }))

    const [recipe, _setRecipe] = useState(pp.default)
    const [amount, _setAmount] = useState(pp.default.amount.toString())

    const setID = (d: Dish) => {
        _setRecipe(d)
        pp.onChange({
            id: d.id,
            name: d.name,
            amount: Number(amount),
        })
    }

    const setAmount = (a: string) => {
        _setAmount(a)
        pp.onChange({
            id: recipe.id,
            name: recipe.name,
            amount: Number(a)
        })
    }

    return (
        <div className='DishPicker' key={pp.key}>
            <input
                type="number"
                min="0"
                value={amount}
                onClick={(e) => { e.target instanceof HTMLInputElement && e.target.select() }}
                onChange={s => setAmount(s.target.value)}
                contentEditable={recipe.id !== 0}
            />
            <Select
                className='Select'
                styles={{
                    control: (base) => ({
                        ...base,
                        height: 40,
                        fontSize: 18,
                        padding: 0,
                        margin: 0,
                        borderRadius: 0,
                        border: "0px",
                        textAlign: "left",
                    }),
                    menu: (base) => ({
                        ...base,
                        fontSize: 18,
                        padding: 0,
                        margin: 0,
                        borderRadius: 0,
                        borderBottomRightRadius: 10,
                        borderBottomLeftRadius: 10,
                    }),
                    container: (base) => ({
                        ...base,
                        padding: 0,
                        margin: 0,
                    }),
                    input: (base) => ({
                        ...base,
                        padding: 0,
                        margin: 0,
                    }),
                }}
                components={{ DropdownIndicator:() => null, IndicatorSeparator:() => null }}
                onChange={selected => {
                    if (selected == null) {
                        return setID({ id: 0, name: '', amount: 0 })
                    }
                    setID(selected.value)
                }}
                value={{ value: recipe, label: recipe.name }}
                options={options}
                isSearchable
            />
            <button id='highlight'
                onClick={pp.onRemove}
            >X</button>
        </div>
    )
}

