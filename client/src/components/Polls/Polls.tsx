import React from 'react'
import { IPoll } from '../../models/poll'
import Poll from '../Poll/Poll'

interface IProps {
    polls: IPoll[],
}

class Polls extends React.Component<IProps> {
    render = () => {
        return (
            <div >
                {this.props.polls.map(poll => (
                    <Poll poll={poll} withVote={false} submitSelected={() => { }}></Poll>
                ))}
            </div>
        )
    }
}

export default Polls