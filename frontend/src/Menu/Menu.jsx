import React, { useState } from 'react'
import MealPicker from '../MealPicker/MealPicker.jsx'

function Menu(pp) {
    const [dishes, setDishes] = useState(new Map())
    const [dishComp, setDishComp] = useState("")

    const updateMeal = (id, name, amount) => {
        const newDishes = dishes.set(id, {name: name, amount: amount})
        setDishes(newDishes)

        let m = new Map()
        Array
            .from(newDishes.values())
            .forEach(dish => {m[dish.name] = (m[dish.name] || 0) + dish.amount})
        
        const newComp = Object
            .keys(m)
            .filter(key => key !== undefined && key !== "")
            .filter(key => m[key] > 0)
            .sort()
            .map(key => <li key={key}>{key}: {m[key]}</li>)
        
        console.log(m)
        setDishComp(newComp)
    }

    const days = ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"]
    const eatings = ["Breakfast", "Lunch", "Snack", "Dinner"]

    const newRow = (eating) => {
        return (<tr>
            <td>{eating}</td>
            {days.map(day => {
                return (
                    <td>
                        <MealPicker
                            recipes={pp.recipes}
                            onChange={(name, amount) => updateMeal(
                                new MealID(day, eating, 0).toString(), name, amount)
                            }
                        />
                    </td>
                )
            })}
        </tr>
        )
    }

    return (
        <>
            <table>
                <tbody>
                    <tr>
                        <th>Meal</th>
                        {days.map(day => (<th>{day}</th>))}
                    </tr>
                    {eatings.map(eating => newRow(eating))}
                </tbody>
            </table>
            <div>
                {dishComp}
            </div>
        </>
    )
}

class MealID {
    constructor(day, eating, pos) {
        this.day = day
        this.eating = eating
        this.pos = pos
    }

    toString() {
        return `${this.day}-${this.eating}-${this.pos}`
    }
}

export default Menu
