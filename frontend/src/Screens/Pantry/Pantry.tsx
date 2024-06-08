import React, { useState } from 'react'
import { Ingredient, State } from '../../State/State.tsx';
import Backend from '../../Backend/Backend.tsx';
import TopBar from '../../TopBar/TopBar.tsx';
import SaveButton from '../../SaveButton/SaveButton.tsx';
import { FocusIngredient, RowIngredient } from './PantryIngredient.tsx';
import { asEuro } from '../../Numbers/Numbers.ts'
import { downloadShoppingList } from '../ShoppingList/ShoppingList.tsx';
import { IngredientUsage } from '../../Backend/endpoints/IngredientUse.tsx';

interface Props {
    backend: Backend;
    globalState: State;
    onBackToMenu: () => void;
    onGotoHome: () => void;
    onComplete: () => void;
}

export default function Pantry(pp: Props) {
    const [savings, setSavings] = useState(computeSavings(pp.globalState.inNeed.ingredients))
    const [help, setHelp] = useState(false)

    const [focussed, setFocussed] = useState<{
        ingredient: Ingredient
        usage: IngredientUsage[]
     } | undefined>(undefined)


     const tableStyle: React.CSSProperties = {}
    if (focussed || help) {
        tableStyle.filter = 'blur(5px)'
    }

    const ingr = pp.globalState.inNeed.ingredients
    const updateSavings = () => {
        setSavings(computeSavings(ingr))
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
                                        pp.globalState.inNeed.ingredients[idx].have = value
                                        updateSavings()
                                    }}
                                    onClick={() => {
                                        if (focussed) {
                                            setFocussed(undefined)
                                        } else {
                                            pp.backend.IngredientUse().POST({
                                                MenuName: pp.globalState.menu.name,
                                                IngredientName: i.name
                                            }).then((usage) => {
                                                setFocussed({
                                                    ingredient: i,
                                                    usage: usage
                                                })
                                            }).catch((reason) => {
                                                console.log('Error getting ingredient usage: ', reason || 'Unknown error')
                                            })
                                        }
                                    }}
                                />
                            ))
                        }
                    </tbody>
                    <tfoot id='header2'>
                        <tr><td colSpan={2} id='header1' /></tr>
                        <tr>
                            <td id='left'>T'estalvies</td>
                            <td id='right'>{asEuro(savings)}</td>
                        </tr>
                    </tfoot>
                </table>
                {
                    focussed && <FocusIngredient
                        ingredient={focussed.ingredient}
                        usage={focussed.usage}
                        onClose={() => setFocussed(undefined)}
                        onChange={(value: number) => {
                            focussed.ingredient.have = value
                            updateSavings()
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
                            <p>
                                Si fas clic en un ingredient, veuràs quins dies, àpats i receptes l'utilitzen en el teu menu.
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

function computeSavings(ingredients: Ingredient[]): number {
    return ingredients
        .map(i => Math.min(i.have, i.need) * i.price / i.batch_size)
        .reduce((acc, x) => acc + x, 0)
}