import './App.css';
import React from 'react';
import StateMachine from './StateMachine/StateMachine.jsx';
import SelectBackend from './Backend/Backend.js'

function App() {
  const backend = SelectBackend()
  return <StateMachine backend={backend}/>
}

export default App;
