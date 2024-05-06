import React, { useEffect, useState } from 'react'
import Select from 'react-select'


function MealPicker(pp) { 
    const options = pp.recipes.map(recipe => ({ value: recipe, label: recipe }))

    const [rec, setRec] = useState("")
    const [amount, setAmount] = useState("")
    const [loading, setLoading] = useState(true)

    useEffect(() => {
        try {
            setRec(pp.default.name)
            setAmount(pp.default.amount)
            console.log(pp.default.name, pp.default.amount)
            setLoading(false)
        } catch {
            setRec("")
            setAmount("")
        }
    }, [pp.default])

    const onMealChange = (newRec, newAmount) => {
        switch (newRec) {
            case "":
                setRec(undefined)
                setAmount("")
                pp.onChange(undefined, 0)
                return
            case undefined:
                setRec(undefined)
                setAmount("")
                pp.onChange(undefined, 0)
                return
            default:
                setRec(newRec)
                setAmount(newAmount)
                pp.onChange(newRec, parseNum(newAmount))
                return
        }
    }
   
    return (
        <table>
            <tbody>
            <tr>
                <td>
                    <Select
                        styles={{ control: (base) => ({ ...base, width: '200px' }) }}
                        onChange={selected => {
                            if (selected == null) {
                                return onMealChange("", 0)
                            }
                            onMealChange(selected.value, 0)
                        }}
                        value={{ value: rec, label: rec }}
                        options={options}
                        isClearable
                    />
                </td>
                <td>
                    <input
                        type="number"
                        min="1"
                        value={amount}
                        onChange={s => onMealChange(rec, s.target.value)}
                        contentEditable={rec !== undefined && rec !== ""}
                        style={{ width: '50px', height: '25px' }}
                    />
                </td>
            </tr>
            </tbody>
        </table>
    )
}

function parseNum(num) {
    let n = parseInt(num)
    if (!isNaN(n)) {
        return n
    }

    n = parseFloat(num)
    if (!isNaN(n)) {
        return n
    }

    if (num === "") {
        return 0
    }

    if (num === undefined) {
        return 0
    }

    console.error(`Failed to parse number ${num}`)
    return 0
}


export default MealPicker
