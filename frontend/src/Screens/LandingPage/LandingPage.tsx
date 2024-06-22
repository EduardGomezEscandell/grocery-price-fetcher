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
                </p>
                <p>
                    Pots afegir productes del teu supermercat preferit, i crear receptes amb ells.
                </p>
                <p>
                    A partir de les receptes, pots crear-te un menú setmanal, i a partir d'aquest menú,
                    es calcularà la llista de la compra, tot tenint en compte el menjar que ja tens a casa.
                </p>
                <p>
                    Per informació més detallada, totes les pàgines tenen un botó d'ajuda al menú d'opcions.
                </p>
            </div>
            <div id='footer'><button onClick={props.onClose}>D'acord</button></div>
        </dialog>
    )
}