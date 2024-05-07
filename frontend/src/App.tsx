import './App.css';
import { useState } from 'react';

import ScreenStateMachine from './Screens/StateMachine.tsx';
import Backend from './Backend/Backend.ts'
import { State } from './State/State.tsx'
import { Menu } from './State/State.tsx'

function App(): JSX.Element {
  const [menu, setMenu] = useState(new Menu())

  return ScreenStateMachine({
    backend: Backend.New(),
    state: new State().attachMenu(menu, setMenu)
  })
}

export default App;
