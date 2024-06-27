import React, { useState } from 'react'
import TopBar from '../../TopBar/TopBar'
import Sidebar from '../../SideBar/Sidebar'
import Backend from '../../Backend/Backend';
import RecipeEditor from './RecipeEditor';
import { useNavigate } from 'react-router-dom';
import ComparableString from '../../ComparableString/ComparableString';
import { Recipe } from '../../Backend/endpoints/Recipe';
import { Dish } from '../../State/State';

interface Props {
    backend: Backend;
    sessionName: string;
}

interface SearchableDish {
    id: number;
    name: ComparableString;
}

export default function Recipes(props: Props) {
    const [sideBar, setSidebar] = useState(false)
    const [help, setHelp] = useState(false)

    const [recipes, setRecipes] = useState<SearchableDish[]>([])
    const [loaded, setLoaded] = useState(false)
    const [query, setQuery] = useState(new ComparableString(''))
    const [hidden, setHidden] = useState<string[]>([])

    const result = recipes
        .filter(r => !hidden.includes(r.name.displayName))
        .filter((r) => r.name.contains(query))

    if (!loaded) {
        props.backend.Dishes()
            .GET()
            .then((d) => d.map(r => {
                return { id: r.id, name: new ComparableString(r.name) } as SearchableDish
            }
            ))
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
            />
            <div className='search-table-search'>
                <input id={result.length === 0 ? 'error' : 'search'}
                    type='text'
                    placeholder='Cerca receptes...'
                    value={query.displayName}
                    onChange={(q) => setQuery(new ComparableString(q.target.value))}
                />
            </div>
            <section>
                <div className='search-table'>
                    <div id='body' key={query.compareName}>
                        {
                            loaded &&
                            <NewRecipe
                                onClick={() => {
                                    const name = (() => {
                                        if (query.displayName !== '' && !recipes.find(a => query.equals(a.name))) {
                                            return query.displayName
                                        }

                                        const name = `Nova recepta`
                                        if (!recipes.find(a => a.name.displayName === name)) {
                                            return name
                                        }

                                        for (let i = 1; ; i++) {
                                            const name = `Nova recepta ${i}`
                                            if (!recipes.find(a => a.name.displayName === name)) {
                                                return name
                                            }
                                        }
                                    })()

                                    props.backend
                                        .Recipe(props.sessionName, 0)
                                        .POST(new Recipe(0, name, []))
                                        .then((id: number) => {
                                            setRecipes([
                                                { id: id, name: new ComparableString(name) },
                                                ...recipes
                                            ])
                                        })
                                }}
                            />
                        }
                        {
                            result.map((r) => (
                                <RecipeEditor
                                    key={r.id}
                                    backend={props.backend}
                                    sessionName={props.sessionName}
                                    dish={new Dish(r.id, r.name.displayName, 0)}
                                    setHidden={() => setHidden([...hidden, r.name.displayName])}
                                    onRename={(newName: string) => {
                                        const idx = recipes.findIndex(a => a.id === r.id)
                                        recipes[idx].name = new ComparableString(newName)
                                        setRecipes(recipes)
                                    }}
                                />
                            ))
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

function NewRecipe(props: { onClick: () => void }): JSX.Element {
    return (
        <div className='search-table-row' key={'recipe-editor'}>
            <div className='title' onClick={props.onClick}>
                <span>Afegir recepta...</span>
            </div>
        </div>
    )
}

function HelpDialog(props: { onClose: () => void }): JSX.Element {
    return (
        <dialog open>
            <h2 id="header">Les meves receptes</h2>
            <div id="body">
                <p>
                    Aquesta pàgina pàgina et permet veure, editar, o eliminar qualsevol les teves receptes.
                </p>
                <p>
                    Per afegir una nova recepta, clica "Afegir recepta..."
                </p>
                <p>
                    Les receptes es creen a utilitzant els productes que tens a <i>Els meus productes</i> com
                    a ingredients.
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
