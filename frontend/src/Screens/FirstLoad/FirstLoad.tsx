import React, { useEffect , useState } from 'react'
import Backend from '../../Backend/Backend.ts';
import { State, Menu } from '../../State/State.tsx';

interface Props {
    backend: Backend;
    state: State;
    onComplete: () => void
}

function FirstLoad(props: Props) {
    const [dishesReady, setDishesReady] = useState(false)
    const [menuReady, setMenuReady] = useState(false)
    
    useEffect(() => {
        props.backend
            .GetDishes()
            .then((d: string[]) => {props.state.dishes = d})
            .finally(() => setDishesReady(true))
    })

    useEffect(() => {
        props.backend
            .GetMenu()
            .then((m: Menu[]) => m[0])
            .then((m: Menu) => props.state.menu = m)
            .finally(() => setMenuReady(true))
    })

    useEffect(() => {
        if (dishesReady && menuReady) {
            props.onComplete()
        }
    })
    
    return (
        <p>Loading...</p>
    )

}

export default FirstLoad
