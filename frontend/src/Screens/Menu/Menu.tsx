import React from 'react'
import MealPicker from './MealPicker.tsx'
import Backend from '../../Backend/Backend.ts';
import Optional from '../../Optional/Optional.ts';
import { State, Day, Dish, Meal } from '../../State/State.tsx';
import TopBar from '../../TopBar/TopBar.tsx';
import './Menu.css'

interface Props {
    backend: Backend;
    globalState: State;
    onComplete: () => void
}

export default class MenuTable extends React.Component<Props> {
    constructor(props: Props) {
        super(props)
        this.days = props.globalState.menu.days.map(d => d.name)
        this.meals = Array.from(
            new Set<string>(
                props.globalState.menu.days.flatMap(d => d.meals.map(m => m.name))
            )
        )
        this.onComplete = props.onComplete
    }

    days: string[];
    meals: string[];
    onComplete: () => void

    render(): JSX.Element {
        return (
            <>
                <TopBar components={[
                    () => <p key='1' className='TopBar.Text'>Grocery Price Fetcher</p>,
                    () => <p key='2' className='TopBar.Text'>Menu</p>,
                    () => <button key='3'
                        className='TopBar.Button'
                        onClick={this.props.onComplete}
                    >Guardar i continuar</button>,
                ]}
                ></TopBar>
                <div className='Menu' key='menu-table'>
                    <table>
                        <tbody>
                            <tr className='Header'>
                                <th key='meal'>Meal</th>
                                {
                                    this.days.map(day => (<th key={day}>{day}</th>))
                                }
                            </tr>
                            {
                                this.meals.map(meal => (
                                    <tr className='Row' key={meal}>
                                        {this.RenderRow(meal)}
                                    </tr>
                                ))
                            }
                        </tbody>
                    </table >
                </div>
            </>
        )
    }

    private RenderRow(mealName: string): JSX.Element {
        return (
            <>
                <td className='MealName' key={mealName}>{mealName}</td>
                {
                    this.days
                        .map((dayName: string) => {
                            return new Optional(this.props.globalState.menu.days.find(d => d.name === dayName))
                                .then(day => new Optional(day.meals.find(m => m.name === mealName))
                                    .then(meal => this.RenderMeal(day, meal))
                                    .else(<td ></td>)
                                )
                                .else(<td></td>)
                        })
                }
            </>
        )
    }

    private RenderMeal(day: Day, meal: Meal): JSX.Element {
        return (
            <td className='Meal' key={day.name + meal.name}>
                <table>
                    <tbody>
                        {
                            meal.dishes.map((dish: Dish, i: number) => {
                                return this.RenderDish(day, meal, i, dish)
                            })
                        }
                    </tbody>
                </table>
            </td>
        )
    }

    private RenderDish(day: Day, meal: Meal, id: number, dish: Dish): JSX.Element {
        return (
            <tr key={id} >
                <td>
                    <MealPicker
                        recipes={this.props.globalState.dishes}
                        default={dish}
                        onChange={(newDish) => {
                            new Optional(this.props.globalState.menu)
                                .then(menu => menu.days.find(d => d.name === day.name))
                                .elseLog(`Could not find day ${day.name}`)
                                .then(day => day.meals.find(m => m.name === meal.name))
                                .elseLog(`Could not find meal ${meal.name}`)
                                .then(meal => meal.dishes[id] = newDish)
                                .then(() => this.props.globalState.setMenu(this.props.globalState.menu))
                        }}
                    />
                </td>
            </tr>
        )
    }
}
