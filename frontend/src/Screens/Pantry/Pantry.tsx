import React, { useState } from 'react'
import { Ingredient, State } from '../../State/State.tsx';
import Backend from '../../Backend/Backend.tsx';
import TopBar from '../../TopBar/TopBar.tsx';
import SaveButton from '../../SaveButton/SaveButton.tsx';
import { FocusIngredient, RowIngredient } from './PantryIngredient.tsx';
import { asEuro, positive } from '../../Numbers/Numbers.ts'
import { downloadShoppingList } from '../ShoppingList/ShoppingList.tsx';

interface Props {
    backend: Backend;
    globalState: State;
    onBackToMenu: () => void;
    onGotoHome: () => void;
    onComplete: () => void;
}

export default function Pantry(pp: Props) {
    const total = new Total().compute(pp.globalState.inNeed.ingredients)

    const [available, setAvailable] = useState(total.available)
    const [remaining, setRemaining] = useState(total.remaining)
    total
        .withAvailable(available, setAvailable)
        .withRemaining(remaining, setRemaining)

    const [help, setHelp] = useState(false)
    const [focussed, setFocussed] = useState<RowIngredient | undefined>(undefined)
    const tableStyle: React.CSSProperties = {}
    if (focussed || help) {
        tableStyle.filter = 'blur(5px)'
    }

    return (
        <>
            <TopBar
                left={<SaveButton
                    key='save'

                    baseTxt='Tornar'
                    onSave={() => savePantry(pp.backend, pp.globalState)}
                    onSaveTxt='Desant...'

                    onAccept={() => pp.onBackToMenu()}
                    onAcceptTxt='Desat'

                    onRejectTxt='Error'
                />}
                logoOnClick={() => { savePantry(pp.backend, pp.globalState).then(pp.onGotoHome) }}
                titleOnClick={() => setHelp(true)}
                titleText='El&nbsp;meu rebost'
                right={<SaveButton
                    key='save'
                    baseTxt='Següent'

                    onSave={() => Promise.all([
                        savePantry(pp.backend, pp.globalState),
                        downloadShoppingList(pp.backend, pp.globalState),
                    ])}
                    onSaveTxt='Desant...'

                    onAccept={() => { pp.onComplete() }}
                    onAcceptTxt='Desat'

                    onReject={(reason: any) => console.log('Error saving pantry: ', reason || 'Unknown error')}
                    onRejectTxt='Error'
                />}
            />
            <div className='scroll-table' key='pantry'>
                <table style={tableStyle}>
                    <thead>
                        <tr key='header' id='header1'>
                            <th id="left">Producte</th>
                            <th id="right">Tens</th>
                        </tr>
                    </thead>
                    <tbody>
                        {
                            pp.globalState.inNeed.ingredients.map((i: Ingredient, idx: number) => (
                                <RowIngredient
                                    key={i.name}
                                    id={idx % 2 === 0 ? 'even' : 'odd'}
                                    ingredient={i}
                                    onChange={(value: number) => {
                                        i.have = value
                                        total
                                            .compute(pp.globalState.inNeed.ingredients)
                                            .commit()
                                    }}
                                    onClick={(ri: RowIngredient) => {
                                        if (focussed) {
                                            setFocussed(undefined)
                                        } else {
                                            setFocussed(ri)
                                        }
                                    }}
                                />
                            ))
                        }
                    </tbody>
                    <tfoot id='header2'>
                        <tr><td colSpan={2} id='header1' /></tr>
                        <tr>
                            <td id='left'>Total a comprar</td>
                            <td id='right'>{asEuro(total.purchased)}</td>
                        </tr>
                        <tr>
                            <td id='left'>Cost del menjar consumit</td>
                            <td id='right'>{asEuro(total.consumed)}</td>
                        </tr>
                    </tfoot>
                </table>
                {
                    focussed && <FocusIngredient
                        ingredient={focussed.props.ingredient}
                        onClose={() => setFocussed(undefined)}
                        onChange={(value: number) => {
                            focussed.props.ingredient.have = value
                            total
                                .compute(pp.globalState.inNeed.ingredients)
                                .commit()
                        }}
                    />
                }
                {
                    help && <dialog open>
                        <h2 id="header">El meu rebost</h2>
                        <div id="body">
                            <p>
                                Aquesta pàgina mostra una llista dels ingredients que necessites per al teu menú setmanal.
                            </p>
                            <p>
                                Per a cada ingredient, indica quant en tens al teu rebost i 
                                així <i>La compra de l'Edu</i> podrà calcular quant en necessites comprar.
                            </p>
                        </div>
                        <div id="footer">
                            <button onClick={() => setHelp(false)}>
                                D'acord
                            </button>
                        </div>
                    </dialog>
                }
            </div>

        </>
    )
}

export async function savePantry(backend: Backend, globalState: State): Promise<void> {
    await backend
        .Pantry()
        .POST({
            name: 'default',
            contents: globalState.inNeed.ingredients
                .filter(i => i.have > 0)
                .map(i_1 => {
                    return { name: i_1.name, amount: i_1.have };
                })
        });
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
            .map(i => positive(i.need) * i.price / i.batch_size)
            .reduce((acc, x) => acc + x, 0)
        this.available = i
            .map(i => positive(i.have) * i.price / i.batch_size)
            .reduce((acc, x) => acc + x, 0)
        this.purchased = i
            .map(i => Math.ceil(positive(i.need - i.have) / i.batch_size) * i.price)
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
