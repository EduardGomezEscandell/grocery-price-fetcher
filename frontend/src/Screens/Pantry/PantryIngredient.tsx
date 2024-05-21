import React from 'react'
import { Ingredient } from '../../State/State.js'
import { asEuro, positive, round2, makePlural } from '../../Numbers/Numbers.ts'
import './Pantry.css'

interface Props {
    ingredient: Ingredient;
    onChange: (newHave: number) => void;
}

class State {
    storage: number
    deficit: number
    packs: number
    cost: number
}

class PantryIngredient<T extends Props, S extends State = State> extends React.Component<T> {
    state: S

    constructor(pp: T) {
        super(pp)
        const def = positive(pp.ingredient.need - pp.ingredient.have)
        const pks = Math.ceil(def / pp.ingredient.batch_size)

        this.state = {
            storage: pp.ingredient.have,
            deficit: def,
            packs: pks,
            cost: pks * pp.ingredient.price
        } as S
    }

    render(): JSX.Element {
        return <div>ERROR</div>
    }

    onChange(s: React.ChangeEvent<HTMLInputElement>) {
        const newStorage = positive(parseFloat(s.target.value))
        const newDeficit = positive(this.props.ingredient.need - newStorage)
        const newPackCount = Math.ceil(newDeficit / this.props.ingredient.batch_size)
        this.setState({
            ...this.state,
            storage: newStorage,
            deficit: newDeficit,
            packs: newPackCount,
            cost: newPackCount * this.props.ingredient.price
        })
        this.props.onChange(newStorage)
    }
}

interface FocusIngredientProps extends Props {
    onClose: () => void;
}

export class FocusIngredient extends PantryIngredient<FocusIngredientProps> {
    render(): JSX.Element {
        return (
            <div className="Dialog">
                <dialog open>
                    <h1>{this.props.ingredient.name}</h1>
                    <p>
                        Tens <span className='Amount'>{round2(this.props.ingredient.have)}</span> {makePlural(this.state.storage, "unitat", "unitats")} al teu rebost. En necessites{' '}
                        <span className='Amount'>{round2(this.props.ingredient.need)}</span>, i per tant te'n falten {' '}
                        <span className='Amount'>{round2(this.state.deficit)}</span>. Aquest producte es ven en
                        paquets de <span className='Amount'>{round2(this.props.ingredient.batch_size)}</span> {makePlural(this.props.ingredient.batch_size, "unitat", "unitats")},
                        i per tant has de comprar  <span className='Amount'>{round2(this.state.packs)}</span> {makePlural(this.state.packs, "paquet", "paquets")}.
                    </p>
                    <p>
                        Cada paquet costa <span className='Amount'>{asEuro(this.props.ingredient.price)}</span>, i per
                        tant et costarà <span className='Amount'>{asEuro(this.state.cost)}</span>
                    </p>
                    <div className='OK'>
                        <button onClick={this.props.onClose}>OK</button>
                    </div>
                </dialog>
            </div>
        )
    }
}

interface RowIngredientProps extends Props {
    id: string;
    onClick: (self: RowIngredient) => void;
}

class RowIngredientState extends State {
    id: string
}

export class RowIngredient extends PantryIngredient<RowIngredientProps, RowIngredientState> {
    defaultID: string;

    constructor(pp: RowIngredientProps) {
        super(pp)
        this.defaultID = pp.id
        this.state.id = pp.id
    }

    render(): JSX.Element {
        return (
            <tr key={this.props.ingredient.name}
                id={this.state.id || 'odd'}
                onMouseEnter={() => this.setState({ ...this.state, id: 'mouseover' })}
                onMouseLeave={() => this.setState({ ...this.state, id: this.defaultID })}
                onClick={(e) => {
                    if (e.target instanceof HTMLInputElement) return
                    this.props.onClick(this)
                }}
            >
                <td className='Label' key='name'> {this.props.ingredient.name}  </td>
                <td className='Select' key='have'>
                    <input
                        type="number"
                        value={this.state.storage}
                        onClick={(e) => { e.target instanceof HTMLInputElement && e.target.select() }}
                        onChange={(s) => this.onChange(s)}
                        datatype='number'
                    />
                </td>
                <td className='Number' key='price-total'> {asEuro(this.state.cost)} </td>
            </tr>
        )
    }
}



