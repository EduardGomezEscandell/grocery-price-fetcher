import React from 'react'
import Backend from "../Backend/Backend.tsx";
import LandingPage from "./LandingPage/LandingPage.tsx";
import RenderMenu from "./Menu/Menu.tsx";
import Pantry from "./Pantry/Pantry.tsx";
import ShoppingList from './ShoppingList/ShoppingList.tsx';

interface Props {
    backend: Backend;
    sessionName: string;
}

export default class StateMachine extends React.Component<Props> {
    state: { screen: Screen }

    constructor(props: Props) {
        super(props)
        const baseScreen = new Screen({
            ...props,
            setScreen: (s: Screen) => this.setState({ screen: s })
        })

        this.state = {
            screen: new HomeScreen(baseScreen)
        }
    }

    render(): JSX.Element {
        return (
            <div style={{
                display: 'flex',
                flexDirection: 'column',
                height: '100vh'
            }}>
                {this.state.screen.render()}
            </div>
        )
    }
}

interface ScreenProps extends Props {
    setScreen: (s: Screen) => void;
}

class Screen extends React.Component<ScreenProps> {
    name: string;
    backend: Backend;
    sessionName: string;
    setScreen: (s: Screen) => void;

    constructor(pp: ScreenProps) {
        super(pp)
        this.name = "BASE Screen"
        this.backend = pp.backend
        this.sessionName = pp.sessionName
        this.setScreen = pp.setScreen
    }

    render(): JSX.Element {
        return <>ERROR</>
    }
}

class HomeScreen extends Screen {
    constructor(pp: Screen) {
        super(pp)
        this.name = "LoadingScreen"
    }

    render(): JSX.Element {
        return <LandingPage
            backend={this.backend}
            onGotoMenu={() => this.setScreen(new MenuScreen(this))}
        />
    }
}

class MenuScreen extends Screen {
    constructor(pp: Screen) {
        super(pp)
        this.name = "MenuScreen";
    }

    render() {
        return <RenderMenu
            backend={this.backend}
            sessionName={this.sessionName}
            onComplete={() => this.setScreen(new PantryScreen(this))}
            onGotoHome={() => this.setScreen(new HomeScreen(this))}
        />
    }
}

class PantryScreen extends Screen {
    constructor(pp: Screen) {
        super(pp)
        this.name = "PantryScreen";
    }

    render() {
        return <Pantry
            backend={this.backend}
            sessionName={this.sessionName}
            onBackToMenu={() => this.setScreen(new MenuScreen(this))}
            onComplete={() => this.setScreen(new ShoppingListScreen(this))}
            onGotoHome={() => this.setScreen(new HomeScreen(this))}
        />
    }
}

class ShoppingListScreen extends Screen {
    constructor(pp: Screen) {
        super(pp)
        this.name = "ShoppingListScreen";
    }

    render() {
        return <ShoppingList
            backend={this.backend}
            sessionName={this.sessionName}
            onBackToPantry={() => this.setScreen(new PantryScreen(this))}            
            onGotoHome={() => this.setScreen(new HomeScreen(this))}
        />
    }
}