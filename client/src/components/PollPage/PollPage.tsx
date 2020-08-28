import React from 'react'
import { ApiGetPoll, ApiPutPollChoices } from '../../api/api'
import { RouteComponentProps } from 'react-router'
import { IPoll } from '../../models/poll'
import Poll from '../Poll/Poll'
import { isExpired } from '../../util/date'

interface MatchParams {
    uuid: string,
}

interface IState {
    poll: IPoll | null,
}


class PollPage extends React.Component<RouteComponentProps<MatchParams>, IState> {
    constructor(props: Readonly<RouteComponentProps<MatchParams, import("react-router").StaticContext, import("history").History.UnknownFacade>>) {
        super(props)
        this.state = {
            poll: null,
        }
    }

    timeout?: ReturnType<typeof setTimeout> = undefined
    clearTimeout = () => {
        this.timeout && clearTimeout(this.timeout)
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

    getPoll = (once: boolean = false) => {
        ApiGetPoll(this.props.match.params.uuid).then(r => {
            let poll: IPoll = r.data
            this.setState({ poll })

            if (once) return

            if (!isExpired(poll.expires)) {
                this.clearTimeout()
                this.timeout = setTimeout(() => {
                    this.getPoll()
                }, 1000)
            }
        })
    }

    componentDidMount = () => {
        this.getPoll()
    }

    componentWillUnmount = () => {
        this.clearTimeout()
    }

    render = () => {
        if (!this.state.poll) return null
        return (
            <div>
                {this.state.poll
                    ? <Poll poll={this.state.poll} withVote={true} submitSelected={this.submitSelected}></Poll>
                    : <div>No such poll</div>
                }
            </div>
        )
    }
}

export default PollPage