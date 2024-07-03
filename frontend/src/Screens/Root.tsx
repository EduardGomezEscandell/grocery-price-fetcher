import React from 'react'
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
    return <LoginRouter
        sessionName="default"
        loginPage={r => <LoginPage onLogIn={creds => r.logIn(creds)} />}
        routes={new Map<string, ElementFactory>([
            ["/", r => <LandingPage backend={r.backend()} onLogout={() => r.logOut()} />],
            ["/products", r => <Products backend={r.backend()} sessionName={r.props.sessionName} />],
            ["/recipes", r => <Recipes backend={r.backend()} sessionName={r.props.sessionName} />],
            ["/menu", r => <RenderMenu backend={r.backend()} sessionName={r.props.sessionName} />],
            ["/pantry", r => <RenderPantry backend={r.backend()} sessionName={r.props.sessionName} />],
            ["/shopping-list", r => <ShoppingList backend={r.backend()} sessionName={r.props.sessionName} />],
        ])
        }
    />
}

type ElementFactory = (r: LoginRouter) => React.ReactNode;

interface routerProps {
    sessionName: string;
    routes: Map<string, ElementFactory>;
    loginPage: ElementFactory;
}

// This class wraps a router with the requirement of logging in before accessing the routes.
class LoginRouter extends React.Component<routerProps, { backend: Backend | undefined }> {
    constructor(props: routerProps) {
        super(props);
        this.state = { backend: undefined };
    }

    isLoggedIn(): boolean {
        return this.state.backend !== undefined;
    }

    logOut() {
        this.setState({ backend: undefined });
    }

    logIn(auth: string) {
        this.setState({ backend: new Backend(auth) });
    }

    backend(): Backend {
        if (this.state.backend === undefined) {
            throw new Error("Backend not initialized");
        }
        return this.state.backend;
    }

    render(): JSX.Element {
        return (
            <BrowserRouter>
                <Routes>
                    {
                        Array.from(this.props.routes.entries()).map((entry: [string, ElementFactory]) => {
                            return <Route
                                key={entry[0]}
                                path={entry[0]}
                                element={this.isLoggedIn()
                                    ? entry[1](this)
                                    : this.props.loginPage(this)
                                }
                            />
                        })
                    }
                    <Route path="*" element={<NotFound />} />
                </Routes>
            </BrowserRouter>
        );
    }
}
