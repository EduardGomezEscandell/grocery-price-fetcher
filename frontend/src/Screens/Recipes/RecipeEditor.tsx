import React, { useState } from 'react'
import Select from 'react-select'
import Backend from '../../Backend/Backend'
import RecipeEndpoint, { Ingredient, Recipe } from '../../Backend/endpoints/Recipe'
import { asEuro, round2 } from '../../Numbers/Numbers'
import './RecipeEditor.css'
import ProductsEndpoint from '../../Backend/endpoints/Products'
import { Product } from '../../State/State'

interface Props {
    backend: Backend
    sessionName: string
    dish: Dish;

    setHidden: () => void
    onRename: (r: string) => void
}

interface Dish {
    name: string
    id: number
}

export default function RecipeEditor(props: Props): JSX.Element {
    const [folded, setFolded] = useState(true)
    const [title, _setTitle] = useState<string>(props.dish.name)
    const setTitle = (t: string) => {
        props.onRename(t)
        _setTitle(t)
    }

    if (folded) {
        return (
            <div
                className='search-table-row'
                key={'recipe-editor'}
                id='folded'
                onClick={() => setFolded(!folded)}
            >
                <div className='title'>
                    <span>{title}</span>
                </div>
            </div>
        )
    }

    return (
        <div
            className='search-table-row'
            key={'recipe-editor'}
            id='expanded'
        >
            <RecipeCard
                recipeEP={props.backend.Recipe(props.sessionName, props.dish.id)}
                productsEP={props.backend.Products(props.sessionName)}
                dish={{id: props.dish.id, name: title}}
                key={title}
                setTitle={setTitle}
                setDeleted={props.setHidden}
                setFolded={() => setFolded(true)}
            />
        </div>
    )
}

interface RecipeCardProps {
    recipeEP: RecipeEndpoint
    productsEP: ProductsEndpoint

    dish: Dish
    setTitle: (r: string) => void
    setFolded: () => void
    setDeleted: () => void
}

function RecipeCard(props: RecipeCardProps): JSX.Element {
    const [title, setTitle] = useState(props.dish.name)
    const [ingredients, _setIngredients] = useState<Ingredient[]>([])
    const [loaded, setLoaded] = useState(false)
    const [total, _setTotal] = useState(0)
    const [backup, setBackup] = useState(new Recipe(props.dish.id, props.dish.name, ingredients))
    const [editing, setEditing] = useState(false)
    const [deletePage, setDeletePage] = useState(false)

    const setIngredients = (i: Ingredient[]) => {
        _setIngredients(i)
        _setTotal(i.reduce((acc, x) => acc + x.unit_price * x.amount, 0))
    }

    if (!loaded) {
        props.recipeEP
            .GET()
            .then((r) => {
                setBackup(deepNewRecipe(r.id, r.name, r.ingredients))
                setIngredients(r.ingredients)
            })
            .then(() => setLoaded(true))

        return <>
            <div className='title' onClick={props.setFolded}><div>{title}</div></div>
            <div id='body' key='body'><div><p>Descarregant ingredients...</p></div></div>
        </>
    }

    if (deletePage) {
        return (
            <span key='box'>
                <div className='title'><div>{title}</div></div>
                <div id='body' key='body'>
                    <div>
                        Segur que vols eliminar la recepta?
                        <div key='buttons' id='buttons'>
                            <button id='happy' onClick={() => {
                                setDeletePage(false)
                            }
                            }>No</button>
                            <button id='delete' onClick={() => {
                                props.recipeEP
                                    .DELETE()
                                    .then(() => props.setDeleted(), () => alert("No s'ha pogut eliminar"))
                                    .finally(() => setEditing(false))
                            }}>Sí</button>
                        </div>
                    </div>
                </div>
            </span>
        )
    }

    return (
        <span key='box' id={editing ? 'editing' : undefined}>
            <div className='title' onClick={() => {
                if (editing) {
                    return
                }
                props.setTitle(title)
                props.setFolded()
            }}>
                {
                    editing
                        ? <input
                            key='title'
                            value={title}
                            onChange={(e) => {
                                setTitle(e.target.value)
                            }}
                        />
                        : <div>{title}</div>
                }
            </div>
            <div className='body' key='body' >
                <div onClick={() => {
                    if (editing) {
                        return
                    }
                    props.setTitle(title)
                    props.setFolded()
                }}>
                    {
                        ingredients.map((ing, idx) => (
                            <IngredientRow
                                key={ing.id.toString() + idx.toString() + editing}
                                ingredient={ing}
                                editing={editing}
                                onChange={(newData: Ingredient) => {
                                    setIngredients(ingredients.map((x, i) => i === idx ? newData : x))
                                }} />
                        ))
                    }
                    {
                        editing &&
                        <NewIngredientRow
                            key='new-ingredient'
                            productsEP={props.productsEP}
                            onChange={(newData: Ingredient) => {
                                setIngredients([...ingredients, newData])
                            }} />
                    }
                    <span id='total'>
                        {loaded &&
                            <IngredientRow key={'total-' + total.toFixed()} ingredient={{
                                id: -1,
                                name: 'Total',
                                unit_price: total,
                                amount: NaN
                            }}
                                editing={false}
                                isTotal={true}
                                onChange={() => { }} />
                        }
                    </span>
                </div>
                <EditButtons
                    key='buttons'
                    onEdit={() => {
                        console.log('Editing')
                        setBackup(deepNewRecipe(props.dish.id, props.dish.name, ingredients))
                        setEditing(true)
                    }}
                    onRestore={() => {
                        console.log('Restoring')
                        setTitle(backup.name)
                        setIngredients(backup.ingredients)
                        setEditing(false)
                    }}
                    onSave={() => {
                        console.log("Saving")
                        saveRecipe(props.recipeEP, new Recipe(props.dish.id, title, ingredients), backup.name)
                            .then((r) => {
                                console.log(`Recipe saved: ${backup.name} -> ${r.name}`)
                                setTitle(r.name)
                                setIngredients(r.ingredients)
                                if(backup.name !== r.name) {
                                    props.setTitle(r.name)
                                }
                            }, (e) => {
                                console.log("Failure: ", e)
                                if (e instanceof Response) {
                                    e.text().then(
                                        (t) => alert(`No s'ha pogut desar:\nError ${e.status}. ${t}`),
                                        () => {
                                            setTitle(backup.name)
                                            setIngredients(backup.ingredients)
                                            alert("No s'ha pogut desar")
                                        }
                                    )
                                }
                                setTitle(backup.name)
                                setIngredients(backup.ingredients)
                            })
                            .finally(() => setEditing(false))
                    }}
                    onDelete={() => {
                        console.log('Deleting')
                        setTitle(backup.name)
                        setIngredients(backup.ingredients)
                        setEditing(false)
                        setDeletePage(true)
                    }}
                    editing={false}
                />
            </div>
        </span>
    )
}

const HTTPStatusConflict = 409

// Save the recipe, retrying if the name is already taken
// Returns the saved recipe
async function saveRecipe(recipeEP: RecipeEndpoint, r: Recipe, altName: string | null): Promise<Recipe> {
    return recipeEP
        .POST(r)
        .then(
            () => {
                return r
            }, (e) => {
                if (altName && e instanceof Response && e.status === HTTPStatusConflict) {
                    alert("Ja existeix una recepta amb aquest nom")
                    return saveRecipe(recipeEP, new Recipe(r.id, altName, r.ingredients), null)
                }
                return Promise.reject(e)
            })
}

interface ButtonsProps {
    onEdit: () => void
    onRestore: () => void
    onSave: () => void
    onDelete: () => void
    editing: boolean
}

function EditButtons(props: ButtonsProps): JSX.Element {
    const [editing, setEditing] = useState(props.editing)

    return (
        <div key='buttons' id='buttons'>
            {
                editing
                    ? <button
                        id='revert' key='revert'
                        onClick={() => {
                            setEditing(false)
                            props.onRestore()
                        }}>
                        Cancel.la
                    </button>
                    : <button id='happy' key='happy'
                        onClick={() => {
                            setEditing(true)
                            props.onEdit()
                        }}>
                        Edita
                    </button>
            }
            {
                editing && <button
                    id='happy' key='happy'
                    onClick={() => {
                        setEditing(false)
                        props.onSave()
                    }}>
                    Desa
                </button>
            }
            <button
                id='delete' key='delete'
                onClick={props.onDelete}>
                Elimina
            </button>
        </div>
    )
}

interface InfredientProps {
    ingredient: Ingredient
    editing: boolean
    isTotal?: boolean
    onChange: (i: Ingredient) => void
}

function IngredientRow(props: InfredientProps): JSX.Element {
    const [amount, setAmount] = useState(round2(props.ingredient.amount))
    const [hidden, setHidden] = useState(false)
    const isTotal = props.isTotal || false

    if (hidden) {
        return <></>
    }

    return (
        <div className='ingredient' key={props.ingredient.name}>
            <div id='amount' key='amount'>
                {props.editing
                    && <input
                        id={atof(amount) === 0 ? 'error' : ''}
                        key='amount'
                        value={amount}
                        onChange={(e) => {
                            if (e.target.value === '') {
                                props.ingredient.amount = 0
                            } else {
                                props.ingredient.amount = atof(e.target.value)
                            }
                            if (props.ingredient.amount != 0) {
                                props.onChange(props.ingredient)
                            }
                            setAmount(e.target.value)
                        }}
                        type='number'
                    />
                    || <span>{isTotal ? '' : amount}</span>
                }
            </div>
            <div id='name' key='name'>
                {props.ingredient.name}
            </div>
            <div id='price' key='price'>{asEuro(
                (isTotal ? 1.0 : props.ingredient.amount) *
                props.ingredient.unit_price)}
            </div>
            {props.editing && <button id='remove' key='remove'
                onClick={() => {
                    setAmount('0')
                    props.onChange({ ...props.ingredient, amount: 0 })
                    setHidden(true)
                }}
                style={{ width: '40px', fontSize: 'inherit' }}
            >
                x
            </button>}
        </div>
    )
}

interface NewIngredientProps {
    productsEP: ProductsEndpoint
    onChange: (i: Ingredient) => void
}

interface selectItem {
    value: Product
    label: string
}

function NewIngredientRow(props: NewIngredientProps): JSX.Element {
    const [amount, setAmount] = useState('0')
    const [selected, setSelected] = useState<Product | null>(null)
    const [prods, setProducts] = useState<selectItem[]>([])
    const [loaded, setLoaded] = useState(false)

    if (!loaded) {
        props.productsEP
            .GET()
            .then(prods => prods.map(p => ({ value: p, label: p.name })).sort((a, b) => a.label.localeCompare(b.label)))
            .then(setProducts)
            .then(() => setLoaded(true))
    }

    return (
        <div className='ingredient new-ingredient' key='new-ingredient'>
            <div id={'amount'} key='amount'>
                <input
                    key='amount'
                    id={(selected && atof(amount) === 0) ? 'error' : 'ok'}
                    value={amount}
                    onChange={(e) => {
                        setAmount(e.target.value)
                    }}
                    type='number'
                />
            </div>
            <div id='name' key='name'>
                <Select
                    className='Select'
                    styles={{
                        control: (base) => ({
                            ...base,
                            fontSize: 'inherit',
                        }),
                        menu: (base) => ({
                            ...base,
                            font: 'inherit',
                            fontSize: 'inherit',
                        }),
                        container: (base) => ({
                            ...base,
                            padding: 0,
                            margin: 0,
                            width: '100%',
                        }),
                        input: (base) => ({
                            ...base,
                            padding: 0,
                            margin: 0,
                        }),
                    }}
                    components={{ DropdownIndicator: () => null, IndicatorSeparator: () => null }}
                    onChange={selected => {
                        if (selected == null) {
                            return setSelected(null)
                        }
                        setSelected(selected.value)
                    }}
                    value={selected ? { value: selected, label: selected.name } : null}
                    options={prods}
                    isSearchable
                    placeholder='Afegeix...'
                />
            </div>
            <div id='price' key='price'>{selected && asEuro(atof(amount) * selected.price / selected.batch_size) || '0.00'}</div>
            <button id={(selected && atof(amount) != 0) ? 'add' : 'disabled'}
                onClick={() => {
                    if (!selected || atof(amount) === 0) {
                        return
                    }
                    props.onChange(new Ingredient(selected.id, selected.name, selected.price / selected.batch_size, atof(amount)))
                    setSelected(null)
                }}
                style={{ width: '40px', fontSize: 'inherit' }}
            >
                +
            </button>

        </div >
    )
}

function deepNewRecipe(id: number, name: string, ingredients: Ingredient[]): Recipe {
    return new Recipe(id, name, ingredients.map(i => ({ ...i })))
}

function atof(s: string): number {
    return s === ''
        ? 0
        : parseFloat(s)
}