import React from 'react'
import { ApiGetPolls } from '../../api/api'
import { RouteComponentProps } from 'react-router'
import { IPoll } from '../../models/poll'
import Poll from '../Poll/Poll'

interface MatchParams {
    uuid: string,
}

interface IState {
    polls: IPoll[] | null,
}

const pageSize = 10

class PollsPage extends React.Component<RouteComponentProps<MatchParams>, IState> {
    constructor(props: Readonly<RouteComponentProps<MatchParams, import("react-router").StaticContext, import("history").History.UnknownFacade>>) {
        super(props)
        this.state = {
            polls: null,
        }
    }

    getPolls = (page: number = 0) => {
        ApiGetPolls(pageSize, page).then(r => {
            this.setState({ polls: r.data })
        }).catch(err => alert(err))
    }

    componentDidMount = () => {
        this.getPolls(0)
    }

    render = () => {
        if (!this.state.polls) return <div>No polls</div>
        return (
            <div>
                {this.state.polls.map(poll => <Poll key={poll.uuid} poll={poll} withVote={false} submitSelected={() => { }}></Poll>)}
            </div>
        )
    }
}

export default PollsPage