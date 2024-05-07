import { useState, useEffect } from 'react'
import Backend from '../../Backend/Backend'
import { State } from '../../State/State'

interface Props {
    backend: Backend
    state: State
    onComplete: () => void
}

function PantryLoad(pp: Props) {
    useEffect(() => {
        pp.backend
            .PostMenu(pp.state.menu)
            .then(sh => pp.state.shoppingList = sh)
            .finally(() => pp.onComplete())
    })
   
    return (
        <p>Loading...</p>
    )

}

export default PantryLoad
