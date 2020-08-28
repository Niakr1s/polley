import React from 'react'
import { ApiGetPoll, ApiPutPollChoices } from '../../api/api'
import { IPoll } from '../../models/poll'
import Poll from './Poll'
import { isExpired } from '../../util/date'

interface IProps {
    uuid: string,
}

interface IState {
    poll: IPoll | null,
}

class PollWithApi extends React.Component<IProps, IState> {
    submitSelected = (selected: boolean[]) => {
        if (!this.state.poll) return
        let choices: string[] = []
        selected.forEach((choice, idx) => {
            if (!choice) return
            choices.push(this.state.poll!.choices[idx].text)
        })
        ApiPutPollChoices(this.state.poll.uuid, choices).then(() => { this.getPoll(true) }).catch(err => alert("couldn't commit votes: " + err))
    }

    pollTimeout?: ReturnType<typeof setTimeout>
    startPollTimeout = () => {
        this.clearPollTimeout()
        this.pollTimeout = setTimeout(() => {
            this.getPoll()
        }, 5000)
    }
    clearPollTimeout = () => {
        this.pollTimeout && clearTimeout(this.pollTimeout)
    }

    getPoll = (once: boolean = false) => {
        ApiGetPoll(this.props.uuid).then(r => {
            let poll: IPoll = r.data
            this.setState({ poll })

            if (once) return

            if (!isExpired(poll.expires)) {
                this.startPollTimeout()
            }
        })
    }

    componentDidMount = () => {
        this.getPoll()
    }

    componentWillUnmount = () => {
        this.clearPollTimeout()
    }

    render = () => {
        if (!this.state.poll) return null
        return (
            this.state.poll && <Poll poll={this.state.poll} withVote={true} submitSelected={this.submitSelected}></Poll>
        )
    }
}

export default PollWithApi