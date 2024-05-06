import React, { useState, useEffect, useCallback } from 'react'
import MealPicker from './MealPicker.jsx'

function Menu(pp) {
    const [loadingRecipes, setLoadingRecipes] = useState(true)
    const [loadingMenu, setLoadingMenu] = useState(true)

    const [recipes, setRecipes] = useState([""])
    const [menu, setMenu] = useState({ name: "loading", menu: []})

    const days = ["Dilluns", "Dimarts", "Dimecres", "Dijous", "Divendres", "Dissabte", "Diumenge"]
    const eatings = ["Esmorzar", "Dinar", "Berenar", "Sopar"]

    useEffect(() => {
        pp.backend.GetDishes()
            .then(response => response.json())
            .then(data => setRecipes(data))
            .finally(() => setLoadingRecipes(false))
    }, [pp.backend])

    useEffect(() => {
        pp.backend.GetMenu()
            .then(response => response.json())
            .then(data => setMenu(data[0]))
            .finally(() => setLoadingMenu(false))
    }, [pp.backend])

    const [dishes, setDishes] = useState(new Map())
    const [dishComp, setDishComp] = useState("")


    const updateDish = useCallback((id, name, amount) => {
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
    }, [dishes])

    useEffect(() => {
        console.log(menu)
        for (let day of menu.menu) {
            for (let meal of day.meals) {
                for (let i = 0; i < meal.dishes.length; i++) {
                    const id = new MealID(day.name, meal.name, i)
                    updateDish(id, meal.dishes[i].name, meal.dishes[i].amount)
                }
            }
        }
    }, [updateDish, menu])

    if (loadingRecipes || loadingMenu) {
        return <p>Loading...</p>
    }

    const newPlate = (day, meal, pos) => {
        const id = new MealID(day, meal, pos)
        return (
            <tr key={id.toString()} >
                <MealPicker
                    recipes={recipes}
                    default={dishes.get(id.toString())}
                    onChange={(name, amount) => updateDish(id, name, amount)}
                />
            </tr>
        )
    }

    const newMeal = (day, meal) => {
        return (
            <td>
            <table>
            <tbody>
            {[0,1,2,3].map(i => newPlate(day, meal, i))}
            </tbody>
            </table>
            </td>
        )
    }

    const newRow = (meal) => {
        return (<tr>
            <td>{meal}</td>
            {days.map(day => newMeal(day, meal))}
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
