import React from 'react'
import { Ingredient } from '../../State/State.js'
import { asEuro, positive, round2, makePlural } from '../../Numbers/Numbers.ts'

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
            <dialog open>
                <h2 id="header">{this.props.ingredient.name}</h2>
                <div id="body">
                    <p>
                        Tens <b>{round2(this.props.ingredient.have)}</b> {makePlural(this.state.storage, "unitat", "unitats")} al teu rebost. En necessites{' '}
                        <b>{round2(this.props.ingredient.need)}</b>, i per tant te'n falten {' '}
                        <b>{round2(this.state.deficit)}</b>. Aquest producte es ven en
                        paquets de <b>{round2(this.props.ingredient.batch_size)}</b> {makePlural(this.props.ingredient.batch_size, "unitat", "unitats")},
                        i per tant has de comprar  <b>{round2(this.state.packs)}</b> {makePlural(this.state.packs, "paquet", "paquets")}.
                    </p>
                    <p>
                        Cada paquet costa <b>{asEuro(this.props.ingredient.price)}</b>, i per
                        tant et costarà <b>{asEuro(this.state.cost)}</b>
                    </p>
                </div>
                <div id="footer">
                    <button onClick={this.props.onClose}>OK</button>
                </div>
            </dialog>
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
                onMouseEnter={() => this.setState({ ...this.state, id: 'highlight' })}
                onMouseLeave={() => this.setState({ ...this.state, id: this.defaultID })}
                onClick={(e) => {
                    if (e.target instanceof HTMLInputElement) return
                    this.props.onClick(this)
                }}
            >
                <td id='left' key='name'> {this.props.ingredient.name}  </td>
                <td key='have'>
                    <input 
                        type="number"
                        value={this.state.storage}
                        onClick={(e) => { e.target instanceof HTMLInputElement && e.target.select() }}
                        onChange={(s) => this.onChange(s)}
                        datatype='number'
                    />
                </td>
                <td id='right' key='price-total'> {asEuro(this.state.cost)} </td>
            </tr>
        )
    }
}



