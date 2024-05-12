import React from 'react'
import MealPicker from './MealPicker.tsx'
import Backend from '../../Backend/Backend.ts';
import Optional from '../../Optional/Optional.ts';
import { State, Day, Dish, Meal } from '../../State/State.tsx';

interface Props {
    backend: Backend;
    globalState: State;
    onComplete: () => void
}

export default class MenuTable extends React.Component<Props> {
    constructor(props: Props) {
        super(props)
        this.globalState = props.globalState
        this.days = this.globalState.menu.days.map(d => d.name)
        this.meals = Array.from(
            new Set<string>(
                this.globalState.menu.days.flatMap(d => d.meals.map(m => m.name))
            )
        )
        this.onComplete = props.onComplete
    }

    globalState: State;
    days: string[];
    meals: string[];
    onComplete: () => void

    render(): JSX.Element {
        return (
            <>
                <table>
                    <tbody>
                        <tr>
                            <th>
                                Meal
                            </th>
                            {
                                this.days.map(day => (<th>{day}</th>))
                            }
                        </tr>
                        {
                            this.meals.map(meal => this.RenderRow(meal))
                        }
                    </tbody>
                </table>
                <div>
                    <button onClick={this.onComplete.bind(this)}>Finish</button>
                </div>
            </>
        )
    }

    private RenderRow(mealName: string): JSX.Element {
        return (
            <tr>
                <td>{mealName}</td>
                {
                    this.days
                        .map((dayName: string) => {
                            return new Optional(this.globalState.menu.days.find(d => d.name === dayName))
                                .then(day => new Optional(day.meals.find(m => m.name === mealName))
                                    .then(meal => this.RenderMeal(day, meal))
                                    .else(<td></td>)
                                )
                                .else(<td></td>)
                        })
                }
            </tr>
        )
    }

    private RenderMeal(day: Day, meal: Meal): JSX.Element {
        return (
            <td>
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
                <MealPicker
                    recipes={this.globalState.dishes}
                    default={dish}
                    onChange={(newDish) => {
                        new Optional(this.globalState.menu)
                            .then(menu => menu.days.find(d => d.name === day.name))
                            .elseLog(`Could not find day ${day.name}`)
                            .then(day => day.meals.find(m => m.name === meal.name))
                            .elseLog(`Could not find meal ${meal.name}`)
                            .then(meal => meal.dishes[id] = newDish)
                            .then(() => this.globalState.setMenu(this.globalState.menu))
                    }}
                />
            </tr>
        )
    }
}
