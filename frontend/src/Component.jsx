import React, {useState} from 'react'

function Component() {
    let [recipes, setRecipes] = useState('Fetching...')

    fetch('/hello-world').then(response => {
        return response.text()
    }, error => {
        return error.message
    }).then(data => {
        console.log(data)
        setRecipes(data)
    })

    return (
        <div> HELLO! {recipes} </div>
    )
}

export default Component
