import React from 'react'
import Backend from '../../Backend/Backend.ts';
import Optional from '../../Optional/Optional.ts';
import { State, Day, Meal, Dish, Menu } from '../../State/State.tsx';
import TopBar from '../../TopBar/TopBar.tsx';
import DishPicker from './DishPicker.tsx'
import './Menu.css'
import { round2 } from '../../Numbers/Numbers.ts';
import SaveButton from '../../SaveButton/SaveButton.tsx';
import DownloadPantry from '../Pantry/PantryLoad.ts';

interface Props {
    backend: Backend;
    globalState: State;
    onComplete: () => void
    onGotoHome: () => void
}

export default class MenuTable extends React.Component<Props> {
    state: {
        days: string[],
        mealSizes: number[]
        focus: { day: Day, meal: Meal } | undefined
        help: boolean
        hover: string | undefined
    }

    constructor(props: Props) {
        super(props)
        this.state = {
            focus: undefined,
            hover: undefined,
            help: false,
            days: props.globalState.menu.days.map(d => d.name),
            mealSizes: this.computeMealSizes(props.globalState.menu)
        }
    }

    get days(): string[] {
        return this.state.days
    }

    render(): JSX.Element {
        const tableStyle: React.CSSProperties = {}
        if (this.state.focus !== undefined || this.state.help) {
            tableStyle.filter = 'blur(5px)'
        }

        return (
            <>
                <TopBar
                    left={<SaveButton
                        key='goback'

                        baseTxt='Tornar'

                        onSave={() => saveMenu(this.props.backend, this.props.globalState)}
                        onSaveTxt='Desant...'

                        onAcceptTxt='Desat'
                        onAccept={this.props.onGotoHome}

                        onRejectTxt='Error'
                    />}
                    logoOnClick={() => saveMenu(this.props.backend, this.props.globalState).then(this.props.onGotoHome)}
                    titleOnClick={() => this.DisplayHelp()}
                    titleText='El&nbsp;meu menú'
                    right={<SaveButton
                        key='save'

                        baseTxt='Següent'

                        onSave={() => DownloadPantry(this.props.backend, this.props.globalState)}
                        onSaveTxt='Desant...'

                        onAcceptTxt='Desat'
                        onAccept={this.props.onComplete}

                        onRejectTxt='Error'

                    />}
                />
                <div className='Menu'>
                    <table key='menu-table' style={tableStyle}>
                        <tbody>
                            <tr>
                                {
                                    this.props.globalState.menu.days.map((day, i) =>
                                        <td key={`day-col-${i}`}>
                                            {this.DayCol(day)}
                                        </td>
                                    )
                                }

                            </tr>
                        </tbody>
                    </table>
                    {this.RenderFocus()}
                    {this.RenderHelp()}
                </div>
            </>
        )
    }

    private Focus(day: Day, meal: Meal) {
        this.setState({
            ...this.state,
            help: false,
            focus: {
                day: day,
                meal: meal
            }
        })
    }

    private Unfocus() {
        this.setState({
            ...this.state,
            focus: undefined
        })
    }

    private Highlight(dish: Dish) {
        if (this.state.focus !== undefined) {
            return
        }
        if (this.state.help) {
            return
        }
        this.setState({
            ...this.state,
            hover: dish.name
        })
    }

    private Unhighlight() {
        if (this.state.focus !== undefined) {
            return
        }
        if (this.state.help) {
            return
        }
        this.setState({
            ...this.state,
            hover: undefined
        })
    }

    private DisplayHelp() {
        if (this.state.focus !== undefined) {
            return
        }
        this.setState({
            ...this.state,
            highlight: undefined,
            help: true
        })
    }

    private HideHelp() {
        if (this.state.focus !== undefined) {
            return
        }
        this.setState({
            ...this.state,
            help: false
        })
    }

    private setMenu(menu: State['menu']) {
        this.props.globalState.setMenu(menu)
        this.setState({
            ...this.state,
            days: menu.days.map(d => d.name),
            mealSizes: this.computeMealSizes(menu)
        })
    }

    private DayCol(day: Day): JSX.Element {
        return (
            <div className='Day'>
                <div className='Header' id='header1'>
                    <input onChange={(event) => {
                        day.name = event.target.value
                        this.setMenu(this.props.globalState.menu)
                    }}
                        defaultValue={day.name}
                    />
                </div>
                {
                    day.meals.map((meal, idx) =>
                        <div className="Meal" key={idx}>
                            <div className='MealHeader' key='MealName' id='header2'>
                                <input
                                    onChange={(event) => {
                                        meal.name = event.target.value
                                        this.setMenu(this.props.globalState.menu)
                                    }}
                                    defaultValue={meal.name}
                                />
                            </div>
                            <div className="Body" key='MealBody' style={{
                                minHeight: (this.state.mealSizes[idx] * 35 || 0) + 15
                            }} onClick={() => {
                                if (this.state.focus !== undefined) {
                                    return
                                }
                                this.Focus(day, meal)
                            }}>
                                {
                                    meal.dishes.map((dish, i) =>
                                        <DishItem
                                            key={dish.name}
                                            name={dish.name}
                                            amount={dish.amount}
                                            id={dish.name === this.state.hover
                                                ? 'highlight' :
                                                i % 2 === 0
                                                    ? 'odd' : 'even'
                                            }
                                            onMouseEnter={() => this.Highlight(dish)}
                                            onMouseLeave={() => this.Unhighlight()}
                                        />
                                    )
                                }
                            </div>
                        </div>
                    )
                }
            </div>
        )
    }

    private RenderFocus(): JSX.Element {
        const f = this.state.focus
        if (f === undefined) {
            return <></>
        }

        const day = f.day
        const meal = f.meal

        return (
            <dialog open>
                <h2 id='header'>
                    {meal.name} de {day.name}
                </h2>
                <div id="body">
                    {
                        meal.dishes.map((dish, i) =>
                            <DishPicker
                                key={`dish-${i}`}
                                recipes={this.props.globalState.dishes}
                                default={dish}
                                onChange={(newDish) => {
                                    new Optional(this.props.globalState.menu)
                                        .then(menu => menu.days.find(d => d.name === day.name))
                                        .elseLog(`Could not find day ${day.name}`)
                                        .then(day => day.meals.find(m => m.name === meal.name))
                                        .elseLog(`Could not find meal ${meal.name}`)
                                        .then(meal => meal.dishes[i] = newDish)
                                        .then(() => this.setMenu(this.props.globalState.menu))
                                }}
                                onRemove={() => {
                                    meal.dishes.splice(i, 1)
                                    this.forceUpdate()
                                }}
                            />
                        )
                    }
                    <button className='AddOne' onClick={() => {
                        meal.dishes.push(new Dish("", 1))
                        this.forceUpdate()
                    }}> + </button>
                </div>
                <div id='footer'>
                    <button onClick={() => {
                        this.setMenu(this.props.globalState.menu) // Trigger a cleanup
                        this.Unfocus()
                    }
                    }>Tancar</button>
                </div>
            </dialog>
        )

    }
    private RenderHelp(): JSX.Element {
        if (!this.state.help) {
            return <></>
        }

        return (
            <dialog open>
                <h2 id='header'>
                    Menú
                </h2>
                <div id="body">
                    <p>Aquesta pàgina et permet planificar els àpats de la setmana.</p>
                    <p>Pots clicar sobre qualsevol àpat per editar els seus continguts</p>
                    <p>Quan estigui llest, clica següent!</p>
                </div>
                <div id='footer'>
                    <button onClick={() => this.HideHelp()}>
                        D'acord
                    </button>
                </div>
            </dialog>
        )
    }

    private computeMealSizes(menu: Menu): number[] {
        return menu.days.map((day: Day) => {
            return day.meals.map(m => m.dishes.length)
        }).reduce((acc: number[], val: number[]): number[] => {
            return acc.map((v, i) => Math.max(v, val[i] || 0)).concat(val.slice(acc.length))
        })
    }

}

function DishItem(pp: { name: string, amount: number, id: string, onMouseEnter: () => void, onMouseLeave: () => void }) {
    return (
        <div
            className='DishItem'
            key={pp.name}
            onMouseEnter={pp.onMouseEnter}
            onMouseLeave={pp.onMouseLeave}
            id={pp.id}
        >
            <span id='amount' key='Amount'>{round2(pp.amount)}</span>
            <span id='name' key='Name'>{pp.name}</span>
        </div>
    )
}

async function saveMenu(backend: Backend, globalState: State): Promise<void> {
    backend.Menu().POST(globalState.menu)
}