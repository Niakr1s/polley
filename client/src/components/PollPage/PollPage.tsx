import React from 'react'
import { RouteComponentProps } from 'react-router'
import PollWithApi from '../Poll/PollWithApi'

interface MatchParams {
    uuid: string,
}

class PollPage extends React.Component<RouteComponentProps<MatchParams>> {
    render = () => {
        return (
            <PollWithApi uuid={this.props.match.params.uuid}></PollWithApi>
        )
    }
}

export default PollPage