import React from 'react'
import Backend from '../../Backend/Backend';
import './LandingPage.css';
import { Outlet, useNavigate } from 'react-router-dom';

interface Props {
    backend: Backend;
    sessionName: string;
}

export default function LandingPage(props: Props) {
    const [help, setHelp] = React.useState(false)

    const baseStyle: React.CSSProperties = {}
    if (help) {
        baseStyle.filter = 'blur(5px)'
    }

    const navigate = useNavigate()

    return (
        <div className='LandingPage'>
            <div id="title" style={baseStyle}>
                <img src='/logo1024.png' alt='logo' className="Logo" />
                <h1>La&nbsp;compra de&nbsp;l'Edu</h1>
            </div>
            <div id="content" style={baseStyle}>
                <div id="iconrow">
                    <button onClick={() => {
                        setHelp(true)
                    }}>
                        Com funciona?
                    </button>
                    <button onClick={() => { navigate("/products") }}>
                        Els meus productes
                    </button>
                    <button onClick={() => { navigate('/recipes') }}>
                        Les meves receptes
                    </button>
                </div>
                <div id="iconrow">
                    <button onClick={() => navigate('/menu')}>
                        El meu menú
                    </button>
                    <button onClick={() => navigate('/pantry')}>
                        El meu rebost
                    </button>
                    <button onClick={() => navigate('/shopping-list')}>
                        La meva llista de la compra
                    </button>
                </div>
                <Outlet />
            </div>
            {help && <HelpDialog onClose={() => setHelp(false)} />}
            <div id="footer" style={baseStyle}>
                <p>
                    La compra de l'Edu és un projecte de codi obert desenvolupat
                    per <a href='https://www.linkedin.com/in/eduard-gomez' target="_blank" rel="noreferrer">Eduard&nbsp;Gómez&nbsp;Escandell</a> i
                    disponible a <a href="https://www.github.com/EduardGomezEscandell/grocery-price-fetcher" target="_blank" rel="noreferrer">GitHub</a></p>
            </div>
        </div >
    )
}

function HelpDialog(props: { onClose: () => void }): JSX.Element {
    return (
        <dialog open>
            <h2 id='header'>Com funciona?</h2>
            <div id="body">
                <p>
                    <b>La compra de l'Edu</b> t'ajuda a planificar la teva compra setmanal.
                    Tingues en compte que està en fase experimental.
                    Dins de cada pàgina, pots obtindre més ajuda clicant el títol de la pàgina.
                </p>
                <p>A <b>Els meus productes</b> pots afegir productes del supermercat que prefereixis.</p>
                <p>A <b>Les meves receptes</b> pots afegir receptes utilitzant els teus productes com a ingredients.</p>
                <p>A <b>La meva compra</b> pots crear un menú setmanal. A partir d'aquest menú, <i>La compra de l'Edu</i> calcularà
                    quant en necessites de cada ingredient i et preguntarà quant en tens de cada al teu rebost. Tot seguit,
                    et preparà la llista de la compra només amb allò que et faci falta.</p>
            </div>
            <div id='footer'><button onClick={props.onClose}>D'acord</button></div>
        </dialog>
    )
}