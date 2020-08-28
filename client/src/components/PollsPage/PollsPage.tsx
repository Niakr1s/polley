import React from 'react'
import { ApiGetUUIDs } from '../../api/api'
import { RouteComponentProps } from 'react-router'
import PollWithApi from '../Poll/PollWithApi'

interface MatchParams {
    uuid: string,
}

interface IState {
    uuids: string[],
    total: number,
}

const pageSize = 10

class PollsPage extends React.Component<RouteComponentProps<MatchParams>, IState> {
    constructor(props: Readonly<RouteComponentProps<MatchParams, import("react-router").StaticContext, import("history").History.UnknownFacade>>) {
        super(props)
        this.state = {
            uuids: [],
            total: 0,
        }
    }

    getPolls = (page: number = 0) => {
        ApiGetUUIDs(pageSize, page).then(r => {
            console.log(`PollsPage:got`, r)
            let {uuids, total} = r.data
            this.setState(prevState => ({
                uuids: [...prevState.uuids, ...uuids],
                total,
            }))
        }).catch(err => alert(err))
    }

    componentDidMount = () => {
        this.getPolls(0)
    }

    render = () => {
        if (!this.state.uuids) return <div>No polls</div>
        return (
            <div>
                {this.state.uuids.map(uuid => <PollWithApi key={uuid} uuid={uuid}></PollWithApi>)}
            </div>
        )
    }
}

export default PollsPage