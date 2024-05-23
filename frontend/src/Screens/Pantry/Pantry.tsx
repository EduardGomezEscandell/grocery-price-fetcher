import React, { useState } from 'react'
import { Ingredient, ShoppingList, ShoppingNeeds, State } from '../../State/State.tsx';
import Backend from '../../Backend/Backend.tsx';
import TopBar from '../../TopBar/TopBar.tsx';
import SaveButton from '../../SaveButton/SaveButton.tsx';
import { FocusIngredient, RowIngredient } from './PantryIngredient.tsx';
import { asEuro, positive } from '../../Numbers/Numbers.ts'

interface Props {
    backend: Backend;
    globalState: State;
    onBackToMenu: () => void;
    onComplete: () => void;
}

export default function Pantry(pp: Props) {
    const total = new Total().compute(pp.globalState.inNeed.ingredients)

    const [available, setAvailable] = useState(total.available)
    const [remaining, setRemaining] = useState(total.remaining)

    total
        .withAvailable(available, setAvailable)
        .withRemaining(remaining, setRemaining)

    return (
        <>
            <TopBar
                left={<button onClick={pp.onBackToMenu} key='go-back'>Tornar al menú</button>}
                right={<SaveButton
                    key='save'
                    baseTxt='Desar'

                    onSave={() => savePantry(pp.backend, pp.globalState)}
                    onSaveTxt='Desant...'

                    onAccept={() => { pp.onComplete() }}
                    onAcceptTxt='Desat'

                    onReject={(reason: any) => console.log('Error saving pantry: ', reason || 'Unknown error')}
                    onRejectTxt='Error'
                />}
            />
            <PantryTable
                inNeed={pp.globalState.inNeed}
                total={total}
            />
        </>
    )
}

export function savePantry(backend: Backend, globalState: State): Promise<void> {
    return Promise.all([
        backend
            .Pantry()
            .POST({
                name: '', // Let the backend handle the name for now
                contents: globalState.inNeed.ingredients
                    .filter(i => i.have > 0)
                    .map(i => {
                        return { name: i.name, amount: i.have }
                    })
            }).then(() => { }),
        backend
            .Shopping()
            .GET()
            .then((s: ShoppingList[]) => s.length > 0 ? s[0] : new ShoppingList())
            .then(s => globalState.shoppingList = makeShoppingList(s, globalState.inNeed))
            .then(() => { })
    ]).then(() => { })
}

function makeShoppingList(list: ShoppingList, needs: ShoppingNeeds): ShoppingList {
    return {
        name: list.name,
        timeStamp: list.timeStamp,
        items: needs.ingredients
            .filter(i => i.need > 0)
            .map(i => {
                const alreadyBought = list.items.findIndex(si => si.name === i.name) !== -1
                const mustBuy = Math.max(0, i.need - i.have)

                return {
                    name: i.name,
                    done: alreadyBought,
                    units: mustBuy,
                    packs: Math.ceil(mustBuy / i.batch_size),
                    cost: Math.ceil(mustBuy / i.batch_size) * i.price,
                }
            })
    }
}

class PantryTableProps {
    inNeed: ShoppingNeeds
    total: Total
}

function PantryTable(pp: PantryTableProps): JSX.Element {
    const [focussed, setFocussed] = useState<RowIngredient | undefined>(undefined)

    const tableStyle: React.CSSProperties = {}
    if (focussed) {
        tableStyle.filter = 'blur(5px)'
    }

    return (
        <div className='scroll-table' key='pantry'>
            <table style={tableStyle}>
                <thead>
                    <tr key='header' id='header1'>
                        <th id="left">Producte</th>
                        <th id="center">Tens</th>
                        <th id="right">A comprar</th>
                    </tr>
                </thead>
                <tbody>
                    {
                        pp.inNeed.ingredients.map((i: Ingredient, idx: number) => (
                            <RowIngredient
                                key={i.name}
                                id={idx % 2 === 0 ? 'even' : 'odd'}
                                ingredient={i}
                                onChange={(value: number) => {
                                    i.have = value
                                    pp.total
                                        .compute(pp.inNeed.ingredients)
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
                    <tr><td colSpan={3} id='header1' /></tr>
                    <tr>
                        <td colSpan={2} id='left'>Total a comprar</td>
                        <td id='right'>{asEuro(pp.total.purchased)}</td>
                    </tr>
                    <tr>
                        <td colSpan={2} id='left'>Cost del menjar consumit</td>
                        <td id='right'>{asEuro(pp.total.consumed)}</td>
                    </tr>
                </tfoot>
            </table>
            {
                focussed && <FocusIngredient
                    ingredient={focussed.props.ingredient}
                    onClose={() => setFocussed(undefined)}
                    onChange={(value: number) => {
                        focussed.props.ingredient.have = value
                        pp.total
                            .compute(pp.inNeed.ingredients)
                            .commit()
                    }}
                />
            }
        </div>
    )
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
