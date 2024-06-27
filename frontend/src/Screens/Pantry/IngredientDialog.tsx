import React from 'react';
import { Ingredient } from '../../State/State';
import { IngredientUsage } from '../../Backend/endpoints/IngredientUse';


interface Props {
    item: Ingredient;
    usage: IngredientUsage[];
    onClose: () => void;
}

export default function IngredientDialog(props: Props): JSX.Element {
    return (
        <dialog open id='pantry-ingredient'>
            <h2 id="header">{props.item.name}</h2>
            <div id="body">
                <p>
                    L'ingredient <b>{props.item.name}</b> apareix en els següents plats:
                </p>
                <div className='vert-scroll'>
                    <div className='scroll-table'>
                        <table>
                            <tbody>
                                {props.usage.map((u, idx) =>
                                    <tr key={u.day + u.meal + u.dish} id={idx % 2 === 0 ? 'even' : 'odd'}>
                                        <td id="left">{u.meal} de {u.day}</td>
                                        <td id="left">{u.dish}</td>
                                        <td id="right">{u.amount}</td>
                                    </tr>
                                )}
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>
            <div id="footer">
                <button onClick={props.onClose}>OK</button>
            </div>
        </dialog>
    )
}