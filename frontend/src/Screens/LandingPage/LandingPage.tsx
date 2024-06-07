import React from 'react'
import Backend from '../../Backend/Backend.ts';
import { State, Menu } from '../../State/State.tsx';
import './LandingPage.css';

interface Props {
    backend: Backend;
    globalState: State;
    onGotoMenu: () => void
}

enum DialogState {
    None,
    CommingSoon,
    Help,
}

export default function LandingPage(props: Props) {
    const [commingSoon, setDialog] = React.useState(DialogState.None)

    const baseStyle: React.CSSProperties = {}
    if (commingSoon) {
        baseStyle.filter = 'blur(5px)'
    }

    return (
        <div className='LandingPage'>
            <div id="title" style={baseStyle}>
                <img src='/logo1024.png' alt='logo' className="Logo" />
                <h1>La&nbsp;compra de&nbsp;l'Edu</h1>
            </div>
            <div id="content" style={baseStyle}>
                <button  onClick={() => {
                    setDialog(DialogState.Help)
                }}>
                    Com funciona?
                </button>
                <button id="inactive" onClick={() => {
                    setDialog(DialogState.CommingSoon)
                }}>
                    Els meus productes
                </button>
                <button id="inactive" onClick={() => {
                    setDialog(DialogState.CommingSoon)
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
            {commingSoon === DialogState.CommingSoon && (
                <dialog open>
                    <h2 id='header'>No encara!</h2>
                    <div id="body"><p>Aquesta funcionalitat encara no està disponible</p></div>
                    <div id='footer'><button onClick={() => setDialog(DialogState.None)}>Entesos</button></div>
                </dialog>
            )}
            {commingSoon === DialogState.Help && (
                <dialog open>
                    <h2 id='header'>Com funciona?</h2>
                    <div id="body">
                        <p><b>La compra de l'Edu</b> t'ajuda a planificar la teva compra setmanal. Tingues en compte que està en fase experimental.</p>
                        <p>A <b>Els meus productes</b> pots afegir productes del supermercat que prefereixis.</p>
                        <p>A <b>Les meves receptes</b> pots afegir receptes utilitzant els teus productes com a ingredients.</p>
                        <p>A <b>La meva compra</b> pots crear un menú setmanal. A partir d'aquest menú, <i>La compra de l'Edu</i> calcularà
                        quant en necessites de cada ingredient i et preguntarà quant en tens de cada al teu rebost. Tot seguit,
                        et preparà la llista de la compra amb només allò que et falti.</p>
                    </div>
                    <div id='footer'><button onClick={() => setDialog(DialogState.None)}>Entesos</button></div>
                </dialog>
            )}
            <div id="footer" style={baseStyle}>
                <p>
                    La compra de l'Edu és un projecte de codi obert desenvolupat
                    per <a href='https://www.linkedin.com/in/eduard-gomez' target="_blank" rel="noreferrer">Eduard&nbsp;Gómez&nbsp;Escandell</a> i
                    disponible a <a href="https://www.github.com/EduardGomezEscandell/grocery-price-fetcher" target="_blank" rel="noreferrer">GitHub</a></p>
            </div>
        </div >
    )
}
