import React from 'react'
import './TopBar.css'
import { useNavigate } from 'react-router-dom'

interface Props {
    left: JSX.Element
    right: JSX.Element
    logoOnClick?: () => void
    titleText?: string|any
}

export default function TopBar(pp: Props): JSX.Element {
    const style: React.CSSProperties = {
        width: '33%',
        display: 'flex',
    }

    const navigate = useNavigate()

    return (
        <div className='TopBar'>
            <div style={{...style, justifyContent: 'flex-start'}}>
                {pp.left}
            </div>
            <div style={{...style, justifyContent: 'center'}}>
                <Title
                    onClick={pp.logoOnClick || (() => navigate("/"))}
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
    onClick?: () => void
    titleText?: string|null
}

function Title(pp: TitleProps): JSX.Element {
    return <div key='1' className='Title'  onClick={pp.onClick}>
        <img src='/logo96.png' alt='logo' className="Logo"/>
        <div className='Text'>
            {pp.titleText || "La\xa0compra de\xa0l'Edu" /* \xa0 is a non-breaking space */} 
        </div>
    </div>
}