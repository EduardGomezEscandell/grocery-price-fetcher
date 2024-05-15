import React, { useState } from 'react'
import { Ingredient, ShoppingList, State } from '../../State/State.tsx';
import Backend from '../../Backend/Backend.tsx';
import TopBar from '../../TopBar/TopBar.tsx';
import RenderIngredient, { Numbers } from './Ingredient.tsx';
import SaveButton from './SaveButton.tsx';

interface Props {
    backend: Backend;
    globalState: State;
    onBackToMenu: () => void;
}

export default function Pantry(pp: Props) {
    const total = new Total().compute(pp.globalState.shoppingList.ingredients)

    const [available, setAvailable] = useState(total.available)
    const [remaining, setRemaining] = useState(total.remaining)

    total
        .withAvailable(available, setAvailable)
        .withRemaining(remaining, setRemaining)

    const baseStyle: React.CSSProperties = {
        width: '800px',
    }

    return (
        <>
            <div>
                <TopBar
                    components={[
                        () => (
                            <td key='go-back' style={{ textAlign: 'left' }}>
                                <button
                                    className='TopBar.Button'
                                    onClick={pp.onBackToMenu}
                                >Tornar al menú</button>
                            </td>
                        ),
                        () => (<text key='pantry' className='TopBar.Text'>Rebost</text>),
                        () => (
                            <td key='save' style={{ textAlign: 'right' }}>
                            <SaveButton
                                className='TopBar__button'
                                backend={pp.backend}
                                globalState={pp.globalState}
                            />
                        </td>
                        )
                    ]}
                />
            </div>
            <div style={{ 
                margin: 'auto',
                alignItems: 'center',
                width: 'fit-content',
             }}>
                <PantryTable
                    shop={pp.globalState.shoppingList}
                    total={total}
                    style={baseStyle}
                />
            </div>
        </>
    )
}

class PantryTableProps {
    shop: ShoppingList
    total: Total
    style: React.CSSProperties
}

class PantryTable extends React.Component<PantryTableProps> {
    shop: ShoppingList
    total: Total
    style: React.CSSProperties

    constructor(pp: PantryTableProps) {
        super(pp)
        this.shop = pp.shop
        this.total = pp.total
        this.style = pp.style ? pp.style : {}
    }

    render(): JSX.Element {
        const blue0 = '#1a237e'
        const blue1 = '#405adb'
        const blue2 = '#6aa2f0'

        const headerStyle: React.CSSProperties = {
            textAlign: 'center',
            borderCollapse: 'collapse',
            backgroundColor: blue0,
            fontSize: '20px',
            fontWeight: 'bold',
            color: 'white',
        }

        const subHeaderStyle: React.CSSProperties = {
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

        const footerStyle: React.CSSProperties = {
            ...this.style,
            fontSize: '20px',
            background: 'orange',
        }

        return (
            <div key='pantry'>
                <table style={{
                    ...this.style,
                    textAlign: 'left',
                    borderCollapse: 'collapse',
                }}>
                    <tbody style={{ ...this.style, border: '1px solid black' }} key='pantry-body-1'>
                        <tr style={headerStyle} key='header-1'>
                            <td rowSpan={2}>Producte</td>
                            <td colSpan={3} style={{ background: blue1 }}> <b>Unitats</b></td>
                            <td colSpan={3} style={{ background: blue2 }}> <b>Paquets</b></td>
                            <th rowSpan={2}>Preu</th>
                        </tr>
                        <tr style={{
                            ...subHeaderStyle,
                            borderBottom: '1px solid black',
                        }} key='header-2'>
                            <td style={{ background: blue1 }}>Tens</td>
                            <td style={{ background: blue1 }}>Necessites</td>
                            <td style={{ background: blue1 }}>Manquen</td>

                            <td style={{ background: blue2 }}>Tamany</td>
                            <td style={{ background: blue2 }}>Manquen</td>
                            <td style={{ background: blue2 }}>Preu</td>
                        </tr>
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
                        <tr style={{
                            ...footerStyle,
                            borderTop: '1px solid black',
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
                        <tr style={footerStyle}>
                            <td colSpan={6} style={{ paddingLeft: '20px' }}>
                                Menjar que tens al rebost
                            </td>
                            <td style={{ textAlign: 'right' }}>
                                +
                            </td>
                            <td style={{
                                textAlign: 'right',
                                paddingRight: '20px'
                            }}>{Numbers.asEuro(this.total.available)}</td>
                        </tr>
                        <tr style={footerStyle}>
                            <td colSpan={6} style={{ paddingLeft: '20px' }}>
                                Menjar que quedarà al rebost
                            </td>
                            <td style={{ textAlign: 'right' }}>
                                -
                            </td>
                            <td style={{
                                textAlign: 'right',
                                paddingRight: '20px'
                            }}>{Numbers.asEuro(this.total.remaining)}</td>
                        </tr>
                        <tr style={{
                            ...footerStyle,
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
    setPurchased: (x: number) => void

    available: number;
    setAvailable: (x: number) => void

    consumed: number;
    setConsumed: (x: number) => void

    remaining: number;
    setRemaining: (x: number) => void

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
            .map(i => Numbers.positive(i.need) * i.price / i.batch_size)
            .reduce((acc, x) => acc + x, 0)
        this.available = i
            .map(i => Numbers.positive(i.have) * i.price / i.batch_size)
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
