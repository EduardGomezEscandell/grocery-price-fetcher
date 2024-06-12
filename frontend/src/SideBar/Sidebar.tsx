import React from 'react'
import './Sidebar.css'
import { useNavigate } from 'react-router-dom'

interface Props {
    onHelp: () => void
    onNavigate: () => void
}

export default function Sidebar(props: Props) {
    const n = useNavigate()
    const navigate = (path: string) =>{
        if (props.onNavigate) {
            props.onNavigate()
        }
        n(path)
    }

    return (
        <div className='sidebar'>
            <div id='header' onClick={() => navigate("/")}>
                <h1>La&nbsp;compra de&nbsp;l'Edu</h1>
            </div>
            <div id='body'>
                <p id="inactive">
                    Els meus productes
                </p>
                <p id="inactive">
                    Les meves receptes
                </p>
                <p onClick={() => navigate('/menu')}>
                    El meu menú
                </p>
                <p onClick={() => navigate('/pantry')}>
                    El meu rebost
                </p>
                <p onClick={() => navigate('/shopping-list')}>
                    La meva llista de la compra
                </p>
            </div>
            <div id='footer'>
                <p onClick={props.onHelp}>
                    Ajuda
                </p>
            </div>
        </div>
    )
}

