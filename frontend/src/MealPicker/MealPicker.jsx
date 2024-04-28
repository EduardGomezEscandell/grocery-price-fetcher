import React, { useState } from 'react'
import Select from 'react-select'


function MealPicker(pp) { 
    const options = pp.recipes.map(recipe => ({ value: recipe, label: recipe }))

    const [rec, setRec] = useState(undefined)
    const [amount, setAmount] = useState("")

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
                        onChange={selected => {
                            if (selected == null) {
                                return onMealChange("", 0)
                            }
                            onMealChange(selected.value, 0)
                        }}
                        defaultValue={{ value: "", label: "" }}
                        options={options}
                        style={{ width: "1000px" }}
                        isClearable
                    />
                </td>
                <td>
                    <input
                        type="number"
                        min="1"
                        value={amount}
                        onChange={s => onMealChange(rec, s.target.value)}
                        style={{ width: "50px" }}
                        contentEditable={rec !== undefined && rec !== ""}
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
