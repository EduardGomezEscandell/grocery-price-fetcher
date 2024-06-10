import React, { useState } from 'react'
import { ShoppingListItem } from '../../State/State.tsx';
import { asEuro, round2 } from '../../Numbers/Numbers.ts';

interface Props {
    item: ShoppingListItem;
    idx: number;
    show: Column;
    setSelection: (v: boolean) => void;
}

export enum Column {
    UNITS = 0, PACKS, COST
}

export default function ShoppingItem(props: Props) {
    const { item, idx } = props;

    const [selected, _setSelected] = useState(props.item.done)

    const flip = () => {
        const s = !selected
        _setSelected(s)
        props.setSelection(s)
    }

    const show = ((): string => {
        switch(props.show) {
            case Column.UNITS: return round2(item.units)
            case Column.PACKS: return round2(item.packs)
            case Column.COST: return asEuro(item.cost * item.packs)
        }
    })()

    return (
        <tr
            key={idx}
            id={item.units === 0 || selected ? 'lowlight' : 'even'}
            onClick={flip}
        >
            <td>
                <input
                    type='checkbox'
                    checked={selected}
                    onChange={flip}
                />
            </td>
            <td id='left'>{item.name}</td>
            <td id='right'>{show}</td>
        </tr>
    )
}