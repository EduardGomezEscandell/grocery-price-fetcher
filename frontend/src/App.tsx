import './App.css';
import React, { useState } from 'react';
import ScreenStateMachine from './Screens/StateMachine.tsx';
import Backend from './Backend/Backend.ts'
import { State, Menu } from './State/State.tsx'

function App(): JSX.Element {
  const [menu, setMenu] = useState(new Menu())

  return (
    <ScreenStateMachine
      backend={Backend.New()}
      globalState={new State().attachMenu(menu, setMenu)}
    ></ScreenStateMachine>
  );
}

export default App;
