import React from 'react'
import { RouteComponentProps } from 'react-router'
import PollWithApi from '../Poll/PollWithApi'
import ToMainPage from '../Links/ToMainPage'

interface MatchParams {
    uuid: string,
}

class PollPage extends React.Component<RouteComponentProps<MatchParams>> {
    render = () => {
        return (
            <div>
                <ToMainPage />
                <PollWithApi uuid={this.props.match.params.uuid}></PollWithApi>
            </div>
        )
    }
}

export default PollPage