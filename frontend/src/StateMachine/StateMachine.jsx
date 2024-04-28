import Menu from "../Menu/Menu";

import React, { useState } from 'react'
import Pantry from "../Pantry/Pantry";

function StateMachine(pp) {
    const [state, setState] = useState(new BaseState("default", pp.backend))
    const updateState = (newState) => {
        setState(newState)
    }

    return (
        state.getComponent(updateState)
    )
}

export default StateMachine


class BaseState {
    constructor(name, backend) {
        this.name = name;
        this.backend = backend;
    }

    getComponent(setState) {
        setState(this.transition());
        return <div>Loading...</div>
    }

    transition() {
        return new MenuState(this.backend);
    }
}


class MenuState extends BaseState {
    constructor(backend) {
        super("MenuState", backend);
    }

    getComponent(setState) {
        return <Menu backend={this.backend} onComplete={(recipes) => setState(this.transition(recipes))} />
    }

    transition(recipes) {
        throw Error("Not implemented")
    }
}