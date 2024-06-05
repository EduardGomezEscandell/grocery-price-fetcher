import React from 'react'
import './TopBar.css'

interface Props {
    left: JSX.Element
    right: JSX.Element
}

export default function TopBar(pp: Props): JSX.Element {
    const style: React.CSSProperties = {
        width: '33%',
        display: 'flex',
    }

    return (
        <div className='TopBar'>
            <div style={{...style, justifyContent: 'flex-start'}}>
                {pp.left}
            </div>
            <div style={{...style, justifyContent: 'center'}}>
                <Title />
            </div>
            <div style={{...style, justifyContent: 'flex-end'}}>
                {pp.right}
            </div>
        </div>
    )
}

function Title(): JSX.Element {
    return <div key='1' className='Title'>
        <img src='/logo64.png' alt='logo' className="Logo"/>
        <div className='Text'>
            La&nbsp;compra de&nbsp;l'Edu
        </div>
    </div>
}