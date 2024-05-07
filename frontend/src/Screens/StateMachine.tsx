import React, { useState } from 'react'
import Backend from "../Backend/Backend.tsx";
import { State } from "../State/State.tsx";
import FirstLoad from "./FirstLoad/FirstLoad.tsx";
import RenderMenu from "./Menu/Menu.tsx";
import Pantry from "./Pantry/Pantry.tsx";
import PantryLoad from './PantryLoad/PantryLoad.tsx';

interface Props {
    backend: Backend;
    state: State;
}

function StateMachine(pp: Props): JSX.Element {
    const [screen, setScreen] = useState(new LoadMenuScreen(pp.backend, pp.state))
    const updateState = (newState) => {
        setScreen(newState)
    }

    return screen.render(updateState)
}

export default StateMachine

type ScreenSetter = React.Dispatch<React.SetStateAction<Screen>>

class Screen {
    backend: Backend;
    name: string;
    state: State;

    constructor(backend: Backend, state: State) {
        this.backend = backend
        this.state = state
    }

    render(setScreen: ScreenSetter): JSX.Element {
        throw Error("Not implemented")
    }

    next(): Screen {
        throw Error("Not implemented")
    }
}

class LoadMenuScreen extends Screen {
    constructor(backend: Backend, state: State) {
        super(backend, state)
        this.name = "LoadingScreen"
    }

    render(setScreen: ScreenSetter): JSX.Element {
        return <FirstLoad
            backend={this.backend}
            state={this.state}
            onComplete={() => setScreen(this.next())}
        />
    }

    next(): Screen {
        return new MenuScreen(this.backend, this.state);
    }
}

class MenuScreen extends Screen {
    constructor(backend: Backend, state: State) {
        super(backend, state)
        this.name = "MenuScreen";
    }

    render(setScreen: ScreenSetter) {
        return <RenderMenu
            backend={this.backend}
            state={this.state}
            onComplete={() => setScreen(this.next())}
        />
    }

    next(): Screen {
        return new LoadPantryScreen(this.backend, this.state)
    }
}

class LoadPantryScreen extends Screen {
    constructor(backend: Backend, state: State) {
        super(backend, state)
        this.name = "LoadPantryScreen"
    }

    render(setScreen: ScreenSetter): JSX.Element {
        return <PantryLoad
            backend={this.backend}
            state={this.state}
            onComplete={() => setScreen(this.next())}
        />
    }

    next(): Screen {
        return new PantryScreen(this.backend, this.state);
    }
}


class PantryScreen extends Screen {
    constructor(backend: Backend, state: State) {
        super(backend, state)
        this.name = "PantryScreen";
    }

    render(setScreen: ScreenSetter) {
        return <Pantry
            backend={this.backend}
            state={this.state}
        />
    }

    transition() {
        throw Error("Not implemented")
    }
}