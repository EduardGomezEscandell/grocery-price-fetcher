import React, { useEffect, useState } from 'react'
import Backend from '../../Backend/Backend';
import TopBar from '../../TopBar/TopBar';
import { ShoppingList } from '../../State/State';
import SaveButton from '../../SaveButton/SaveButton';
import ShoppingItem, { Column } from './ShoppingItem';
import { asEuro } from '../../Numbers/Numbers';
import './ShoppingList.css';
import { useNavigate } from 'react-router-dom';
import Sidebar from '../../SideBar/Sidebar';
import DangerDialog from '../DangerDialog/DangerDialog';

interface Props {
    backend: Backend;
    sessionName: string;
}

enum Dialog {
    OFF,
    RESTORE,
    HELP,
    CLOSING,
}

export default function Shopping(props: Props): JSX.Element {
    const [dialog, _setDialog] = useState(Dialog.OFF);

    const [k, _setK] = useState(0)
    const forceReload = () => _setK(k + 1)

    const setDialog = (d: Dialog): boolean => {
        if (dialog === Dialog.CLOSING) return false
        _setDialog(d)
        return true
    }

    const [column, setColumn] = useState(Column.UNITS)

    const [shoppingList, setShoppingList] = useState<ShoppingList>(new ShoppingList())
    useEffect(() => {
        props.backend
            .ShoppingList(props.sessionName, props.sessionName)
            .GET()
            .then((sl) => setShoppingList(sl))
    }, [props])

    const tableStyle: React.CSSProperties = {}
    if (dialog !== Dialog.OFF) {
        tableStyle.filter = 'blur(5px)'
    }

    const navigate = useNavigate()

    const [sidebar, setSidebar] = useState(false)

    return (
        <div id='rootdiv'>
            <TopBar
                left={<button className='save-button' id='idle'
                    onClick={() => setSidebar(!sidebar)}
                >Opcions </button>
                }
                logoOnClick={() => saveShoppingList(props.backend, shoppingList).then(() => navigate("/"))
                }
                titleOnClick={() => setDialog(Dialog.HELP)}
                titleText="La&nbsp;meva compra"
                right={<SaveButton
                    key='save'

                    baseTxt='Desar'
                    onSave={() => saveShoppingList(props.backend, shoppingList)}
                    onSaveTxt='Desant...'
                    onAcceptTxt='Desat'
                    onRejectTxt='Error'
                />}
            />
            <section>
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
                                shoppingList.items.map((item, idx) =>
                                    <ShoppingItem
                                        item={item}
                                        idx={idx}
                                        key={`${k}-${idx}-${column}`}
                                        show={column}
                                        setSelection={(v: boolean) => {
                                            item.done = v
                                            setShoppingList(shoppingList)
                                        }}
                                    />
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
                                            const cost = shoppingList.items.reduce((acc, i) => acc + i.cost, 0)
                                            return (<>{asEuro(cost)}</>)
                                        })()
                                    }
                                </td>
                            </tr>
                        </tfoot>
                    </table>
                    {(dialog === Dialog.RESTORE || dialog === Dialog.CLOSING) &&
                        <DangerDialog
                            onAccept={() => {
                                if (!setDialog(Dialog.CLOSING)) {
                                    return
                                }
                                shoppingList.items.forEach(i => i.done = false)
                                props.backend
                                    .ShoppingList(props.sessionName, props.sessionName)
                                    .DELETE()
                                    .then(() => setDialog(Dialog.OFF))
                                    .then(() => forceReload())
                            }}
                            onReject={() => setDialog(Dialog.OFF)}
                        >
                            <h3 id='header'>Restaurar la llista de la compra?</h3>
                            <div id='body'>
                                <p>Tots els elements marcats com a comprats es desmarcaran</p>
                                <p>Prem <i>No</i> per a tornar enrere, prem <i>Sí</i> per a continuar</p>
                            </div>
                        </DangerDialog>
                    }{
                        dialog === Dialog.HELP &&
                        <HelpDialog
                            onClose={() => setDialog(Dialog.OFF)}
                        />
                    }
                </div>
                {sidebar && <Sidebar
                    onHelp={() => {
                        setDialog(Dialog.HELP)
                        setSidebar(false)
                    }}
                    onNavigate={() => saveShoppingList(props.backend, shoppingList)}
                />}
            </section>
        </ div>
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

async function saveShoppingList(backend: Backend, shoppingList: ShoppingList): Promise<void> {
    backend
        .ShoppingList(shoppingList.menu, shoppingList.pantry)
        .PUT(shoppingList.items.filter(i => i.done).map(i => i.name))
        .then(() => { })
}