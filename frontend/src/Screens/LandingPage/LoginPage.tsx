import React from 'react'
import Backend from '../../Backend/Backend';
import { useGoogleLogin, googleLogout as _googleLogout } from '@react-oauth/google';
import { Footer, HelpText, PageHeader } from './LandingPage';
import './LandingPage.css';

interface Props {
    logIn: (creds: string) => Promise<void>;
}

export default function LoginPage(props: Props) {
    return (
        <div className='LandingPage'>
            <PageHeader />
            <Login onLogin={props.logIn} />
            <Footer />
        </div >
    )
}

export function Login(props: { onLogin: (credential: string) => Promise<void> }) {
    const onError = () => {
        alert('Error al iniciar sessió')
    }

    return (
        <div id='login'>
            <div>
                <HelpText />
            </div>
            <div>
                <p>
                    Inicia sessió amb el teu compte de Google per començar a planificar la teva compra setmanal.
                </p>
                <GoogleLogin onSuccess={props.onLogin} onError={onError} />
            </div>

        </div >
    )
}

function GoogleLogin(props: { onSuccess: (c: string) => Promise<void>, onError: () => void }): JSX.Element {
    let login: () => void

    if (Backend.IsMock()) {
        login = () => {
            console.log('Mock login success')
            props.onSuccess('mock-credential-123')
        }
    } else {
        login = useGoogleLogin({
            flow: 'auth-code',
            onSuccess: (creds) => {
                props.onSuccess(creds.code).catch(props.onError)
            },
            onError: props.onError
        })
    }

    return (
        <button id='google' onClick={login}>
            Inicia la sessió amb Google
        </button>
    )
}

export function GoogleLogout(props: { backend: Backend, logOut: () => Promise<void> }) {
    return (
        <div id='login'>
            <button
                onClick={() => props.logOut().then(
                    () => googleLogout(),
                    () => alert('Error al tancar la sessió')
                )}
                id='google'
            >
                Tanca la sessió
            </button>
        </div >
    )
}

function googleLogout() {
    Backend.IsMock() ? console.log('Mock logout') : _googleLogout()
}