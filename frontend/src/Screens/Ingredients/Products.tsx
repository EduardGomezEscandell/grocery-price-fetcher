import React, { useState } from 'react'
import { useNavigate } from 'react-router-dom';
import TopBar from '../../TopBar/TopBar'
import Sidebar from '../../SideBar/Sidebar'
import Backend from '../../Backend/Backend';
import ComparableString from '../../ComparableString/ComparableString';
import { asEuro, makePlural, round2 } from '../../Numbers/Numbers';
import { Product } from '../../State/State';
import './Products.css'

interface Props {
    backend: Backend;
    sessionName: string;
}

export default function Products(props: Props) {
    const [sideBar, setSidebar] = useState(false)
    const [help, setHelp] = useState(false)
    const navigate = useNavigate()

    const [products, setProducts] = useState<product[]>([])
    const [loaded, setLoaded] = useState(false)

    if (!loaded) {
        props.backend.Products(props.sessionName)
            .GET()
            .then((d) => d.map(r => new product(r.name, r.price, r.batch_size, r.provider)))
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
                titleOnClick={() => setHelp(true)}
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
                            result.map(r => <ProductRow product={r} />)
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
                {help && <HelpDialog onClose={() => setHelp(false)} />}
                {sideBar && <Sidebar onHelp={() => setHelp(true)} onNavigate={() => { props.backend.ClearCache() }} />}
            </section >
        </div >
    )
}

function ProductRow(props: { product: product }): JSX.Element {
    const { name, batch_size, price, provider } = props.product

    const text = round2(batch_size)
        + ' '
        + makePlural(batch_size, 'unitat', 'unitats')
        + ' a '
        + asEuro(price)

    return <div key={name} className='search-table-row'>
        <div className='title'>
            {name}
        </div>
        <div className='details'>
            <div>
                {provider}
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

    constructor(name: string, price: number, batch_size: number, provider: string) {
        super(name, price, batch_size, provider)
        this.comp = new ComparableString(name)
    }
}