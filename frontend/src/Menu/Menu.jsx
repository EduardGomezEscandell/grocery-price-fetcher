import React, { useState, useEffect } from 'react'
import MealPicker from './MealPicker.jsx'

function Menu(pp) {
    const [loading, setLoading] = useState(true)
    const [recipes, setRecipes] = useState([""])

    useEffect(() => {
        pp.backend.fetchDishes()
            .then(response => response.json())
            .then(data => setRecipes(data))
            .finally(() => setLoading(false))
    }, [pp.backend])

    const [dishes, setDishes] = useState(new Map())
    const [dishComp, setDishComp] = useState("")

    if (loading) {
        return <p>Loading...</p>
    }

    const updateDish = (id, name, amount) => {
        const newDishes = (() => {
            if (amount === 0) {
                dishes.delete(id.toString())
                return dishes
            }
            return dishes.set(id.toString(), { mealID: id, name: name, amount: amount }) 
        })()
        setDishes(newDishes)

        let m = new Map()
        Array
            .from(newDishes.values())
            .forEach(dish => { m[dish.name] = (m[dish.name] || 0) + dish.amount })

        const newComp = Object
            .keys(m)
            .filter(key => key !== undefined && key !== "")
            .filter(key => m[key] > 0)
            .sort()
            .map(key => <li key={key}>{key}: {m[key]}</li>)

        setDishComp(newComp)
    }

    const days = ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"]
    const eatings = ["Breakfast", "Lunch", "Snack", "Dinner"]

    const newRow = (meal) => {
        return (<tr>
            <td>{meal}</td>
            {days.map(day => {
                const id = new MealID(day, meal, 0)
                return (
                    <td key={id.toString()} >
                        <MealPicker
                            recipes={recipes}
                            onChange={(name, amount) => updateDish(id, name, amount)}
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
            <div>
            <button onClick={() => pp.onComplete(dishes)}>Finish</button>
            </div>
        </>
    )
}

class MealID {
    constructor(day, eating, pos) {
        this.day = day
        this.meal = eating
        this.pos = pos
    }

    toString() {
        return `${this.day}-${this.meal}-${this.pos}`
    }
}

export default Menu
