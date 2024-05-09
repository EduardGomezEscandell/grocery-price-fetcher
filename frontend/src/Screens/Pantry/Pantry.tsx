import React, { useState } from 'react'
import { Ingredient, ShoppingList, State } from '../../State/State.tsx';
import Backend from '../../Backend/Backend.tsx';
import RenderIngredient, { Numbers } from './Ingredient.tsx';

interface Props {
    backend: Backend;
    state: State;
    onComplete: () => void;
}

export default function Pantry(pp: Props) {
    const total = new Total().compute(pp.state.shoppingList.ingredients)
    
    const [available, setAvailable] = useState(total.available)
    const [remaining, setRemaining] = useState(total.remaining)

    total
        .withAvailable(available, setAvailable)
        .withRemaining(remaining, setRemaining)
   
    return new PantryTable(pp.state.shoppingList, total).Render(pp.onComplete)
}

class PantryTable {
    shop: ShoppingList;
    total: Total;
    setTotal: (x: Total) => void;

    constructor(shop: ShoppingList, total: Total) {
        this.shop = shop
        this.total = total
    }

    attach(total: Total, setTotal: (x: Total) => void): PantryTable {
        this.total = total
        this.setTotal = setTotal
        return this
    }


    Render(onComplete: () => void): JSX.Element {
        const blue0 = '#1a237e'
        const blue1 = '#405adb'
        const blue2 = '#6aa2f0'

        const headerStyle: React.CSSProperties = {
            width: '800px',
            textAlign: 'center',
            borderCollapse: 'collapse',
            backgroundColor: blue0,
            fontSize: '20px',
            fontWeight: 'bold',
            color: 'white',
        }

        const subHeaderStyle: React.CSSProperties = {
            width: '800px',
            textAlign: 'center',
            borderCollapse: 'collapse',
            backgroundColor: blue0,
            fontSize: '20px',
            fontWeight: 'normal',
            color: 'white',
        }

        const rowStyle = (idx: number): React.CSSProperties => {
            return {
                background: idx % 2 === 0 ? 'white' : '#eeeeee',
                fontSize: '16px',
            }
        }

        return (
            <div key='pantry'>
                <table style={{
                    width: '800px',
                    textAlign: 'left',
                    borderCollapse: 'collapse',
                }}>
                    <tbody style={{ border: '1px solid black' }} key='pantry-body-1'>
                        <tr style={headerStyle} key='header-1'>
                            <td rowSpan={2}>Producte</td>
                            <td colSpan={3} style={{ background: blue1 }}> <b>Unitats</b></td>
                            <td colSpan={3} style={{ background: blue2 }}> <b>Paquets</b></td>
                            <th rowSpan={2}>Preu</th>
                        </tr>
                        <tr style={subHeaderStyle} key='header-2'>
                            <td style={{ background: blue1 }}>Tens</td>
                            <td style={{ background: blue1 }}>Necessites</td>
                            <td style={{ background: blue1 }}>Manquen</td>

                            <td style={{ background: blue2 }}>Tamany</td>
                            <td style={{ background: blue2 }}>Manquen</td>
                            <td style={{ background: blue2 }}>Preu</td>
                        </tr>
                    </tbody>
                    <tbody style={{ border: '1px solid black' }}>
                        {
                            this.shop.ingredients.map((i: Ingredient, idx: number) => (
                                <RenderIngredient
                                    style={rowStyle(idx)}
                                    ingredient={i}
                                    onChange={(value: number) => {
                                        i.have = value
                                        this.total
                                            .compute(this.shop.ingredients)
                                            .commit()
                                    }}
                                />
                            ))
                        }
                    </tbody>
                    <tbody style={{
                        border: '1px solid black',
                        background: 'orange',
                    }}>
                        <tr style={{
                            fontSize: '20px',
                            fontWeight: 'bold',
                        }}>
                            <td colSpan={7} style={{ paddingLeft: '20px' }}>
                                Total a comprar
                            </td>
                            <td style={{
                                textAlign: 'right',
                                paddingRight: '20px'
                            }}>{Numbers.asEuro(this.total.purchased)}</td>
                        </tr>
                        <tr style={{ fontSize: '20px', }}>
                            <td colSpan={6} style={{ paddingLeft: '20px' }}>
                                Menjar que tens rebost
                            </td>
                            <td style={{ textAlign: 'right'}}>
                                +
                            </td>
                            <td style={{
                                textAlign: 'right',
                                paddingRight: '20px'
                            }}>{Numbers.asEuro(this.total.available)}</td>
                        </tr>
                        <tr style={{ fontSize: '20px', }}>
                            <td colSpan={6} style={{ paddingLeft: '20px' }}>
                                Menjar que quedarà al rebost
                            </td>
                            <td style={{ textAlign: 'right'}}>
                                -
                            </td>
                            <td style={{
                                textAlign: 'right',
                                paddingRight: '20px'
                            }}>{Numbers.asEuro(this.total.remaining)}</td>
                        </tr>
                        <tr style={{
                            fontSize: '20px',
                            fontWeight: 'bold',
                            borderTop: '1px solid black',
                        }}>
                            <td colSpan={7} style={{ paddingLeft: '20px' }}>
                                Cost del menjar consumit
                            </td>
                            <td style={{
                                textAlign: 'right',
                                paddingRight: '20px',
                                width: '300px',
                            }}>{Numbers.asEuro(this.total.consumed)}</td>
                        </tr>
                    </tbody>
                </table>
            </div>
        )
    }


}


class Total {
    purchased: number;
    setPurchased: (x :number) => void

    available: number;
    setAvailable: (x :number) => void

    consumed: number;
    setConsumed: (x :number) => void

    remaining: number;
    setRemaining: (x :number) => void

    withAvailable(a: number, update: (x: number) => void): Total {
        this.available = a
        this.setAvailable = update
        return this
    }

    withRemaining(r: number, update: (x: number) => void): Total {
        this.remaining = r
        this.setRemaining = update
        return this
    }

    compute(i: Ingredient[]): Total {
        this.consumed = i
            .map(i => Numbers.positive(i.need - i.have) * i.price / i.batch_size )
            .reduce((acc, x) => acc + x, 0)
        this.available = i
            .map(i => Numbers.positive(i.have) * i.price)
            .reduce((acc, x) => acc + x, 0)
        this.purchased = i
            .map(i => Math.ceil(Numbers.positive(i.need - i.have) / i.batch_size) * i.price)
            .reduce((acc, x) => acc + x, 0)
        this.remaining = this.purchased + this.available - this.consumed

        return this
    }

    commit(): Total {
        this.setAvailable(this.available)
        this.setRemaining(this.remaining)
        return this
    }
}
