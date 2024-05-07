import React, { useState } from 'react'
import MealPicker from './MealPicker.tsx'
import Backend from '../../Backend/Backend.ts';
import Optional from '../../Optional/Optional.ts';
import { State, Menu, Dish, Meal } from '../../State/State.tsx';

interface Props {
    backend: Backend;
    state: State;
    onComplete: () => void
}

function RenderMenu(pp: Props) {
    return new MenuTable(pp.state.menu, pp.state.dishes)
        .Render(pp.onComplete)
}

export default RenderMenu

class DishEntry {
    constructor(id: number, dish: Dish, setter: React.Dispatch<React.SetStateAction<Dish>>) {
        this.id = id
        this.value = dish
        this.setter = setter
    }

    id: number;
    value: Dish;
    setter: React.Dispatch<React.SetStateAction<Dish>>;

    get(): Dish {
        return this.value
    }

    set(d: Dish) {
        this.setter(d)
    }
}

class MenuTable {
    constructor(menu: Menu, allDishes: string[]) {
        this.menu = menu
        this.allDishes = allDishes
        this.days = this.menu.days.map(d => d.name)
        this.meals = Array.from(
            new Set<string>(
                this.menu.days.flatMap(d => d.meals.map(m => m.name))
            )
        )
    }

    menu: Menu;
    allDishes: string[];
    days: string[];
    meals: string[];


    Render(onComplete: () => void): JSX.Element {
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
                    <button onClick={onComplete}>Finish</button>
                </div>
            </>
        )
    }

    RenderRow(mealName: string): JSX.Element {
        return (
            <tr>
                <td>{mealName}</td>
                {
                    this.days
                        .map((dayName: string) => new Optional(this.menu)
                            .then(menu => menu.days.find(d => d.name === dayName))
                            .then(day => day.meals.find(m => m.name === mealName))
                            .then(meal => this.RenderMeal(meal))
                            .else(<td></td>))
                }
            </tr>
        )
    }

    RenderMeal(meal: Meal): JSX.Element {
        return (
            <td>
                <table>
                    <tbody>
                        {
                            meal.dishes.map((d: Dish, i: number) => {
                                const [dish, setDish] = useState(d)
                                return this.RenderDish(new DishEntry(i, dish, setDish))
                            })
                        }
                    </tbody>
                </table>
            </td>
        )
    }

    RenderDish(dish: DishEntry): JSX.Element {
        return (
            <tr key={dish.id} >
                <MealPicker
                    recipes={this.allDishes}
                    default={dish.get()}
                    onChange={(newDish) => dish.set(newDish)}
                />
            </tr>
        )
    }
}
