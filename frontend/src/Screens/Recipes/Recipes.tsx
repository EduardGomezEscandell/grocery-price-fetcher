import React, { useEffect, useState } from 'react'
import TopBar from '../../TopBar/TopBar'
import Sidebar from '../../SideBar/Sidebar'
import Backend from '../../Backend/Backend';
import './Recipes.css';
import { c } from 'vite/dist/node/types.d-aGj9QkWt';
import { Dish } from '../../State/State';

interface Props {
    backend: Backend;
    sessionName: string;
}

enum Focus {
    NONE,
    SIDEBAR,
    HELP,
}

export default function Recipes(props: Props) {
    const { } = props

    const [recipes, setRecipes] = useState<comparableString[]>([])
    const [loaded, setLoaded] = useState(false)

    const [focus, setFocus] = useState(Focus.NONE)
    const [query, setQuery] = useState(new comparableString(''))

    const result = recipes.filter((r) => r.contains(query))

    if (!loaded) {
        props.backend.Dishes()
            .GET()
            .then((d) => d.map(r => new comparableString(r)))
            .then(setRecipes)
            .then(() => setLoaded(true))
    }

    return (
        <div id='rootdiv'>
            <TopBar
                left={<button onClick={() => setFocus(focus === Focus.SIDEBAR ? Focus.NONE : Focus.SIDEBAR)}> Opcions </button>}
                right={<></>}
                titleText="Les meves receptes"
            />
            <section>
                <div className='Recipes'>
                    <div id='header'>
                        <input id={result.length === 0 ? 'error' : 'search'}
                            type='text'
                            placeholder='Cerca receptes...'
                            value={query.displayName}
                            onChange={(q) => setQuery(new comparableString(q.target.value))}
                        />
                    </div>
                    <div id='body' key={query.compareName}>
                        {
                            result.map((r, i) => {
                                return (
                                    <div
                                        className='recipe'
                                        key={query.compareName + r.compareName}
                                        id='even'>
                                        {r.displayName}
                                    </div>
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
                {renderFocus(focus, setFocus)}
            </section>
        </div>
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


function renderFocus(focus: Focus, setFocus: (focus: Focus) => void): JSX.Element {
    switch (focus) {
        default:
            return <></>
        case Focus.HELP:
            return (
                <HelpDialog onClose={() => setFocus(Focus.NONE)} />
            )
        case Focus.SIDEBAR:
            return (
                <Sidebar
                    onHelp={() => setFocus(Focus.HELP)}
                    onNavigate={() => setFocus(Focus.NONE)}
                />
            )
    }
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
