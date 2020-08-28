import React from 'react'
import Polls from '../Polls/Polls'
import ToCreatePoll from '../Links/ToCreatePoll'

class PollsPage extends React.Component {
    render = () => {
        return (
            <div>
                <ToCreatePoll/>
                <Polls></Polls>
            </div>
        )
    }
}

export default PollsPage