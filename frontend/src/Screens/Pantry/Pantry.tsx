import React, { useState } from 'react'
import { ShoppingList, State } from '../../State/State.tsx';
import Backend from '../../Backend/Backend.tsx';
import RenderIngredient, { Numbers} from './Ingredient.tsx';

interface Props {
    backend: Backend;
    state: State;
    onComplete: () => void;
}

export default function Pantry(pp: Props) {
    const pt = new PantryTable(pp.state.shoppingList)
    const [total, setTotal] = useState(pt.recomputeTotal())
    return pt
        .attach(total, setTotal)
        .Render(pp.onComplete)
}

class PantryTable {
    shop: ShoppingList;
    total: number;
    setTotal: (x: number) => void;

    constructor(shop: ShoppingList) {
        this.shop = shop
    }

    attach(total: number, setTotal: (x: number) => void): PantryTable {
        this.total = total
        this.setTotal = setTotal
        return this
    }

    Render(onComplete: () => void): JSX.Element {
        return (
            <div id='pantry'>
            <table>
                <tbody>
                    <tr>
                        <th>Product</th>
                        <th>Have</th>
                        <th>Need</th>
                        <th>Deficit</th>
                        <th>Unit cost</th>
                        <th>Deficit cost</th>
                    </tr>
                    {
                        this.shop.ingredients.map(i => (
                            <RenderIngredient
                                ingredient={i}
                                onChange={() => this.setTotal(this.recomputeTotal())} />
                        ))
                    }
                    <tr>
                        <td colSpan={5}> <b>Total</b></td>
                        <td><b>{Numbers.asEuro(this.total)}</b></td>
                    </tr>
                </tbody>
            </table>
            </div>
        )
    }

    recomputeTotal(): number {
        return this.shop.ingredients.reduce((acc, i) => acc + i.price * Numbers.positive(i.need - i.have), 0)
    }
}

