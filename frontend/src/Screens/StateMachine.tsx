import React from 'react'
import Backend from "../Backend/Backend.tsx";
import { State } from "../State/State.tsx";
import FirstLoad from "./FirstLoad/FirstLoad.tsx";
import RenderMenu from "./Menu/Menu.tsx";
import Pantry from "./Pantry/Pantry.tsx";
import PantryLoad from './PantryLoad/PantryLoad.tsx';

interface Props {
    backend: Backend;
    globalState: State;
}

export default class StateMachine extends React.Component<Props> {
    state: { screen: Screen }

    constructor(props: Props) {
        super(props)
        this.state = {
            screen: new Screen({
                ...props,
                setScreen: (s: Screen) => this.setState({ screen: s })
            })
        }
    }

    render(): JSX.Element {
        return this.state.screen.render()
    }
}

interface ScreenProps extends Props {
    setScreen: (s: Screen) => void;
}

class Screen extends React.Component<ScreenProps> {
    name: string;
    backend: Backend;
    globalState: State;
    setScreen: (s: Screen) => void;

    constructor(p: ScreenProps) {
        super(p)
        this.name = "BASE Screen"
        this.backend = p.backend
        this.globalState = p.globalState
        this.setScreen = p.setScreen
    }

    render(): JSX.Element {
        this.setScreen(new LoadMenuScreen(this))
        return <>Loading...</>
    }
}

class LoadMenuScreen extends Screen {
    constructor(pp: Screen) {
        super(pp)
        this.name = "LoadingScreen"
    }

    render(): JSX.Element {
        return <FirstLoad
            backend={this.backend}
            globalState={this.globalState}
            onComplete={() => this.setScreen(new MenuScreen(this))}
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
            globalState={this.globalState}
            onComplete={() => this.setScreen(new LoadPantryScreen(this))}
        />
    }
}

class LoadPantryScreen extends Screen {
    constructor(pp: Screen) {
        super(pp)
        this.name = "LoadPantryScreen"
    }

    render(): JSX.Element {
        return <PantryLoad
            backend={this.backend}
            state={this.globalState}
            onComplete={() => this.setScreen(new PantryScreen(this))}
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
            state={this.globalState}
            onBackToMenu={() => this.setScreen(new LoadMenuScreen(this))}
        />
    }
}