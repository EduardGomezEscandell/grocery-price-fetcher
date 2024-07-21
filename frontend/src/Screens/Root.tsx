import React, { useEffect, useState } from 'react'
import { CookiesProvider, useCookies } from 'react-cookie';
import Backend from "../Backend/Backend";
import LandingPage from "./LandingPage/LandingPage";
import RenderMenu from "./Menu/Menu";
import RenderPantry from "./Pantry/Pantry";
import ShoppingList from './ShoppingList/ShoppingList';
import { BrowserRouter, Route, Routes } from "react-router-dom";
import NotFound from './404/NotFound.tsx/NotFound';
import Recipes from './Recipes/Recipes';
import Products from './Products/Products';
import LoginPage from './LandingPage/LoginPage';

export default function Root(): JSX.Element {
    const sessionName = "default"

    return <LoginRouter
        loginPage={r => <LoginPage logIn={creds => r.logIn(creds)} />}
        routes={new Map<string, ElementFactory>([
            ["/", r => <LandingPage backend={r.backend} logOut={() => r.logOut()} />],
            ["/products", r => <Products backend={r.backend} sessionName={sessionName} />],
            ["/recipes", r => <Recipes backend={r.backend} sessionName={sessionName} />],
            ["/menu", r => <RenderMenu backend={r.backend} sessionName={sessionName} />],
            ["/pantry", r => <RenderPantry backend={r.backend} sessionName={sessionName} />],
            ["/shopping-list", r => <ShoppingList backend={r.backend} sessionName={sessionName} />],
        ])
        }
    />
}

type ElementFactory = (r: Auth) => React.ReactNode;

interface routerProps {
    routes: Map<string, ElementFactory>;
    loginPage: ElementFactory;
}

// This class wraps a router with the requirement of logging in before accessing the routes.
function LoginRouter(props: routerProps): JSX.Element {
    const [refreshOnAuth, _f] = useState(0)
    const forceRefresh = () => _f(refreshOnAuth + 1)

    const cookie = new Cookie('GROCERY_PRICE_FETCHER_AUTH', forceRefresh)
    const auth = new Auth(cookie)

    const [loading, setLoading] = useState(true)
    if (loading) {
        auth.validate().then(() => setLoading(false))
    }

    return (
        <CookiesProvider>
            <BrowserRouter>
                <Routes key={refreshOnAuth.toString()}>
                    {
                        Array.from(props.routes.entries()).map((entry: [string, ElementFactory]) => {
                            return <Route
                                key={entry[0]}
                                path={entry[0]}
                                element={
                                    loading
                                        ? <>Carregant...</> 
                                        : auth.loggedIn()
                                            ?   entry[1](auth)
                                            : props.loginPage(auth)
                                }
                            />
                        })
                    }
                    <Route path="*" element={<NotFound />} />
                </Routes>
            </BrowserRouter>
        </CookiesProvider>
    );
}


class Auth {
    backend: Backend;
    private cookie: Cookie;

    constructor(cookie: Cookie) {
        this.cookie = cookie;
        this.backend = new Backend();
    }

    loggedIn(): boolean {
        return this.cookie.Get() !== undefined
    }

    async validate() {
        return this.backend.AuthRefresh()
            .POST()
            .then(v => this.cookie.Set(v))
            .catch(() => this.cookie.Remove())
    }

    async logIn(code: string): Promise<void> {
        return this.backend.AuthLogin()
            .POST(code)
            .then((auth: string) => this.cookie.Set(auth))
    }

    async logOut(): Promise<void> {
        return this.backend.AuthLogout()
            .POST()
            .then(() => this.cookie.Remove())
    }
}

class Cookie {
    Set: (v: string) => void
    Remove: () => void
    Get: () => string | undefined

    constructor(name: string, onChange: () => void) {
        const [cookies, setCookie, removeCookie] = useCookies([name]);

        var cached = cookies[name]
        const _onChange = (v?: string) => {
            if (cached === v) {
                return
            }
            cached = v
            onChange()
        }

        this.Set = (v: string) => {
            setCookie(name, v, { sameSite: 'strict', secure: true })
            _onChange(v)
        }

        this.Remove = () => {
            removeCookie(name);
            _onChange(undefined)
        }

        this.Get = () => {
            const value = cookies[name]
            _onChange(value)
            return value
        }
    }
}