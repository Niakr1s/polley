import React from 'react'
import { ApiGetPoll, ApiPutPollChoices } from '../../api/api'
import { IPoll } from '../../models/poll'
import Poll from './Poll'
import { isExpired } from '../../util/date'
import { Loader } from 'semantic-ui-react'

interface IProps {
    uuid: string,
}

interface IState {
    poll: IPoll | null,
    error?: string,
}

class PollWithApi extends React.Component<IProps, IState> {
    constructor(props: Readonly<IProps>) {
        super(props)
        this.state = {
            poll: null,
        }
    }

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
    startPollInterval = () => {
        this.clearPollInterval();
        this.getPoll();
        this.pollTimeout = setInterval(() => {
            this.getPoll()
        }, 5000);
    }
    clearPollInterval = () => {
        this.pollTimeout && clearTimeout(this.pollTimeout)
    }

    getPoll = (once: boolean = false) => {
        console.log('getPoll', this.props.uuid)
        ApiGetPoll(this.props.uuid)
            .then(r => {
                let poll: IPoll = r.data
                this.setState({ poll })

                if (once || isExpired(poll.expires)) {
                    this.clearPollInterval();
                }
            })
            .catch((error) => {
                console.log('error', error.message)
                this.setState({ error: error.message })
            })
    }

    componentDidMount = () => {
        this.startPollInterval()
    }

    componentWillUnmount = () => {
        this.clearPollInterval()
    }

    render = () => {
        if (this.state.error) {
            return (
                <div>Error occured while loading poll: {this.state.error}</div>
            )
        }
        if (!this.state.poll) {
            return (
                <Loader />
            )
        }
        return (
            <Poll poll={this.state.poll} withVote={this.state.poll.voteAllowed} submitSelected={this.submitSelected}></Poll>
        )
    }
}

export default PollWithApi