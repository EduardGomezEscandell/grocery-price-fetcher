import React, { useEffect, useState } from 'react'
import Select from 'react-select'
import { Dish } from '../../State/State'
import './DishPicker.css'

interface Props {
    recipes: string[];
    default: Dish;
    onChange: (d: Dish) => void;
    onRemove: () => void;
    key: string;
}

export default function DishPicker(pp: Props) {
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
        <div className='DishPicker' key={pp.key}>
            <input
                type="number"
                min="0"
                value={amount}
                onClick={(e) => { e.target instanceof HTMLInputElement && e.target.select() }}
                onChange={s => setAmount(s.target.value)}
                contentEditable={name !== undefined && name !== ""}
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
                        return setName("")
                    }
                    setName(selected.value)
                }}
                value={{ value: name, label: name }}
                options={options}
                isSearchable
            />
            <button id='highlight'
                onClick={pp.onRemove}
            >X</button>
        </div>
    )
}

