import React from 'react'
import Backend from '../../Backend/Backend.tsx';
import TopBar from '../../TopBar/TopBar.tsx';
import { State } from '../../State/State.tsx';
import SaveButton from '../../SaveButton/SaveButton.tsx';

interface Props {
    backend: Backend;
    globalState: State;
    onBackToPantry: () => void;
}

export default function ShoppingList(props: Props): JSX.Element {
    return (
        <>
            <TopBar
                left={<button onClick={props.onBackToPantry} key='go-back'>Tornar al rebost</button>}
                right={<SaveButton
                    key='save'

                    baseTxt='Desar'
                    onSave={() => savePantry(props.backend, props.globalState)}
                    onSaveTxt='Desant...'
                    onAcceptTxt='Desat'
                    onRejectTxt='Error'
                />}
            />
            <div className='shopping-list'>
                <h1 id='header1'>Llista de la compra</h1>
            </div>
        </>
    )
}


export function savePantry(backend: Backend, globalState: State): Promise<void> {
    return backend
        .Pantry()
        .POST({
            name: '', // Let the backend handle the name for now
            contents: globalState.shoppingList.ingredients
                .filter(i => i.have > 0)
                .map(i => {
                    return { name: i.name, amount: i.have }
                })
        })
}