import React, { useState } from 'react'
import Backend from '../../Backend/Backend';
import { Day, Meal, Dish, Menu } from '../../State/State';
import TopBar from '../../TopBar/TopBar';
import DishPicker from './DishPicker'
import './Menu.css'
import { round2 } from '../../Numbers/Numbers';
import SaveButton from '../../SaveButton/SaveButton';
import { useNavigate } from 'react-router-dom';
import Sidebar from '../../SideBar/Sidebar';

interface Props {
    backend: Backend;
    sessionName: string;
}

export default function RenderMenu(props: Props) {
    const [state, setState] = useState(new State())

    const tableStyle: React.CSSProperties = {}
    if (state.focus !== undefined || state.help) {
        tableStyle.filter = 'blur(5px)'
    }

    if (!state.loaded) {
        Promise.all([
            props.backend.Dishes().GET(),
            props.backend.Menu(props.sessionName).GET(),
        ]).then(([dishes, menu]) => setState(state.With({ dishes: dishes, loaded: true }).WithMenu(menu)))
    }

    const navigate = useNavigate()
    const [sidebar, setSidebar] = useState(false)

    return (<div id='rootdiv'>
        <TopBar
            left={<button className='save-button' id='idle'
                onClick={() => setSidebar(!sidebar)}
            >Opcions </button>
            }
            logoOnClick={() => saveMenu(props.backend, state.menu).then(() => navigate('/'))}
            titleText='El&nbsp;meu menú'
            right={<SaveButton
                key='save'

                baseTxt='Següent'

                onSave={() => saveMenu(props.backend, state.menu).then(() => navigate('/'))}
                onSaveTxt='Desant...'

                onAcceptTxt='Desat'
                onAccept={() => navigate('/pantry')}

                onRejectTxt='Error'

            />}
        />
        <section>
            <div className='Menu'>
                <table key='menu-table' style={tableStyle}>
                    <tbody>
                        <tr>
                            {
                                state.menu.days.map((_, i) =>
                                    <td key={`day-col-${i}`}>
                                        <DayColumn state={state} setState={setState} path={new Path(i)} />
                                    </td>
                                )
                            }

                        </tr>
                    </tbody>
                </table>
            </div>
            <FocusDialog state={state} setState={setState} key={(state.focus && state.focus.toString()) || 'none'} />
            <HelpDialog state={state} setState={setState} key={(state.help && 'T') || 'F'} />
            {sidebar && <Sidebar
                onHelp={() => {
                    setSidebar(false)
                    setState(state.WithHelp())
                }
                }
                onNavigate={() => saveMenu(props.backend, state.menu)}
            />}
        </section>
    </div>
    )
}

interface SubProps {
    state: State
    setState: (s: State) => void
}

interface DayColumnProps extends SubProps {
    path: Path
}

function DayColumn({ state, setState, path }: DayColumnProps): JSX.Element {
    const m = state.menu

    return (
        <div className='Day'>
            <div className='Header' id='header1'>
                <input onChange={(event) => {
                    path.Day(m).name = event.target.value
                    setState(state.WithMenu(m))
                }}
                    defaultValue={path.Day(m).name}
                />
            </div>
            {
                path.Day(m).meals.map((meal, idx) => {
                    const p = new Path(path.day, idx)

                    return (
                        <div className="Meal" key={idx}>
                            <div className='MealHeader' key='MealName' id='header2'>
                                <input
                                    onChange={(event) => {
                                        meal.name = event.target.value
                                        setState(state.WithMenu(m))
                                    }}
                                    defaultValue={meal.name}
                                />
                            </div>
                            <div className="Body" key='MealBody' style={{
                                minHeight: (state.mealSizes[idx] * 35 || 0) + 15
                            }} onClick={() => {
                                if (state.focus !== undefined) {
                                    return
                                }
                                setState(state.WithFocus(p))
                            }}>
                                {
                                    meal.dishes.map((dish, i) =>
                                        <DishItem
                                            key={dish.name}
                                            name={dish.name}
                                            amount={dish.amount}
                                            id={dish.name === state.highlight
                                                ? 'highlight' :
                                                i % 2 === 0 ? 'odd' : 'even'
                                            }
                                            onMouseEnter={() => setState(state.WithHighlight(dish.name))}
                                            onMouseLeave={() => setState(state.WithoutHighlight())}
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

function FocusDialog({ state, setState }: SubProps): JSX.Element {
    const f = state.focus
    if (f === undefined) {
        return <></>
    }

    const day = f.Day(state.menu)
    const meal = f.Meal(state.menu)

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
                            recipes={state.dishes}
                            default={dish}
                            onChange={(newDish) => {
                                meal.dishes[i] = newDish
                            }}
                            onRemove={() => {
                                meal.dishes.splice(i, 1)
                                setState(state.WithMenu(state.menu))
                            }}
                        />
                    )
                }
                <button className='AddOne' onClick={() => {
                    meal.dishes.push(new Dish("", 1))
                    setState(state.WithMenu(state.menu))
                }}> + </button>
            </div>
            <div id='footer'>
                <button onClick={() => {
                    meal.dishes = meal.dishes.filter(d => d.name !== "" && d.amount !== 0)
                    setState(state.WithMenu(state.menu).WithoutFocus())
                }
                }>Tancar</button>
            </div>
        </dialog>
    )

}
function HelpDialog({ state, setState }: SubProps): JSX.Element {
    if (!state.help) {
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
                <button onClick={() => setState(state.WithoutHelp())}>
                    D'acord
                </button>
            </div>
        </dialog>
    )
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
    backend.Menu(menu.name).PUT(menu)
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

    toString(): string {
        return `${this.day}/${this.meal}/${this.dish}`
    }
}

interface IState {
    // Data loaded
    loaded?: boolean

    // Menu data
    days?: string[]
    mealSizes?: number[]
    dishes?: string[]
    menu?: Menu

    // UI state
    focus?: Path
    help?: boolean
    highlight?: string
}

class State {
    // Data loaded
    loaded: boolean

    // Menu data
    days: string[]
    mealSizes: number[]
    dishes: string[]
    menu: Menu

    // UI state
    focus: Path | undefined
    help: boolean
    highlight: string | undefined

    constructor(argv: IState = {}) {
        this.loaded = argv.loaded || false
        this.focus = argv.focus || undefined
        this.help = argv.help || false
        this.highlight = argv.highlight || undefined
        this.days = argv.days || []
        this.mealSizes = argv.mealSizes || []
        this.menu = argv.menu || new Menu()
        this.dishes = argv.dishes || []
    }

    With(argv: IState): State {
        return new State({
            ...this,
            ...argv
        })
    }

    WithMenu(menu: Menu): State {
        return this.With({
            menu: menu,
            days: menu.days.map(d => d.name),
            mealSizes: State.computeMealSizes(menu),
        })
    }

    private static computeMealSizes(menu: Menu): number[] {
        return menu.days.map((day: Day) => {
            return day.meals.map(m => m.dishes.length)
        }).reduce((acc: number[], val: number[]): number[] => {
            return acc.map((v, i) => Math.max(v, val[i] || 0)).concat(val.slice(acc.length))
        }, [])
    }

    WithFocus(path: Path): State {
        return this.With({ focus: path, help: false, })
    }

    WithoutFocus(): State {
        return this.With({ focus: undefined })
    }

    WithHighlight(dish: string): State {
        if (this.focus !== undefined) { return this }
        if (this.help) { return this }
        return this.With({ highlight: dish })
    }

    WithoutHighlight(): State {
        if (this.focus !== undefined) { return this }
        if (this.help) { return this }
        return this.With({ highlight: undefined })
    }

    WithHelp() {
        if (this.focus !== undefined) { return this }
        return this.With({ help: true, highlight: undefined })
    }

    WithoutHelp() {
        if (this.focus !== undefined) { return this }
        return this.With({ help: false })
    }
}

