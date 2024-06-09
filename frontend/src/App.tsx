import './App.css';
import React from 'react';
import ScreenStateMachine from './Screens/StateMachine.tsx';
import Backend from './Backend/Backend.ts'
import { State } from './State/State.tsx'

function App(): JSX.Element {
  return (
    <ScreenStateMachine
      backend={Backend.New()}
      globalState={new State()}
    ></ScreenStateMachine>
  );
}

export default App;
