import React, { useState } from 'react'
import Backend from '../../Backend/Backend.tsx';
import TopBar from '../../TopBar/TopBar.tsx';
import { ShoppingNeeds, ShoppingList, State } from '../../State/State.tsx';
import SaveButton from '../../SaveButton/SaveButton.tsx';
import ShoppingItem, { Column } from './ShoppingItem.tsx';
import { asEuro } from '../../Numbers/Numbers.ts';
import './ShoppingList.css';

interface Props {
    backend: Backend;
    globalState: State;
    onBackToPantry: () => void;
    onGotoHome: () => void;
}

enum Dialog {
    OFF,
    RESTORE,
    HELP,
    CLOSING,
}

export default function Shopping(props: Props): JSX.Element {
    const [dialog, _setDialog] = useState(Dialog.OFF);
    const [k, setK] = useState(0)
    const forceChildUpdate = () => setK(k + 1)

    const setDialog = (d: Dialog): boolean => {
        if (dialog === Dialog.CLOSING) return false
        _setDialog(d)
        return true
    }

    const [column, setColumn] = useState(Column.UNITS)

    const tableStyle: React.CSSProperties = {}
    if (dialog !== Dialog.OFF) {
        tableStyle.filter = 'blur(5px)'
    }

    return (
        <>
            <TopBar
                left={<SaveButton
                    key='save'

                    baseTxt='Tornar'
                    onSave={() => saveShoppingList(props.backend, props.globalState)}
                    onSaveTxt='Desant...'

                    onAccept={() => props.onBackToPantry()}
                    onAcceptTxt='Desat'

                    onRejectTxt='Error'
                />}
                logoOnClick={() => { saveShoppingList(props.backend, props.globalState).then(props.onGotoHome) }}
                titleOnClick={() => setDialog(Dialog.HELP)}
                titleText="La&nbsp;meva compra"
                right={<SaveButton
                    key='save'

                    baseTxt='Desar'
                    onSave={() => saveShoppingList(props.backend, props.globalState)}
                    onSaveTxt='Desant...'
                    onAcceptTxt='Desat'
                    onRejectTxt='Error'
                />}
            />
            <div className='scroll-table'>
                <table style={tableStyle}>
                    <thead id='header1'>
                        <tr>
                            <th>
                                <button id='clear' onClick={() => {
                                    setDialog(Dialog.RESTORE)
                                }}>x</button>
                            </th>
                            <th id='left'>Ingredient</th>
                            <th id='right'>
                                <select
                                    value={column}
                                    onChange={e => setColumn(e.target.selectedIndex as Column)}
                                >
                                    <option value={Column.UNITS}>Unitats</option>
                                    <option value={Column.PACKS}>Paquets</option>
                                    <option value={Column.COST}>Cost</option>
                                </select>
                            </th>
                        </tr>
                    </thead>
                    <tbody>
                        {
                            props.globalState.shoppingList.items.map((i, idx) =>
                                <ShoppingItem i={i} idx={idx} key={`${k}-${idx}-${column}`} globalState={props.globalState} show={column}/>
                            )
                        }
                    </tbody>
                    <tfoot id='header2'>
                        <tr>
                            <td></td>
                            <td id='left'>Cost total</td>
                            <td id='right'>
                                {
                                    (() => {
                                        const cost = props.globalState.shoppingList.items.reduce((acc, i) => acc + i.packs * i.cost, 0)
                                        return (<>{asEuro(cost)}</>)
                                    })()
                                }
                            </td>
                        </tr>
                    </tfoot>
                </table>
                {(dialog === Dialog.RESTORE || dialog === Dialog.CLOSING) &&
                    <ResetDialog
                        onReset={() => {
                            if (!setDialog(Dialog.CLOSING)) {
                                return
                            }
                            props.globalState.shoppingList.items.forEach(i => i.done = false)
                            saveShoppingList(props.backend, props.globalState)
                                .then(() => forceChildUpdate())
                                .then(() => setDialog(Dialog.OFF))
                        }}
                        onExit={() => setDialog(Dialog.OFF)}
                    />
                }{
                    dialog === Dialog.HELP &&
                    <HelpDialog
                        onClose={() => setDialog(Dialog.OFF)}
                    />
                }
            </div>
        </>
    )
}

function ResetDialog(props: {
    onReset: () => void
    onExit: () => void
}): JSX.Element {
    return (
        <dialog open>
            <h3 id='header'>
                Restaurar la llista de la compra?
            </h3>
            <div id='body'>
                <p>Tots els elements marcats com a comprats es desmarcaran</p>
                <p>Aquesta acció és irreversible, prem Tornar si no vols realitzar-la</p>
            </div>
            <div id='footer'>
                <button id='dialog-left' onClick={props.onExit}>Tornar</button>
                <button id='dialog-right' onClick={props.onReset}>Restaurar</button>
            </div>
        </dialog>
    )
}

function HelpDialog(props: {
    onClose: () => void
}): JSX.Element {
    return (
        <dialog open>
            <h2 id="header">La meva compra</h2>
            <div id="body">
                <p>
                    Aquesta pàgina mostra una llista dels ingredients que necessites comprar per al
                    teu menú setmanal, descomptant-li el que ja tens al teu rebost. Pots fer clic a qualsevol
                    ingredient per marcar-lo com a comprat.
                </p>
                <p>
                    Per a cada ingredient, pots escollir quina informació vols veure. Expliquem-ho amb un exemple: si
                    necessites 9 ous que es venen a 2€ cada mitja dotzena:
                </p>
                <p>
                    <b>Unitats:</b> Nombre d'unitats que necessites comprar. En aquest cas, 9 ous.
                </p>
                <p>
                    <b>Paquets:</b> Nombre de paquets que necessites comprar. En aquest cas, dues mitges dotzenes.
                </p>
                <p>
                    <b>Cost:</b> Cost total de la compra. En aquest cas, 4€.
                </p>
            </div>
            <div id="footer">
                <button onClick={props.onClose}>Tancar</button>
            </div>
        </dialog>
    )
}


function saveShoppingList(backend: Backend, globalState: State): Promise<void> {
    return backend
        .Shopping()
        .POST({
            name: globalState.shoppingList.name,
            items: globalState.shoppingList.items
                .filter(i => i.done)
                .map(i => i.name)
        })
}

export async function downloadShoppingList(backend: Backend, globalState: State): Promise<void> {
    const lists = await backend.Shopping().GET();
    const s = lists[0] || new ShoppingList();
    globalState.shoppingList = makeShoppingList(s, globalState.inNeed);
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