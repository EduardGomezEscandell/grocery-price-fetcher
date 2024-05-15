import React from 'react'
import Backend from '../../Backend/Backend'
import { State } from '../../State/State'

enum Phase {
    Idle = 0,
    Saving,
    Done,
}

interface Props {
    backend: Backend
    globalState: State
    className?: string
}

export default class SaveButton extends React.Component<Props> {
    state: { phase: Phase }

    constructor(pp: Props) {
        super(pp)
        this.state = { phase: Phase.Idle }
    }

    render(): JSX.Element {
        const { text, style } = (() => {
            switch (this.state.phase) {
                case Phase.Idle:
                    return { text: 'Guardar', style: {} }
                case Phase.Saving:
                    return { text: 'Guardant...', style: { backgroundColor: 'gray' } }
                case Phase.Done:
                    return { text: 'Guardat', style: { color: 'green' } }
            }
        })()

        return (
            <button
                className={this.props.className}
                onClick={this.onClick.bind(this)}
                style={{
                    ...style,
                    width: '100px',
                }}
            >{text}</button>
        )
    }

    private async onClick(): Promise<void> {
        if (this.state.phase !== Phase.Idle) {
            return
        }

        this.setState({ phase: Phase.Saving })
        await Promise.all([
            this.save(),
            this.wait(200)
        ])
        this.setState({ phase: Phase.Done })
        await this.wait(1000)
        this.setState({ phase: Phase.Idle })
    }

    private async save(): Promise<any> {
        return this.props.backend
            .Pantry()
            .POST({
                name: '', // Let the backend handle the name for now
                contents: this.props.globalState.shoppingList.ingredients
                    .filter(i => i.have > 0)
                    .map(i => {
                        return { name: i.name, amount: i.have }
                    })
            })
    }

    private async wait(ms: number): Promise<any> {
        return new Promise(resolve => setTimeout(resolve, ms))
    }
}