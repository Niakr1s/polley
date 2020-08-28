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
    expires: Date,
    allowMultiple: number,
    voteAllowed: boolean,
    choices: Choice[],
}

export interface Choice {
    text: string,
    votes: number,
}