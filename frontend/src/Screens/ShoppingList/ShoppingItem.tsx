import React, { useState } from 'react'
import { ShoppingListItem, State } from '../../State/State.tsx';
import { round2 } from '../../Numbers/Numbers.ts';

interface Props {
    globalState: State;
    i: ShoppingListItem;
    idx: number;
}

export default function ShoppingItem(props: Props) {
    const { i, idx } = props;

    const [selected, _setSelected] = useState(props.i.done)

    const flip = () => {
        const s = !selected
        _setSelected(s)
        props.globalState.shoppingList.items[idx].done = s
    }

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
            <td id='right'>{round2(i.units)}</td>
            <td id='right'>{round2(i.packs)}</td>
        </tr>
    )
}