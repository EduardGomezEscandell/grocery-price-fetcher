import React from 'react'
import './Sidebar.css'
import { useLocation, useNavigate } from 'react-router-dom'

interface Props {
    onHelp: () => void
    onNavigate: () => void
}

export default function Sidebar(props: Props) {
    const n = useNavigate()
    const navigate = (path: string) => {
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
                <Goto path={'/products'}>Els meus productes</Goto>
                <Goto path={'/recipes'}>Les meves receptes</Goto>
                <Goto path={'/menu'}>El meu menú</Goto>
                <Goto path={'/pantry'}>El meu rebost</Goto>
                <Goto path={'/shopping-list'}>La meva llista de la compra</Goto>
            </div>
            <div id='footer'>
                <p onClick={props.onHelp}>
                    Ajuda
                </p>
            </div>
        </div>
    )
}

function Goto(props: {
    path?: string,
    children: string,
    disabled?: boolean,
    onNavigate?: () => void
}) {

    const location = useLocation()
    const navigate = useNavigate()

    const onClick = () => { 
        if (props.disabled) {
            return
        }

        if (!props.path) {
            return
        }

        if (props.path === location.pathname) {
            return
        }

        if (props.onNavigate) {
            props.onNavigate()
        }

        navigate(props.path)
    }
    
    
    let id = props.disabled ? 'disabled' : 'enabled'
    if (props.path === location.pathname) {
        id = 'current'
    }

    return (
        <p onClick={onClick} id={id}>
            {props.children}
        </p>
    )
}

