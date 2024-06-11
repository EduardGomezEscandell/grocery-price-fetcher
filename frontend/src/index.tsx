import React from 'react';
import ReactDOM from 'react-dom/client';
import Root from './Screens/Root.tsx';
import './index.css';

const root = document.getElementById('root')
if (!root) {
  throw new Error("No root element found")
}

ReactDOM.createRoot(root).render(
  <React.StrictMode>
    <Root />
  </React.StrictMode>
);

