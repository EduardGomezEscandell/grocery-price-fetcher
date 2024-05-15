import React, { useState } from 'react'
import { Ingredient, ShoppingList, State } from '../../State/State.tsx';
import Backend from '../../Backend/Backend.tsx';
import TopBar from '../../TopBar/TopBar.tsx';
import RenderIngredient, { Numbers } from './Ingredient.tsx';
import SaveButton from './SaveButton.tsx';
import './Pantry.css'

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
            <TopBar
                components={[
                    () => <button
                        onClick={pp.onBackToMenu}
                        key='go-back'>
                        Tornar al menú
                    </button>,
                    () => (<p key='pantry'>Rebost</p>),
                    () => (<SaveButton
                        key='save'
                        backend={pp.backend}
                        globalState={pp.globalState}
                    />)
                ]}
            />
            <PantryTable
                shop={pp.globalState.shoppingList}
                total={total}
                style={baseStyle}
            />
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
        return (
            <div key='pantry' className='Pantry'>
                <table className='Table'>
                    <thead>
                        <tr className='Header' key='header'>
                            <td rowSpan={2}>Producte</td>
                            <td colSpan={3} id='units'> <b>Unitats</b></td>
                            <td colSpan={3} id='packs'> <b>Paquets</b></td>
                            <th rowSpan={2}>Preu</th>
                        </tr>
                        <tr className='SubHeader' key='subheader' style={{
                            borderBottom: '1px solid black',
                        }}>
                            <td key='1' id='units'>Tens</td>
                            <td key='2' id='units'>Necessites</td>
                            <td key='3' id='units'>Manquen</td>

                            <td key='4' id='packs'>Tamany</td>
                            <td key='5' id='packs'>Manquen</td>
                            <td key='6' id='packs'>Preu</td>
                        </tr>
                    </thead>
                    <tbody>
                        {
                            this.shop.ingredients.map((i: Ingredient, idx: number) => (
                                <RenderIngredient 
                                    key={i.name}
                                    id={idx % 2 === 0 ? 'even' : 'odd'}
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
                    <tfoot>
                        <tr style={{
                            fontWeight: 'bold',
                        }}>
                            <td colSpan={7} style={{ paddingLeft: '20px' }}>Total a comprar</td>
                            <td className='Number'>{Numbers.asEuro(this.total.purchased)}</td>
                        </tr>
                        <tr >
                            <td colSpan={6} style={{ paddingLeft: '20px' }}>Menjar que tens al rebost</td>
                            <td className='Number'>+</td>
                            <td className='Number'>{Numbers.asEuro(this.total.available)}</td>
                        </tr>
                        <tr >
                            <td colSpan={6} style={{ paddingLeft: '20px' }}>Menjar que quedarà al rebost</td>
                            <td className='Number'>-</td>
                            <td className='Number'>{Numbers.asEuro(this.total.remaining)}</td>
                        </tr>
                        <tr>
                            <td colSpan={8} style={{background: 'black'}}>    </td>
                        </tr>
                        <tr style={{
                            fontWeight: 'bold',
                        }}>
                            <td colSpan={7} style={{ paddingLeft: '20px' }}>Cost del menjar consumit</td>
                            <td className='Number'>{Numbers.asEuro(this.total.consumed)}</td>
                        </tr>
                    </tfoot>
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
