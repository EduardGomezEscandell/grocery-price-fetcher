import React from 'react'
import Backend from '../../Backend/Backend.tsx';
import TopBar from '../../TopBar/TopBar.tsx';
import { ShoppingNeeds, ShoppingList, State } from '../../State/State.tsx';
import SaveButton from '../../SaveButton/SaveButton.tsx';
import ShoppingItem from './ShoppingItem.tsx';
import { asEuro } from '../../Numbers/Numbers.ts';

interface Props {
    backend: Backend;
    globalState: State;
    onBackToPantry: () => void;
}

export default function Shopping(props: Props): JSX.Element {
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
                            <th></th>
                            <th>Ingredient</th>
                            <th>Unitats</th>
                            <th>Paquets</th>
                            <th>Cost</th>
                        </tr>
                    </thead>
                    <tbody>
                        {
                            props.globalState.shoppingList.items.map((i, idx) =>
                                <ShoppingItem i={i} idx={idx} key={idx} globalState={props.globalState} />
                            )
                        }
                    </tbody>
                    <tfoot id='header2'>
                        <tr>
                            <td></td>
                            <td id='left' colSpan={3}>Total</td>
                            <td id='right' >{
                                asEuro(
                                    props.globalState.shoppingList.items.reduce((acc, i) => acc + i.cost, 0)
                                )
                            }</td>
                        </tr>
                    </tfoot>
                </table>
            </div>
        </>
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