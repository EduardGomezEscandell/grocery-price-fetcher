import React from 'react'
import Backend from '../../Backend/Backend.ts';
import { State, Menu } from '../../State/State.tsx';
import './LandingPage.css';

interface Props {
    backend: Backend;
    globalState: State;
    onGotoMenu: () => void
}

export default function LandingPage(props: Props) {
    const [commingSoon, setCommingSoon] = React.useState(false)

    const tableStyle: React.CSSProperties = {}
    if (commingSoon) {
        tableStyle.filter = 'blur(5px)'
    }

    return (
        <div className='LandingPage'>
            <div id="title" style={tableStyle}>
                <img src='/logo1024.png' alt='logo' className="Logo" />
                <h1>El&nbsp;rebost</h1>
            </div>
            <div id="content" style={tableStyle}>
                <button onClick={() => {
                    setCommingSoon(true)
                }}>
                    Els meus productes
                </button>
                <button onClick={() => {
                    setCommingSoon(true)
                }}>
                    Les meves receptes
                </button>
                <button onClick={() => {
                    Promise.all([
                        props.backend
                            .Dishes()
                            .GET()
                            .then((d: string[]) => { props.globalState.dishes = d }),
                        props.backend
                            .Menu()
                            .GET()
                            .then((m: Menu[]) => m[0])
                            .then((m: Menu) => props.globalState.menu = m)
                    ]).finally(props.onGotoMenu)
                }}>
                    La meva compra
                </button>
            </div>
            {commingSoon && (
                <dialog open>
                    <h2 id='header'>No encara!</h2>
                    <div id="body"><p>Aquesta funcionalitat encara no està disponible</p></div>
                    <div id='footer'><button onClick={() => setCommingSoon(false)}>Entesos</button></div>
                </dialog>
            )}
            <div id="footer" style={tableStyle}>
                <p>
                    El rebost és un projecte de codi obert desenvolupat
                    per <a href='https://www.linkedin.com/in/eduard-gomez' target="_blank" rel="noreferrer">Eduard&nbsp;Gómez&nbsp;Escandell</a> i
                    disponible a <a href="https://www.github.com/EduardGomezEscandell/grocery-price-fetcher" target="_blank" rel="noreferrer">GitHub</a></p>
            </div>
        </div >
    )
}
