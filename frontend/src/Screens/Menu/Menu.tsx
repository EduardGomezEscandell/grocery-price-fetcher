import React from 'react'
import Backend from '../../Backend/Backend.ts';
import Optional from '../../Optional/Optional.ts';
import { State, Day, Meal, Dish } from '../../State/State.tsx';
import TopBar from '../../TopBar/TopBar.tsx';
import DishPicker from './DishPicker.tsx'
import './Menu.css'
import { round2 } from '../../Numbers/Numbers.ts';

interface Props {
    backend: Backend;
    globalState: State;
    onComplete: () => void
}

class MealMetadata {
    name: string;
    size: number;
}

export default class MenuTable extends React.Component<Props> {
    state: {
        days: string[],
        meals: MealMetadata[]
        focus: { day: Day, meal: Meal } | undefined
        hover: string | undefined
    }

    constructor(props: Props) {
        super(props)
        this.state = {
            focus: undefined,
            hover: undefined,
            days: props.globalState.menu.days.map(d => d.name),
            meals: props.globalState.menu.days.map((day: Day, i: number) => {
                return day.meals.map(meal => {
                    return {
                        name: meal.name,
                        size: meal.dishes.length,
                    }
                })
            }).reduce((acc: MealMetadata[], val: MealMetadata[]) => {
                val.forEach((m: MealMetadata) => {
                    var idx = acc.findIndex(x => x.name === m.name)
                    if (idx === -1) {
                        acc.push(m)
                        return
                    }
                    acc[idx].size = Math.max(acc[idx].size, m.size)
                })
                return acc
            }, [])
        }

    }


    get days(): string[] {
        return this.state.days
    }

    get meals(): MealMetadata[] {
        return this.state.meals
    }

    render(): JSX.Element {
        const tableStyle: React.CSSProperties = {}
        if (this.state.focus !== undefined) {
            tableStyle.filter = 'blur(5px)'
        }

        return (
            <>
                <TopBar
                    left={<p key='2' className='Text'>Menu</p>}
                    right={<button key='3' className='Button' onClick={this.props.onComplete}>Guardar i continuar</button>}
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
                </div>
            </>
        )
    }

    private DayCol(day: Day): JSX.Element {
        return (
            <div className='Day'>
                <div className='Header'>
                    {day.name}
                </div>
                {
                    day.meals.map((meal) =>
                        <div className="Meal" key={meal.name} onClick={() => {
                            if (this.state.focus !== undefined) {
                                return
                            }

                            this.setState({
                                ...this.state,
                                hover: undefined,
                                focus: {
                                    day: day,
                                    meal: meal
                                }
                            })
                        }}>
                            <div className='MealHeader' key='MealName'>
                                {meal.name}
                            </div>
                            <div className="Body" key='MealBody' style={{
                                minHeight: new Optional(this.meals.find(m => m.name === meal.name))
                                    .then(m => m.size * 50)
                                    .then(s => s.toString() + 'px')
                                    .else('0px')
                            }}>
                                {
                                    meal.dishes.map((dish, i) =>
                                        <DishItem
                                            name={dish.name}
                                            amount={dish.amount}
                                            id={dish.name === this.state.hover 
                                                ? 'mouseover' : 
                                                i % 2 === 0 
                                                ? 'odd' : 'even'
                                            }
                                            onMouseEnter={() => {
                                                if (this.state.focus !== undefined) {
                                                    return
                                                }
                                                this.setState({
                                                    ...this.state,
                                                    hover: dish.name
                                                })
                                            }}
                                            onMouseLeave={() => {
                                                if (this.state.focus !== undefined) {
                                                    return
                                                }
                                                this.setState({
                                                    ...this.state,
                                                    hover: undefined
                                                })
                                            }}
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
            <div className='Focus'>
                <dialog open>
                    <div className='FocusHeader'>
                        {meal.name} de {day.name}
                    </div>
                    <div className="FocusBody" style={{
                        minHeight: new Optional(this.meals.find(m => m.name === meal.name))
                            .then(m => m.size * 30)
                            .then(s => s.toString() + 'px')
                            .else('0px')
                    }}>
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
                                            .then(() => this.props.globalState.setMenu(this.props.globalState.menu))
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
                    <div className='FocusFooter'>
                        <button className='Button' onClick={() => {
                            this.props.globalState.setMenu(this.props.globalState.menu) // Trigger a cleanup
                            this.setState({
                                ...this.state,
                                focus: undefined
                            })
                        }
                        }>Tancar</button>
                    </div>
                </dialog>
            </div>
        )
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
            <span className='Amount' key='Amount'>{round2(pp.amount)}</span>
            <span className='Name' key='Name'>{pp.name}</span>
        </div>
    )
}