import React from 'react'
import './TopBar.css'

interface Props {
    components: { (): JSX.Element }[]
}

export default class TopBar extends React.Component<Props> {
    separator: JSX.Element = <></>

    render(): JSX.Element {
        return (
            <div className='TopBar'>
                {this.props.components.map(f => f())}
            </div>
        )
    }
}