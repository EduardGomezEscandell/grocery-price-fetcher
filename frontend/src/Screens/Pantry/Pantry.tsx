import React, { useEffect, useState } from 'react'
import { Pantry, PantryItem, ShoppingNeeds, ShoppingNeedsItem } from '../../State/State.tsx';
import Backend from '../../Backend/Backend.tsx';
import TopBar from '../../TopBar/TopBar.tsx';
import SaveButton from '../../SaveButton/SaveButton.tsx';
import IngredientRow from './PantryIngredient.tsx';
import IngredientDialog from './IngredientDialog.tsx';
import { asEuro } from '../../Numbers/Numbers.ts'
import { IngredientUsage } from '../../Backend/endpoints/IngredientUse.tsx';

interface Props {
    backend: Backend;
    sessionName: string;
    onBackToMenu: () => void;
    onGotoHome: () => void;
    onComplete: () => void;
}

interface Focus {
    item: PantryItem
    usage: IngredientUsage[]
}

export default function RenderPantry(pp: Props) {
    const [needs, setNeeds] = useState<ShoppingNeeds>(new ShoppingNeeds())
    const [pantry, setPantry] = useState<Pantry>(new Pantry())
    const [help, setHelp] = useState(false)
    const [focussed, setFocussed] = useState<Focus | undefined>(undefined)

    const tableStyle: React.CSSProperties = {}
    if (focussed || help) {
        tableStyle.filter = 'blur(5px)'
    }

    const computeSavings = (): number => {
        return merge(
            needs.contents,
            pantry.contents,
            (need, have) => need.name.localeCompare(have.name),
            (need, have) => ((need && need.price) || 0) * Math.max(0,
                ((have && have.amount) || 0) - ((need && need.amount) || 0))
        )
    }

    useEffect(() => {
        Promise.all([
            pp.backend.Needs(pp.sessionName).GET(),
            pp.backend.Pantry(pp.sessionName).GET(),
        ])
            .then(([needs, pantry]) => {
                setNeeds(needs)
                setPantry(pantry)
            })
            .catch((reason) => {
                console.log('Error getting pantry: ', reason || 'Unknown error')
            })
    }, [pp.backend, pp.sessionName])

    return (
        <>
            <TopBar
                left={<SaveButton
                    key='save'

                    baseTxt='Tornar'
                    onSave={() => pp.backend.Pantry(pp.sessionName).PUT(pantry)}
                    onSaveTxt='Desant...'

                    onAccept={() => pp.onBackToMenu()}
                    onAcceptTxt='Desat'

                    onRejectTxt='Error'
                />}
                logoOnClick={() => { pp.backend.Pantry(pp.sessionName).PUT(pantry).then(pp.onGotoHome) }}
                titleOnClick={() => setHelp(true)}
                titleText='El&nbsp;meu rebost'
                right={<SaveButton
                    key='save'
                    baseTxt='Següent'

                    onSave={() => pp.backend.Pantry(pp.sessionName).PUT(pantry)}
                    onSaveTxt='Desant...'

                    onAccept={() => pp.onComplete()}
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
                            needs.contents.map((i: ShoppingNeedsItem, idx: number) => (
                                <IngredientRow
                                    key={i.name}
                                    id={idx % 2 === 0 ? 'even' : 'odd'}
                                    item={i}
                                    onChange={(value: number) => {
                                        const c = pantry.contents.find(p => p.name === i.name)
                                        c && (c.amount = value)
                                        setPantry(pantry)
                                    }}
                                    onClick={() => {
                                        if (focussed) {
                                            setFocussed(undefined)
                                        } else {
                                            pp.backend.IngredientUse().POST({
                                                MenuName: pp.sessionName,
                                                IngredientName: i.name
                                            }).then((usage) => {
                                                setFocussed({
                                                    item: i,
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
                            <td id='right'>{asEuro(computeSavings())}</td>
                        </tr>
                    </tfoot>
                </table>
                {
                    focussed && <IngredientDialog
                        item={focussed.item}
                        usage={focussed.usage}
                        onClose={() => setFocussed(undefined)}
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

function merge<A, B>(a: A[], b: B[], cmp: ((a: A, b: B) => number), f: ((a: A | undefined, b: B | undefined) => number)): number {
    let i = 0
    let j = 0
    let acc = 0

    while (i < a.length && j < b.length) {
        switch (cmp(a[i], b[j])) {
            case -1:
                i++
                continue
            case 1:
                j++
                continue
            case 0:
                acc += f(a[i], b[j])
                i++
                j++
        }
    }

    for (; i < a.length; i++) {
        acc += f(a[i], undefined)
    }

    for (; j < b.length; j++) {
        acc += f(undefined, b[j])
    }

    return acc
}