import React, { useState } from 'react'
import { ShoppingListItem, State } from '../../State/State.tsx';
import { asEuro, round2 } from '../../Numbers/Numbers.ts';

interface Props {
    globalState: State;
    i: ShoppingListItem;
    idx: number;
    show: Column;
}

export enum Column {
    UNITS = 0, PACKS, COST
}

export default function ShoppingItem(props: Props) {
    const { i, idx } = props;

    const [selected, _setSelected] = useState(props.i.done)

    const flip = () => {
        const s = !selected
        _setSelected(s)
        props.globalState.shoppingList.items[idx].done = s
    }

    const show = ((): string => {
        switch(props.show) {
            case Column.UNITS: return round2(i.units)
            case Column.PACKS: return round2(i.packs)
            case Column.COST: return asEuro(i.cost * i.packs)
        }
    })()

    return (
        <tr
            key={idx}
            id={i.units === 0 || selected ? 'lowlight' : 'even'}
            onClick={flip}
        >
            <td>
                <input
                    type='checkbox'
                    checked={selected}
                    onChange={flip}
                />
            </td>
            <td id='left'>{i.name}</td>
            <td id='right'>{show}</td>
        </tr>
    )
}