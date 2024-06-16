import React, { useState } from 'react'
import TopBar from '../../TopBar/TopBar'
import Sidebar from '../../SideBar/Sidebar'
import Backend from '../../Backend/Backend';
import RecipeEditor from './RecipeEditor';
import './Recipes.css'
import { useNavigate } from 'react-router-dom';

interface Props {
    backend: Backend;
    sessionName: string;
}

export default function Recipes(props: Props) {
    const [sideBar, setSidebar] = useState(false)
    const [help, setHelp] = useState(false)

    const [recipes, setRecipes] = useState<comparableString[]>([])
    const [loaded, setLoaded] = useState(false)
    const [query, setQuery] = useState(new comparableString(''))
    const [hidden, setHidden] = useState<string[]>([])

    const result = recipes
        .filter(r => !hidden.includes(r.displayName))
        .filter((r) => r.contains(query))

    if (!loaded) {
        props.backend.Dishes()
            .GET()
            .then((d) => d.map(r => new comparableString(r)))
            .then(setRecipes)
            .then(() => setLoaded(true))
    }

    const navigate = useNavigate()
    
    return (
        <div id='rootdiv'>
            <TopBar
                left={<button onClick={() => setSidebar(!sideBar)}> Opcions </button>}
                right={<></>}
                titleText="Les meves receptes"
                logoOnClick={() => {
                    props.backend.ClearCache()
                    navigate('/')
                }}
                titleOnClick={() => setHelp(true)}
            />
            <div className='recipe-table-search'>
                <input id={result.length === 0 ? 'error' : 'search'}
                    type='text'
                    placeholder='Cerca receptes...'
                    value={query.displayName}
                    onChange={(q) => setQuery(new comparableString(q.target.value))}
                />
            </div>
            <section>
                <div className='recipe-table'>
                    <div id='body' key={query.compareName}>
                        {
                            result.map((r) => {
                                return hidden.includes(r.displayName) ? null : (
                                    <RecipeEditor
                                        key={r.displayName}
                                        backend={props.backend}
                                        sessionName={props.sessionName}
                                        dish={r.displayName}
                                        setHidden={() => setHidden([...hidden, r.displayName])}
                                    />
                                )
                            })
                        }
                        {
                            result.length === 0 &&
                            <div id='error'>
                                No hi ha resultats
                            </div>
                        }
                        <p></p>
                    </div>
                </div>
                {help && <HelpDialog onClose={() => setHelp(false)} />}
                {sideBar && <Sidebar onHelp={() => setHelp(true)} onNavigate={() => { props.backend.ClearCache() }} />}
            </section>
        </div>
    )
}



function HelpDialog(props: { onClose: () => void }): JSX.Element {
    return (
        <dialog open>
            <h2 id="header">Les meves receptes</h2>
            <div id="body">
                <p>
                    Aquesta pàgina pàgina et permet veure les teves receptes, i crear-ne de noves.
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

class comparableString {
    displayName: string
    compareName: string

    constructor(displayName: string) {
        this.displayName = displayName
        this.compareName = this.localeFold(displayName)
    }

    private localeFold(s: string): string {
        return s.normalize("NFKD")               // Decompose unicode characters
            .replace(/[\u0300-\u036f]/g, "") // Remove accents
            .toLowerCase()
    }

    contains(other: comparableString): boolean {
        return this.compareName.includes(other.compareName)
    }
}

