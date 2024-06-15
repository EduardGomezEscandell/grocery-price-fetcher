import React, { useState } from 'react'
import Select from 'react-select'
import Backend from '../../Backend/Backend'
import RecipeEndpoint, { Ingredient, Recipe } from '../../Backend/endpoints/Recipe'
import { asEuro, round2 } from '../../Numbers/Numbers'
import './RecipeEditor.css'
import ProductsEndpoint from '../../Backend/endpoints/Products'
import { Product } from '../../State/State'

interface RecipeDialogProps {
    backend: Backend
    sessionName: string
    dish: string

    setHidden: () => void
}

export default function RecipeEditor(props: RecipeDialogProps): JSX.Element {
    const [expanded, setExpanded] = useState(false)
    const [editing, setEditing] = useState(false)

    return (
        <div className='recipe-editor'>
            <div key='header' id='header' onClick={() => {
                setExpanded(!expanded)
            }}>
                <div id='title'>
                    <span>{props.dish}</span>
                    {/* TODO: Make editable */}
                </div>
            </div>
            {expanded &&
                <CardBody
                    recipeEP={props.backend.Recipe(props.sessionName, props.dish)}
                    productsEP={props.backend.Products(props.sessionName)}
                    recipe={props.dish}
                    key={props.dish}
                    setEditing={setEditing}
                    setDeleted={props.setHidden}
                />
            }
        </div>
    )
}

interface CardBodyProps {
    recipeEP: RecipeEndpoint
    productsEP: ProductsEndpoint
    recipe: string

    setEditing: (editing: boolean) => void
    setDeleted: () => void
}

function CardBody(props: CardBodyProps): JSX.Element {
    const [ingredients, _setIngredients] = useState<Ingredient[]>([])
    const [loaded, setLoaded] = useState(false)
    const [total, _setTotal] = useState(0)
    const [backup, setBackup] = useState(new Recipe(props.recipe, ingredients))
    const [editing, _setEditing] = useState(false)
    const setEditing = (e: boolean) => { props.setEditing(e); _setEditing(e) }

    const setIngredients = (i: Ingredient[]) => {
        _setIngredients(i)
        _setTotal(i.reduce((acc, x) => acc + x.unit_price * x.amount, 0))
    }

    if (!loaded) {
        props.recipeEP
            .GET()
            .then((r) => setIngredients(r.ingredients))
            .then(() => setLoaded(true))
        
        return <div id='body' key='body'><div><h3>Descarregant ingredients...</h3></div></div>
    }

    return (
        <div id='body' key='body'>
            <div>
            <EditButtons
                key='buttons'
                onEdit={() => {
                    setBackup(deepNewIngredient(props.recipe, ingredients))
                    setEditing(true)
                }}
                onRestore={() => {
                    setIngredients(backup.ingredients)
                    setEditing(false)
                }}
                onSave={() => {
                    props.recipeEP
                        .PUT(new Recipe(props.recipe, ingredients))
                        .then(() => setIngredients(ingredients.filter(i => i.amount > 0)))
                        .catch(() => {
                            alert("No s'ha pogut desar")
                            setIngredients(backup.ingredients)
                        })
                        .finally(() => setEditing(false))
                }}
                onDelete={() => {
                    props.recipeEP
                        .DELETE()
                        .then(() => props.setDeleted(), () => alert("No s'ha pogut eliminar"))
                        .finally(() => setEditing(false))
                }}
                editing={false}
            />
            <h3>Ingredients</h3>
            {
                ingredients.map((ing, idx) => (
                    <IngredientRow
                        key={ing + idx.toString() + editing}
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
        </div>
    )
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
                    id={(selected && atof(amount)===0) ? 'error' : 'ok'}
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
                    placeholder='Selecciona...'
                />
            </div>
            <div id='price' key='price'>{selected && asEuro(atof(amount) * selected.price / selected.batch_size) || '0.00'}</div>
            <button id={(selected && atof(amount) != 0) ? 'add' : 'disabled'}
                onClick={() => {
                    if (!selected || atof(amount) === 0) {
                        return
                    }
                    props.onChange(new Ingredient(selected.name, selected.price / selected.batch_size, atof(amount)))
                    setSelected(null)
                }}
                style={{ width: '40px', fontSize: 'inherit' }}
            >
                +
            </button>

        </div >
    )
}

function deepNewIngredient(name: string, ingredients: Ingredient[]): Recipe {
    return new Recipe(name, ingredients.map(i => ({ ...i })))
}

function atof(s: string): number {
    return s === ''
        ? 0
        : parseFloat(s)
}