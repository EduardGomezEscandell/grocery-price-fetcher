import { useState } from "react";
import { Product } from "../../State/State";
import { asEuro, parseNumber, round2 } from "../../Numbers/Numbers";
import ProviderEndpoint from "../../Backend/endpoints/Provider";
import Backend from "../../Backend/Backend";
import './ProductEditor.css'


interface Props {
    backend: Backend;
    sessionName: string;
    product: Product;
    onHide: () => void;
    onClose: () => void;
    onChange: (newProduct: Product) => void;
}

enum Stage {
    Name,
    Provider,
    ID,
    BatchSize,
    Final,
}

export default function ProductEditor(props: Props) {
    const [found, _setFound] = useState(props.product.price != 0.0)
    const [phase, _setPhase] = useState(computePhase(props.product, found))
    const [p, _setProduct] = useState(props.product)
    const [confirmDelete, setConfirmDelete] = useState(false)

    const setProduct = (newProduct: Product, newFound?: boolean) => {
        _setProduct(newProduct)
        if (newFound !== undefined) {
            _setFound(newFound)
        } else {
            newFound = found
        }

        const newPhase = computePhase(newProduct, newFound)
        newPhase != phase && _setPhase(newPhase)
    }

    if (confirmDelete) {
        return (
            <dialog open className="product-editor">
                <h2 id="header">Editor de productes</h2>
                <div id='body'>
                    <div id='search'>
                        <h3>Estàs segur que vols eliminar el següent producte?</h3>
                        <div>
                            <span>Nom</span>
                            <p>{p.name}</p>
                        </div>
                        <div>
                            <span>Proveidor</span>
                            <p>{p.provider}</p>
                        </div>
                        <div>
                            <span>Codi</span>
                            <p>{p.product_code}</p>
                        </div>
                        <div>
                            <span>Unitats a cada paquet</span>
                            <p>{p.batch_size}</p>
                        </div>
                    </div>
                </div>
                <div id='footer'>
                    <button id='dialog-left' onClick={() => { setConfirmDelete(false) }}>Cancel·lar</button>
                    <button id='dialog-right' onClick={() => {
                        props.backend.Products(props.sessionName).DELETE(p.id).then(
                            () => {
                                props.onHide()
                                props.onClose()
                            },
                            (e) => {
                                setConfirmDelete(false)
                                if (e instanceof Response) {
                                    e.text().then((t) => alert(`No s'ha pogut eliminar el producte.\n${e.status} ${t}`))
                                } else {
                                    alert(`No s'ha pogut eliminar el producte.\n${e}`)
                                }   
                            })
                    }}>Eliminar</button>
                </div>
            </dialog>
        )
    }

    return (
        <dialog open className="product-editor">
            <h2 id="header">Editor de productes</h2>
            <div id='body'>
                <div id='search'>
                    <div>
                        <span>Nom</span>
                        <input value={p.name} onChange={e => setProduct({ ...p, name: e.target.value })} />
                    </div>
                    {phase >= Stage.Provider &&
                        <div>
                            <span>Proveidor</span>
                            <select value={p.provider} onChange={e => setProduct({ ...p, provider: e.target.value })}>
                                <option value=""></option>
                                <option value="Bonpreu">Bonpreu</option>
                                <option value="Mercadona">Mercadona</option>
                            </select>
                        </div>
                    }
                    {
                        phase >= Stage.ID &&
                        <div>
                            <span>Codi</span>
                            <input value={p.product_code} onChange={e => setProduct({ ...p, product_code: e.target.value })} />
                            {<ProviderLink
                                provider={p.provider}
                                providerId={p.product_code}
                                key={`url+${p.provider}+${p.product_code}`}
                            />}
                            <ProviderIDHelper provider={p.provider} />
                            <SearchProduct
                                key={`search+${p.provider}+${p.product_code}`}
                                product={p}
                                onSearch={(price: number, found: boolean) => {
                                    if (found) {
                                        setProduct({ ...p, price: price }, true)
                                    } else {
                                        setProduct(p, false)
                                    }
                                }}
                                endpoint={props.backend.Provider()}
                            />
                        </div>
                    }
                    {
                        phase >= Stage.BatchSize &&
                        <div>
                            <div>
                                <span>Unitats a cada paquet</span>
                                <input type='number'
                                    value={round2(p.batch_size)}
                                    onChange={e => setProduct({ ...p, batch_size: Math.floor(100 * parseNumber(e.target.value)) / 100 })}
                                />
                            </div>
                            <BatchSizeHelper />
                        </div>
                    }
                    {
                        phase >= Stage.Final &&
                        <div>
                            <h3>Resultat</h3>
                            <div id='prices'>
                                <div>
                                    <span>Preu</span>
                                    <p>{asEuro(p.price)}</p>
                                </div>
                                <div>
                                    <span>Preu per unitat</span>
                                    <p>{asEuro(p.price / p.batch_size)}/u</p>
                                </div>
                            </div>
                        </div>
                    }
                </div>
            </div>
            <div id='footer'>
                <button id='dialog-left' onClick={props.onClose}>Tornar sense desar</button>
                <button id='dialog-center' onClick={() => setConfirmDelete(true)}>El·liminar producte</button>
                <button id='dialog-right' onClick={() => {
                    props.onChange(p)
                    props.onClose()
                }
                }>Desar i tornar</button>
            </div>
        </dialog>
    )
}

function computePhase(product: Product, found: boolean): Stage {
    if (!product.name) {
        return Stage.Name
    }
    if (!product.provider) {
        return Stage.Provider
    }
    if (!found) {
        return Stage.ID
    }
    if (!product.batch_size) {
        return Stage.BatchSize
    }
    return Stage.Final
}

function ProviderLink(props: { provider: string, providerId: string }): JSX.Element {
    if (!props.provider || !props.providerId) {
        return <></>
    }

    const { URL, displayURL } = providerURLS(props.provider, props.providerId)
    return <a href={URL} target="_blank" rel="noopener noreferrer">
        {displayURL}
    </a>
}

function appleCode(provider: string): string {
    // Code for an apple in the provider's system (used as an example in the help dialog)
    switch (provider) {
        case 'Bonpreu':
            return '90041'
        case 'Mercadona':
            return '8177'
        default:
            return ''
    }
}

function providerURLS(provider: string, providerId: string): { URL: string, displayURL: string } {
    switch (provider) {
        case 'Bonpreu':
            return {
                URL: `https://www.compraonline.bonpreuesclat.cat/products/${providerId}/details`,
                displayURL: `compraonline.bonpreuesclat.cat/products/${providerId}/details`,
            }
        case 'Mercadona':
            providerId = providerId || '8177' // Default to an apple
            return {
                URL: `https://tienda.mercadona.es/product/${providerId}/`,
                displayURL: `tienda.mercadona.es/product/${providerId}/`,
            }
        default:
            return {
                URL: '',
                displayURL: '',
            }
    }
}

interface searchResult {
    price?: number;
    error?: boolean;
    pending: boolean;
}

interface SearchProductProps {
    product: Product;
    onSearch: (price: number, ok: boolean) => void;
    endpoint: ProviderEndpoint;
}

function SearchProduct(props: SearchProductProps): JSX.Element {
    const [searched, setSearched] = useState<searchResult>({ pending: true })
    const [title, setTitle] = useState('Cerca')

    const button = <button
        onClick={() => {
            setTitle('Cercant...')
            props.endpoint.GET(props.product)
                .then(p => {
                    setSearched({ price: p, pending: false })
                    props.onSearch(p, true)
                }, () => {
                    setSearched({ error: true, pending: false })
                    props.onSearch(0, false)
                })
                .finally(() => setTitle('Cerca'))
        }}
    >
        {title}
    </button>

    if (searched.pending) {
        return <>
            {button}
        </>
    }

    if (searched.error) {
        return <>
            {button}
            <p id='error'>
                No s'ha trobat cap producte a {props.product.provider} amb el codi <b>{props.product.product_code}</b>. Assegura't que el codi sigui correcte.
            </p>
        </>
    }

    return <>
        {button}
        <p id='success'>
            Producte trobat.
        </p>
    </>
}

function ProviderIDHelper(props: { provider: string }): JSX.Element {
    const exampleID = appleCode(props.provider)
    const [expanded, setExpanded] = useState(false)

    return <>
        <p id='expandable' onClick={() => setExpanded(!expanded)}>
            Com saber el codi del producte?
        </p>
        {
            expanded &&
            <span id='help'>
                <p>
                    Per coneixer el codi de qualsevol producte, visita la pàgina web del proveidor i busca el producte desitjat.
                    Mira l'URL de la pàgina i copia l'últim número que apareix a la URL.
                </p>
                <p>
                    Per exemple, per a una poma, l'URL seria
                </p>
                <ProviderLink provider={props.provider} providerId={exampleID} />
                <p>
                    i per tant el codi
                    seria <b>{exampleID}</b>.
                </p>
                <p id='expandable-exit' onClick={() => setExpanded(false)}>
                    D'acord
                </p>
            </span>
        }
    </>
}

function BatchSizeHelper(): JSX.Element {
    const [expanded, setExpanded] = useState(false)

    return <>
        <p id='expandable' onClick={() => setExpanded(!expanded)}>
            Què és aquest número?
        </p>
        {
            expanded &&
            <span id='help'>
                <p>
                    Aquest número es refereix a la quantitat de productes que venen en un paquet. Mira la descripció del producte
                    per saber quantes unitats venen en un paquet.
                </p>
                <p>
                    Per exemple, el valor típic pels ous es 6 ó 12.
                </p>
                <p id='expandable-exit' onClick={() => setExpanded(false)}>
                    D'acord
                </p>
            </span>
        }
    </>
}