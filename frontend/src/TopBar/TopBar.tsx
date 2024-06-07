import React from 'react'
import './TopBar.css'

interface Props {
    left: JSX.Element
    right: JSX.Element
    logoOnClick?: () => void
    titleOnClick?: () => void
    titleText?: string|any
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
                <Title
                    logoOnClick={pp.logoOnClick || (() => {})}
                    titleOnClick={pp.titleOnClick || (() => {})}
                    titleText={pp.titleText}
                    />
            </div>
            <div style={{...style, justifyContent: 'flex-end'}}>
                {pp.right}
            </div>
        </div>
    )
}

interface TitleProps {
    logoOnClick?: () => void
    titleOnClick?: () => void
    titleText?: string|null
}

function Title(pp: TitleProps): JSX.Element {
    return <div key='1' className='Title'>
        <img src='/logo64.png' alt='logo' className="Logo" onClick={pp.logoOnClick}/>
        <div className='Text' onClick={pp.titleOnClick}>
            {pp.titleText || "La\xa0compra de\xa0l'Edu" /* \xa0 is a non-breaking space */} 
        </div>
    </div>
}