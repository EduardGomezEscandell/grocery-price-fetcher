import './App.css';
import ScreenStateMachine from './Screens/StateMachine.tsx';
import Backend from './Backend/Backend.ts'
import { State } from './State/State.tsx'

function App(): JSX.Element {
  return ScreenStateMachine({
    backend: Backend.New(),
    state: new State()
  })
}

export default App;
