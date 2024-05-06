import React, { useEffect, useState } from 'react'
import Ingredient from './Ingredient.jsx'

function Pantry(pp) {
    const [ingredients, setIngredients] = useState({})
    const [total, setTotal] = useState(0)
    const [loading, setLoading] = useState(true)

    useEffect(() => {
        let reqBody = {
            format: 'json',
            menu: []
        }

        Array
            .from(pp.recipes.values())
            .map(entry => {
                // Find or create day
                let iday = reqBody.menu.indexOf(entry.mealID.day)
                if (iday === -1) {
                    reqBody.menu.push({ name: entry.mealID.day, meals: [] })
                    iday = reqBody.menu.length - 1
                }
                return { Entry: entry, Day: reqBody.menu[iday] }
            }).map(x => {
                // Find or create meal
                let iMeal = x.Day.meals.indexOf(x.Entry.mealID.meal)
                if (iMeal === -1) {
                    x.Day.meals.push({ name: x.Entry.mealID.meal, Dishes: [] })
                    iMeal = x.Day.meals.length - 1
                }
                return { Entry: x.Entry, Meal: x.Day.meals[iMeal] }
            }).forEach(x => {
                // Append dish
                x.Meal.Dishes.push({ name: x.Entry.name, amount: x.Entry.amount })
            })

        pp.backend
            .PostMenu(reqBody)
            .then(response => response.json())
            .then(data => data.reduce((map, x) => map.set(x.product, new ShoppingItem(x)), new Map()))
            .then(ingredients => {
                setIngredients(ingredients)
                setTotal(computeTotal(ingredients))
                console.log(ingredients)
                console.log(computeTotal(ingredients))
            }).finally(() => setLoading(false))
    }, [pp.backend, pp.recipes])

    if (loading) {
        return <p>Loading...</p>
    }

    const updateHave = (name, have) => {
        console.log(`Updating ${name} to ${have}`)

        const newIngredients = ingredients.set(
            name,
            ingredients.get(name).setHave(have)
        )
        setIngredients(newIngredients)
        setTotal(computeTotal(newIngredients))
    }

    return (
        <>
            <h1>Pantry</h1>
            <table>
                <tbody>
                    <tr>
                        <th>Product</th>
                        <th>Have</th>
                        <th>Need</th>
                        <th>Deficit</th>
                        <th>Unit cost</th>
                        <th>Deficit cost</th>
                    </tr>
                    {
                        Array.from(ingredients.values())
                            .map(ingredient => (
                                <Ingredient
                                    key={ingredient.name}
                                    ingredient={ingredient}
                                    onChange={have => updateHave(ingredient.name, have)}
                                />
                            ))}
                    <tr name='pantry-total'>
                            <td colSpan="5"> <b>Total</b></td>
                            <td> <b>{total.toFixed(2)} € </b></td>
                    </tr>
                </tbody>
            </table>
        </>
    )
}

function computeTotal(ingredients) {
    return Array
        .from(ingredients.values())
        .map(ingredient => ingredient.unit_cost * (ingredient.need - ingredient.have))
        .filter(x => x > 0)
        .reduce((acc, x) => acc + x, 0)
}

class ShoppingItem {
    constructor(response) {
        this.name = response.product
        this.unit_cost = response.unit_cost
        this.need = parseFloat(response.amount) || 0
        this.have = 0
    }

    setHave(have) {
        this.have = have
        return this
    }
}


export default Pantry
