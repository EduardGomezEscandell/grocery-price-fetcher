import React, { useState } from 'react'
import Backend from '../../Backend/Backend.tsx';
import TopBar from '../../TopBar/TopBar.tsx';
import { ShoppingNeeds, ShoppingList, State } from '../../State/State.tsx';
import SaveButton from '../../SaveButton/SaveButton.tsx';
import ShoppingItem from './ShoppingItem.tsx';

interface Props {
    backend: Backend;
    globalState: State;
    onBackToPantry: () => void;
}

enum Dialog {
    OFF,
    ON,
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
                <table>
                    <thead id='header1'>
                        <tr>
                            <th>
                                <button id='clear' onClick={() => {
                                    setDialog(Dialog.ON)
                                }}>x</button>
                            </th>
                            <th id='left'>Ingredient</th>
                            <th id='right'>Unitats</th>
                            <th id='right'>Paquets</th>
                        </tr>
                    </thead>
                    <tbody>
                        {
                            props.globalState.shoppingList.items.map((i, idx) =>
                                <ShoppingItem i={i} idx={idx} key={`${k}-${idx}`} globalState={props.globalState} />
                            )
                        }
                    </tbody>
                    <tfoot id='header2'>
                        <tr>
                            <td colSpan={4}>
                             Gràcies per utilitzar El Rebost!
                            </td>
                        </tr>
                    </tfoot>
                </table>
                {dialog !== Dialog.OFF &&
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


function saveShoppingList(backend: Backend, globalState: State): Promise<void> {
    return backend
        .Shopping()
        .POST({
            name: globalState.shoppingList.name, // Let the backend handle the name for now
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