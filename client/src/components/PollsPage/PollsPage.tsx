import React from 'react'
import { ApiGetUUIDs } from '../../api/api'
import { RouteComponentProps } from 'react-router'
import PollWithApi from '../Poll/PollWithApi'
import InfiniteScroll from 'react-infinite-scroller'

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
            total: Number.MAX_VALUE,
        }
    }

    getPolls = (page: number) => {
        console.log(`PollsPage:start loading`)
        ApiGetUUIDs(pageSize, page).then(r => {
            console.log(`PollsPage:got`, r)
            let { uuids, total }: { uuids: string[], total: number } = r.data
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
        return (
            <InfiniteScroll
                pageStart={0}
                loadMore={this.getPolls}
                hasMore={this.state.uuids.length < this.state.total}
                loader={<div>Loading...</div>}
                initialLoad={false}
            >
                {this.state.uuids.map(uuid => <PollWithApi key={uuid} uuid={uuid}></PollWithApi>)}
            </InfiniteScroll>
        )
    }
}

export default PollsPage