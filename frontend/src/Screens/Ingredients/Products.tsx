import React, { useState } from 'react'
import { useNavigate } from 'react-router-dom';
import TopBar from '../../TopBar/TopBar'
import Sidebar from '../../SideBar/Sidebar'
import Backend from '../../Backend/Backend';
import './Products.css'

interface Props {
    backend: Backend;
    sessionName: string;
}

export default function Products(props: Props) {
    const [sideBar, setSidebar] = useState(false)
    const [help, setHelp] = useState(false)
    const navigate = useNavigate()

    return (
        <div id='rootdiv'>
            <TopBar
                left={<button onClick={() => setSidebar(!sideBar)}> Opcions </button>}
                right={<></>}
                titleText="Els meus productes"
                logoOnClick={() => {
                    props.backend.ClearCache()
                    navigate('/')
                }}
                titleOnClick={() => setHelp(true)}
            />
            <div className='product-table-search'>
                Hello, world!
            </div>
            <section>
                {help && <HelpDialog onClose={() => setHelp(false)} />}
                {sideBar && <Sidebar onHelp={() => setHelp(true)} onNavigate={() => { props.backend.ClearCache() }} />}
            </section>
        </div>
    )
}

function HelpDialog(props: { onClose: () => void }): JSX.Element {
    return (
        <dialog open>
            <h2 id="header">Els meus productes</h2>
            <div id="body">
                <p>
                    Aquesta pàgina pàgina et permet veure i editar els teus productes, i crear-ne de nous.
                </p>
            </div>
            <div id="footer">
                <button onClick={props.onClose}>
                    D'acord
                </button>
            </div>
        </dialog>
    )
}