import React from 'react'
import { useNavigate } from 'react-router-dom'
import './NotFound.css'

interface Props {}

export default function NotFound(props: Props) {
    const {} = props

    const navigate = useNavigate()
    return (
        <div className='not-found'>
            <h1>Aquesta adreça no existeix</h1>
            <div>
                <p>L'adreça que has introduït no existeix</p>
                <p>{window.location.href}</p>
                <p>Prem el botó per tornar a la pàgina principal</p>
                <button onClick={() => navigate("/")}>Tornar</button>
            </div>
        </div>
    )
}
