export interface IPollToSend {
    choices: string[],
    name: string,
    settings: {
        allowMultiple: number,
        timeoutMinutes: number,
        filter: string, // "", "ip", "cookie"
    }
}

export interface IPoll {
    uuid: string,
    name: string,
    expires: Date,
    allowMultiple: number,
    choices: Choice[],
}

export interface Choice {
    text: string,
    votes: number,
}