import React, { useEffect } from 'react'
import Backend from '../../Backend/Backend.ts';
import { State, Menu } from '../../State/State.tsx';

interface Props {
    backend: Backend;
    globalState: State;
    onComplete: () => void
}

export default function MenuLoad(props: Props) {
    useEffect(() => {
        Promise.all([
        props.backend
            .Dishes()
            .GET()
            .then((d: string[]) => {props.globalState.dishes = d}),
        props.backend
            .Menu()
            .GET()
            .then((m: Menu[]) => m[0] || new Menu())
            .then((m: Menu) => props.globalState.menu = m)
        ]).finally(props.onComplete)
    })
    
    return (
        <p>Loading...</p>
    )
}
