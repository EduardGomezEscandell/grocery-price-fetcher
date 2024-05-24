import React from 'react'
import './SaveButton.css'

enum Phase {
    Idle = 0,
    Busy,
    Done,
    Error,
}

interface Props {
    onSave?: () => Promise<any>
    onAccept?: () => void
    onReject?: (reason: any) => void

    baseTxt: string
    onSaveTxt: string
    onAcceptTxt: string
    onRejectTxt: string
}

export default class SaveButton extends React.Component<Props> {
    state: { phase: Phase }

    constructor(pp: Props) {
        super(pp)
        this.state = { phase: Phase.Idle }
    }

    render(): JSX.Element {
        const { text, id } = this.getPhaseData(this.state.phase)
        return <button className='save-button' id={id} onClick={this.cycle.bind(this)}>{text}</button>
    }

    private getPhaseData(phase: Phase): { text: string, id: string } {
        switch (phase) {
            case Phase.Idle:
                return { text: this.props.baseTxt, id: 'idle' }
            case Phase.Busy:
                return { text: this.props.onSaveTxt, id: 'busy' }
            case Phase.Done:
                return { text: this.props.onAcceptTxt, id: 'done' }
            case Phase.Error:
                return { text: this.props.onRejectTxt, id: 'error' }
        }
    }

    private setPhase(phase: Phase): void {
        this.setState({ phase })
    }

    private async cycle(): Promise<void> {
        if (this.state.phase !== Phase.Idle) {
            return
        }

        this.setPhase(Phase.Busy)

        this.onClick()
            .then(() => {
                this.setPhase(Phase.Done)
                return this.onAccept()
            })
            .catch((reason: any) => {
                this.setPhase(Phase.Error)
                return this.onReject(reason)
            })
            .finally(() => this.setPhase(Phase.Idle))
    }

    private static wait(ms: number): Promise<any> {
        return new Promise(resolve => setTimeout(resolve, ms))
    }

    private onClick(): Promise<void> {
        if (this.props.onSave === undefined) {
            return SaveButton.wait(1000)
        }
        return this.props.onSave()
    }

    private async onAccept(): Promise<Promise<void>> {
        if (this.props.onAccept === undefined) {
            return await SaveButton.wait(1000)
        }
        return this.props.onAccept()
    }

    private async onReject(reason: any): Promise<void> {
        if (this.props.onReject === undefined) {
            console.error(`SaveButton: ${reason}`)
            return await SaveButton.wait(1000)
        }
        return this.props.onReject(reason)
    }
}