import React from 'react'
import { ApiGetUUIDs } from '../../api/api'
import PollWithApi from '../Poll/PollWithApi'
import InfiniteScroll from 'react-infinite-scroller'

interface IState {
    uuids: string[],
    total: number,
}

const pageSize = 10

class Polls extends React.Component<any, IState> {
    constructor(props: any) {
        super(props)
        this.state = {
            uuids: [],
            total: Number.MAX_VALUE,
        }
    }

    getPolls = (page: number) => {
        ApiGetUUIDs(pageSize, page).then(r => {
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
                initialLoad={false}
            >
                {this.state.uuids.map(uuid => <PollWithApi key={uuid} uuid={uuid}></PollWithApi>)}
            </InfiniteScroll>
        )
    }
}

export default Polls