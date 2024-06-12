import React, { useEffect, useState } from 'react'
import { Pantry, PantryItem, ShoppingNeeds, ShoppingNeedsItem } from '../../State/State';
import Backend from '../../Backend/Backend';
import TopBar from '../../TopBar/TopBar';
import SaveButton from '../../SaveButton/SaveButton';
import IngredientRow from './PantryIngredient';
import IngredientDialog from './IngredientDialog';
import { IngredientUsage } from '../../Backend/endpoints/IngredientUse';
import { useNavigate } from 'react-router-dom';
import Sidebar from '../../SideBar/Sidebar';
import DangerDialog from '../DangerDialog/DangerDialog';
import { PantryEndpoint } from '../../Backend/endpoints/Pantry';

interface Props {
    backend: Backend;
    sessionName: string;
}

interface Focus {
    item: PantryItem
    usage: IngredientUsage[]
}

enum Dialog {
    OFF,
    HELP,
    FOCUS,
    CONFIRM,
}

export default function RenderPantry(pp: Props) {
    const [pantry, setPantry] = useState<Pantry>(new Pantry(pp.sessionName))
    const [focussed, setFocussed] = useState<Focus>({ item: new PantryItem("", 0), usage: [] })
    const navigate = useNavigate()

    const [dialog, setDialog] = useState(Dialog.OFF)

    const tableStyle: React.CSSProperties = {}
    if (dialog !== Dialog.OFF) {
        tableStyle.filter = 'blur(5px)'
    }

    useEffect(() => {
        Promise.all([
            pp.backend.Pantry(pp.sessionName)
                .GET()
                .then(p => p, () => new Pantry(pp.sessionName)),
            pp.backend.Needs(pp.sessionName).GET(),
        ])
            .then(([pantry, needs]) => filterPantry(pantry, needs))
            .then(p => setPantry(p))
            .catch((reason) => {
                console.log('Error getting pantry: ', reason || 'Unknown error')
            })
    }, [pp.backend, pp.sessionName])

    const [sidebar, setSidebar] = useState(false)

    return (
        <div id='rootdiv'>
            <TopBar
                left={<button className='save-button' id='idle'
                    onClick={() => setSidebar(!sidebar)}
                >Opcions </button>
                }
                logoOnClick={() => { pp.backend.Pantry(pp.sessionName).PUT(pantry).then(() => navigate("/")) }}
                titleOnClick={() => setDialog(Dialog.HELP)}
                titleText='El&nbsp;meu rebost'
                right={<SaveButton
                    key='save'
                    baseTxt='Següent'

                    onSave={() => pp.backend.Pantry(pp.sessionName).PUT(pantry)}
                    onSaveTxt='Desant...'

                    onAccept={() => navigate("/shopping-list")}
                    onAcceptTxt='Desat'

                    onReject={(reason: any) => console.log('Error saving pantry: ', reason || 'Unknown error')}
                    onRejectTxt='Error'
                />}
            />
            <section>
                <div className='scroll-table' key='pantry'>
                    <table style={tableStyle}>
                        <thead>
                            <tr key='header' id='header1'>
                                <th id="left">
                                    <button id='X' onClick={async () => setDialog(Dialog.CONFIRM)}>X</button>
                                    Producte
                                </th>
                                <th id="right">Tens</th>
                            </tr>
                        </thead>
                        <tbody>
                            {
                                pantry.contents.map((i: ShoppingNeedsItem, idx: number) => (
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
                                            pp.backend.IngredientUse(pp.sessionName, i.name)
                                                .GET()
                                                .then(usage => setFocussed({ item: i, usage: usage }))
                                                .then(() => setDialog(Dialog.FOCUS))
                                                .catch(reason => console.log('Error getting ingredient usage: ', reason || 'Unknown error'))

                                        }}
                                    />
                                ))
                            }
                        </tbody>
                    </table>
                    {renderDialog(
                        { state: dialog, setState: setDialog },
                        focussed,
                        pp.backend.Pantry(pp.sessionName)
                    )}
                </div>
                {sidebar && <Sidebar
                    onHelp={() => {
                        setDialog(Dialog.HELP)
                        setSidebar(false)
                    }}
                    onNavigate={() => pp.backend.Pantry(pp.sessionName).PUT(pantry)}
                />}
            </section>
        </ div>
    )
}

interface Stateful<T> {
    state: T
    setState: (s: T) => void
}

function renderDialog(d: Stateful<Dialog>, focus: Focus, pantryEP: PantryEndpoint) {
    switch (d.state) {
        case Dialog.OFF:
            return null
        case Dialog.FOCUS:
            return <IngredientDialog
                item={focus.item}
                usage={focus.usage}
                onClose={() => d.setState(Dialog.OFF)}
            />
        case Dialog.HELP:
            return <dialog open>
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
                    <button onClick={() => d.setState(Dialog.OFF)}>
                        D'acord
                    </button>
                </div>
            </dialog>
        case Dialog.CONFIRM:
            return <DangerDialog
                onAccept={() => {
                    pantryEP.DELETE().then(() => window.location.reload())
                }}
                onReject={() => d.setState(Dialog.OFF)}
            >
                <h3 id='header'>Confirmació</h3>
                <div id='body'>Segur que vols esborrar el rebost?</div>
            </DangerDialog>
        default:
            return null
    }
}

// This function is used to filter the pantry contents against the shopping needs
// - Items inherit their amounts from the pantry.
// - If an item is in the pantry but not in the needs, it is removed.
// - If an item is in the needs the amount defaults to 0.
// - Items are sorted alphabetically.
function filterPantry(pantry: Pantry, needs: ShoppingNeeds): Pantry {
    const filtered = new Pantry(pantry.name)
    pantry.contents.sort((a, b) => a.name.localeCompare(b.name))
    needs.items.sort((a, b) => a.name.localeCompare(b.name))

    let i = 0;
    let j = 0;

    while (i < pantry.contents.length && j < needs.items.length) {
        const comp = pantry.contents[i].name.localeCompare(needs.items[j].name)
        if (comp < 0) {
            // Ingredient in pantry but not in needs
            i++
        } else if (comp > 0) {
            // Ingredient in needs but not in pantry
            filtered.contents.push({
                name: needs.items[j].name,
                amount: 0,
            })
            j++
        } else {
            // Ingredient in both pantry and needs
            filtered.contents.push(pantry.contents[i])
            i++
            j++
        }
    }

    while (j < needs.items.length) {
        filtered.contents.push({
            name: needs.items[j].name,
            amount: 0,
        })
        j++
    }
    return filtered
}