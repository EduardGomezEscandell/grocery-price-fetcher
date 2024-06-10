import React from 'react'
import Backend from '../../Backend/Backend.ts';
import { Day, Meal, Dish, Menu } from '../../State/State.tsx';
import TopBar from '../../TopBar/TopBar.tsx';
import DishPicker from './DishPicker.tsx'
import './Menu.css'
import { round2 } from '../../Numbers/Numbers.ts';
import SaveButton from '../../SaveButton/SaveButton.tsx';

interface Props {
    backend: Backend;
    sessionName: string;
    onComplete: () => void
    onGotoHome: () => void
}

class Path {
    constructor(day: number, meal: number = 0, dish: number = 0) {
        this.day = day
        this.meal = meal
        this.dish = dish
    }

    day: number
    meal: number
    dish: number

    Day(m: Menu): Day {
        return m.days[this.day]
    }

    Meal(m: Menu): Meal {
        return this.Day(m).meals[this.meal]
    }

    Dish(m: Menu): Dish {
        return this.Meal(m).dishes[this.dish]
    }
}

export default class MenuTable extends React.Component<Props> {
    state: {
        // Data loaded
        loaded: boolean

        // Menu data
        days: string[],
        mealSizes: number[]
        dishes: string[]
        menu: Menu

        // UI state
        focus: Path | undefined
        help: boolean
        hover: string | undefined
    }

    constructor(props: Props) {
        super(props)
        this.state = {
            loaded: false,

            focus: undefined,
            hover: undefined,
            help: false,

            days: [],
            mealSizes: [],
            menu: new Menu(),
            dishes: []
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

        if (!this.state.loaded) {
            Promise.all([
                this.props.backend
                    .Dishes()
                    .GET(),
                this.props.backend.Menu()
                    .GET()
                    .then(menu => menu.find(m => m.name === this.props.sessionName) || new Menu())
            ])
                .then(([dishes, menu]) => {
                    this.setMenu(menu, {
                        dishes: dishes,
                        loaded: true
                    })
                })
        }

        return (
            <>
                <TopBar
                    left={<SaveButton
                        key='goback'

                        baseTxt='Tornar'

                        onSave={() => saveMenu(this.props.backend, this.state.menu)}
                        onSaveTxt='Desant...'

                        onAcceptTxt='Desat'
                        onAccept={this.props.onGotoHome}

                        onRejectTxt='Error'
                    />}
                    logoOnClick={() => saveMenu(this.props.backend, this.state.menu).then(this.props.onGotoHome)}
                    titleOnClick={() => this.DisplayHelp()}
                    titleText='El&nbsp;meu menú'
                    right={<SaveButton
                        key='save'

                        baseTxt='Següent'

                        onSave={() => saveMenu(this.props.backend, this.state.menu).then(this.props.onGotoHome)}
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
                                    this.state.menu.days.map((day, i) =>
                                        <td key={`day-col-${i}`}>
                                            {this.DayCol(new Path(i))}
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

    private Focus(path: Path) {
        this.setState({
            ...this.state,
            help: false,
            focus: path,
        })
    }

    private Highlight(dish: string) {
        if (this.state.focus !== undefined) {
            return
        }
        if (this.state.help) {
            return
        }
        this.setState({
            ...this.state,
            hover: dish
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

    private setMenu(menu: Menu, args = {}) {
        this.setState({
            ...this.state,
            ...args,
            menu: menu,
            days: menu.days.map(d => d.name),
            mealSizes: this.computeMealSizes(menu),
        })
    }

    private DayCol(p: Path): JSX.Element {
        const m = this.state.menu

        return (
            <div className='Day'>
                <div className='Header' id='header1'>
                    <input onChange={(event) => {
                        p.Day(m).name = event.target.value
                        this.setMenu(m)
                    }}
                        defaultValue={p.Day(m).name}
                    />
                </div>
                {
                    p.Day(m).meals.map((meal, idx) => {
                        const path = new Path(p.day, idx)

                        return (
                            <div className="Meal" key={idx}>
                                <div className='MealHeader' key='MealName' id='header2'>
                                    <input
                                        onChange={(event) => {
                                            meal.name = event.target.value
                                            this.setMenu(this.state.menu)
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
                                    this.Focus(path)
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
                                                onMouseEnter={() => this.Highlight(dish.name)}
                                                onMouseLeave={() => this.Unhighlight()}
                                            />
                                        )
                                    }
                                </div>
                            </div>
                        )
                    })
                }
            </div>
        )
    }

    private RenderFocus(): JSX.Element {
        const f = this.state.focus
        if (f === undefined) {
            return <></>
        }

        const day = f.Day(this.state.menu)
        const meal = f.Meal(this.state.menu)

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
                                recipes={this.state.dishes}
                                default={dish}
                                onChange={(newDish) => {
                                    meal.dishes[i] = newDish
                                }}
                                onRemove={() => {
                                    meal.dishes.splice(i, 1)
                                    this.setMenu(this.state.menu)
                                }}
                            />
                        )
                    }
                    <button className='AddOne' onClick={() => {
                        meal.dishes.push(new Dish("", 1))
                        this.setMenu(this.state.menu)
                    }}> + </button>
                </div>
                <div id='footer'>
                    <button onClick={() => {
                        meal.dishes = meal.dishes.filter(d => d.name !== "" && d.amount !== 0)
                        this.setMenu(this.state.menu, { focus: undefined })
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
        }, [])
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

async function saveMenu(backend: Backend, menu: Menu): Promise<void> {
    backend.Menu().PUT(menu)
}