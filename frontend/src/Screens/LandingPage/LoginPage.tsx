import React from 'react'
import Backend from '../../Backend/Backend';
import { useGoogleLogin, googleLogout } from '@react-oauth/google';
import { Footer, HelpText, PageHeader } from './LandingPage';
import './LandingPage.css';

interface Props {
    onLogIn: (creds: string) => void;
}

export default function LoginPage(props: Props) {
    return (
        <div className='LandingPage'>
            <PageHeader />
            <Login setCredential={(c: string) => props.onLogIn(c)} />
            <Footer />
        </div >
    )
}

export function Login(props: { setCredential: (credential: string) => void }) {
    const onError = () => {
        alert('Error al iniciar sessió')
    }

    const onSuccess = (creds?: string) => {
        if (creds) {
            props.setCredential(creds)
        } else {
            onError()
        }
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
                <GoogleLogin onSuccess={onSuccess} onError={onError} />
            </div>

        </div >
    )
}

function GoogleLogin(props: { onSuccess: (c: string) => void, onError: () => void }): JSX.Element {
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
                if (!creds.code) {
                    props.onError()
                }
                new Backend().Login()
                    .POST(creds.code)
                    .then(props.onSuccess)
                    .catch(props.onError)
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

export function GoogleLogout(props: { backend: Backend, onLogout: () => void }) {
    return (
        <div id='login'>
            <button
                onClick={() => {
                    props.backend.Logout().POST()
                        .then(() => {
                            Backend.IsMock()
                                ? console.log('Mock logout')
                                : googleLogout()
                            props.onLogout()
                        }).catch(() => {
                            alert('Error al tancar la sessió')
                        })
                }}
                id='google'
            >
                Tanca la sessió
            </button>
        </div >
    )
}
