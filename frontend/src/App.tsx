import './App.css';
import React from 'react';
import ScreenStateMachine from './Screens/StateMachine.tsx';
import Backend from './Backend/Backend.ts'

function App(): JSX.Element {
  return (
    <ScreenStateMachine
      backend={new Backend()}
      sessionName='default'
    ></ScreenStateMachine>
  );
}

export default App;
