import React, { useEffect , useState } from 'react'
import Backend from '../../Backend/Backend.ts';
import { State, Menu } from '../../State/State.tsx';

interface Props {
    backend: Backend;
    globalState: State;
    onComplete: () => void
}

function FirstLoad(props: Props) {
    useEffect(() => {
        Promise.all([
        props.backend
            .Dishes()
            .GET()
            .then((d: string[]) => {props.globalState.dishes = d}),
        props.backend
            .Menu()
            .GET()
            .then((m: Menu[]) => m[0])
            .then((m: Menu) => props.globalState.menu = m)
        ]).finally(props.onComplete)
    })
    
    return (
        <p>Loading...</p>
    )

}

export default FirstLoad
