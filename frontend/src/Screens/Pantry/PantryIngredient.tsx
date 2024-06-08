import React from 'react'
import { Ingredient } from '../../State/State.tsx'
import { positive } from '../../Numbers/Numbers.ts'
import { IngredientUsage } from '../../Backend/endpoints/IngredientUse.tsx';
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
    usage: IngredientUsage[];
    onClose: () => void;
}

export class FocusIngredient extends PantryIngredient<FocusIngredientProps> {
    render(): JSX.Element {
        return (
            <dialog open id='pantry-ingredient'>
                <h2 id="header">{this.props.ingredient.name}</h2>
                <div id="body">
                    <p>
                        L'ingredient <b>{this.props.ingredient.name}</b> apareix en els següents plats:
                    </p>
                    <div className='vert-scroll'>
                        <div className='scroll-table'>
                            <table>
                                <tbody>
                                    {this.props.usage.map((u, idx) =>
                                        <tr key={u.day + u.meal + u.dish} id={idx%2===0 ? 'even' : 'odd'}>
                                            <td id="left">{u.meal} de {u.day}</td>
                                            <td id="left">{u.dish}</td>
                                            <td id="right">{u.amount}</td>
                                        </tr>
                                    )}
                                </tbody>
                            </table>
                        </div>
                    </div>
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
    onClick: () => void;
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
                    this.props.onClick()
                }}
            >
                <td id='left' key='name'> {this.props.ingredient.name}  </td>
                <td id='right' key='have'>
                    <input
                        type="number"
                        value={this.state.storage}
                        onClick={(e) => { e.target instanceof HTMLInputElement && e.target.select() }}
                        onChange={(s) => this.onChange(s)}
                        datatype='number'
                    />
                </td>
            </tr>
        )
    }
}



