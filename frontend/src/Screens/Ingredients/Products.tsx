import React, { useState } from 'react'
import { useNavigate } from 'react-router-dom';
import TopBar from '../../TopBar/TopBar'
import Sidebar from '../../SideBar/Sidebar'
import Backend from '../../Backend/Backend';
import ComparableString from '../../ComparableString/ComparableString';
import { asEuro, makePlural, round2 } from '../../Numbers/Numbers';
import { Product } from '../../State/State';
import ProductEditor from './ProductEditor';
import './Products.css'

interface Props {
    backend: Backend;
    sessionName: string;
}

enum Dialog {
    None,
    Help,
    Editor
}

export default function Products(props: Props) {
    const [sideBar, setSidebar] = useState(false)
    const [focus, setFocus] = useState(Dialog.None)
    const navigate = useNavigate()

    const [products, setProducts] = useState<product[]>([])
    const [loaded, setLoaded] = useState(false)
    const [currProduct, setCurrProduct] = useState<product | null>(null)

    if (!loaded) {
        props.backend.Products(props.sessionName)
            .GET()
            .then((d) => d.map(r => new product(r.name, r.price, r.batch_size, r.provider, r.product_id)))
            .then(setProducts)
            .then(() => setLoaded(true))
    }

    const [query, setQuery] = useState(new ComparableString(''))
    const [hidden, setHidden] = useState<string[]>([])

    const result = products
        .filter(r => !hidden.includes(r.name))
        .filter((r) => r.comp.contains(query))

    return (
        <div id='rootdiv'>
            <TopBar
                left={<button onClick={() => setSidebar(!sideBar)}> Opcions </button>}
                right={<></>}
                titleText="Els meus productes"
                logoOnClick={() => {
                    props.backend.ClearCache()
                    navigate('/')
                }}
                titleOnClick={() => setFocus(Dialog.Help)}
            />
            <div className='search-table-search'>
                <input id={result.length === 0 ? 'error' : 'search'}
                    type='text'
                    placeholder='Cerca productes...'
                    value={query.displayName}
                    onChange={q => setQuery(new ComparableString(q.target.value))}
                />
            </div>
            <section>
                <div className='search-table'>
                    <div id='body' key={query.compareName}>
                        {
                            result.map(r =>
                                <ProductRow product={r} key={r.name} onClick={() => {
                                    setCurrProduct(r)
                                    setFocus(Dialog.Editor)
                                }} />
                            )
                        }
                        {
                            result.length === 0 &&
                            <div id='error'>
                                No hi ha resultats
                            </div>
                        }
                        <p></p>
                    </div>
                </div>
                {focus === Dialog.Help && <HelpDialog onClose={() => setFocus(Dialog.None)} />}
                {focus === Dialog.Editor && <ProductEditor
                    backend={props.backend}
                    sessionName={props.sessionName}
                    product={currProduct!}
                    onHide={() => { setHidden([...hidden, currProduct!.name]); setFocus(Dialog.None) }}
                    onChange={(p: Product) => {
                        props.backend.Products(props.sessionName).POST(currProduct!.name, p)
                        currProduct!.override(p)
                    }}
                    onClose={() => setFocus(Dialog.None)}
                />}
                {sideBar && <Sidebar onHelp={() => setFocus(Dialog.Help)} onNavigate={() => { props.backend.ClearCache() }} />}
            </section >
        </div >
    )
}

function ProductRow(props: { product: product, onClick: () => void }): JSX.Element {
    const { name, batch_size, price, provider, product_id: provider_id } = props.product

    const text = round2(batch_size)
        + ' '
        + makePlural(batch_size, 'unitat', 'unitats')
        + ' a '
        + asEuro(price)

    return <div key={name} className='search-table-row' onClick={props.onClick}>
        <div className='title'>
            {name}
        </div>
        <div className='details'>
            <div>
                {provider} #{provider_id}
            </div>
            <div>
                {text}
            </div>
            <div>
                {asEuro(price / batch_size)}/u
            </div>
        </div>
    </div>
}

function HelpDialog(props: { onClose: () => void }): JSX.Element {
    return (
        <dialog open>
            <h2 id="header">Els meus productes</h2>
            <div id="body">
                <p>
                    Aquesta pàgina pàgina et permet veure i editar els teus productes, i crear-ne de nous.
                </p>
            </div>
            <div id="footer">
                <button onClick={props.onClose}>
                    D'acord
                </button>
            </div>
        </dialog>
    )
}

class product extends Product {
    comp: ComparableString;

    constructor(name: string, price: number, batch_size: number, provider: string, provider_id: string) {
        super(name, price, batch_size, provider, provider_id)
        this.comp = new ComparableString(name)
    }

    override(p: Product) {
        this.name = p.name
        this.price = p.price
        this.batch_size = p.batch_size
        this.provider = p.provider
        this.product_id = p.product_id
        this.comp = new ComparableString(p.name)
    }
}