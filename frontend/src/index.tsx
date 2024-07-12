import React from 'react';
import ReactDOM from 'react-dom/client';
import { GoogleOAuthProvider } from '@react-oauth/google';
import Root from './Screens/Root';  
import './index.css';

const root = document.getElementById('root')
if (!root) {
  throw new Error("No root element found")
}

ReactDOM.createRoot(root).render(
  <GoogleOAuthProvider
    clientId={import.meta.env.VITE_APP_GOOGLE_CLIENT_ID}
  >
    <React.StrictMode>
      <Root />
    </React.StrictMode>
  </GoogleOAuthProvider>
);

