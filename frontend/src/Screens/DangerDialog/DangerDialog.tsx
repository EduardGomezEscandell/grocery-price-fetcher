import React from 'react'

export default function DangerDialog(props: {
    onAccept: () => void
    onReject: () => void
    children?: JSX.Element[]
}): JSX.Element {
    if (props.children && props.children.length > 2) {
        throw new Error('DangerDialog can only have 2 children')
    }

    const header = (() => {
        const defaultH = <h3 id='header'>Confirmació</h3>
        if (!props.children) {
            return defaultH
        }
        const c = props.children.find(c => c.props.id === 'header')
        if (!c) {
            return defaultH
        }

        return c
    })()

    const body = (() => {
        const defaultH = <div id='body'>Procedir amb l'acció?</div>
        if (!props.children) {
            return defaultH
        }
        const c = props.children.find(c => c.props.id === 'body')
        if (!c) {
            return defaultH
        }

        return c
    })()

    return (
        <dialog open className='danger-dialog'>
            {header}
            {body}
            <div id='footer'>
                <button id='dialog-left' onClick={props.onReject}>No</button>
                <button id='dialog-right' onClick={props.onAccept}>Sí</button>
            </div>
        </dialog>
    )
}