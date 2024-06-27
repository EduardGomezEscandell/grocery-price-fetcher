import React, { useState } from 'react'
import { Ingredient } from '../../State/State'
import './Pantry.css'

interface Props {
    item: Ingredient;
    id: string;
    onChange: (newHave: number) => void;
    onClick: () => void;
}

export default function IngredientRow(props: Props): JSX.Element {
    const [highlight, setHighlight] = useState(false)
    const [have, setHave] = useState(
        (Math.round(props.item.amount*1e4)/1e4) // Clean up floating point errors
        .toString()
    )

    return (
        <tr key={props.item.name}
            id={highlight ? 'highlight' : props.id}
            onMouseEnter={() => setHighlight(true)}
            onMouseLeave={() => setHighlight(false)}
            onClick={(e) => {
                e.target instanceof HTMLInputElement || props.onClick()
            }}
        >
            <td id='left' key='name'> {props.item.name}  </td>
            <td id='right' key='have'>
                <input
                    type="number"
                    value={have}
                    onClick={(e) => { e.target instanceof HTMLInputElement && e.target.select() }}
                    onChange={(s) => {
                        if (s.target.value === '') {
                            setHave('0')
                            props.onChange(0)
                            return
                        } 
                        else if (s.target.valueAsNumber < 0) {
                            setHave('0')
                            props.onChange(0)
                        } else {
                            setHave(s.target.value)
                            props.onChange(s.target.valueAsNumber)
                        }
                    }}
                    datatype='number'
                />
            </td>
        </tr>
    )
}