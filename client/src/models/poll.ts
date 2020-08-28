// send to server
export interface IPollToSend {
    choices: string[],
    name: string,
    settings: {
        allowMultiple: number,
        timeoutMinutes: number,
        filter: string, // "", "ip", "cookie"
    }
}

// got from server
export interface IPoll {
    uuid: string,
    name: string,
    expires: string,
    allowMultiple: number,
    voteAllowed: boolean,
    choices: Choice[],
}

export interface Choice {
    text: string,
    votes: number,
}

export const choicesRemained = (poll: IPoll, selected?: number): number => {
    return (poll.allowMultiple || 0) - (selected || 0)
}

export const getPercentOfNthChoice = (poll: IPoll, idx: number): number => {
    const sum: number = poll.choices.reduce((acc, choice) => acc + choice.votes, 0)
    const res: number = Math.round(poll.choices[idx].votes / sum * 100)
    return res
}