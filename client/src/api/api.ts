import { IPollToSend } from '../models/poll'
import axios from './axios'

export const ApiCreatePoll = (poll: IPollToSend) => {
    return axios.post("/createPoll", poll)
}

export const ApiGetPoll = (uuid: string) => {
    return axios.get(`poll/${uuid}`)
}

export const ApiPutPollChoices = (uuid: string, choices: string[]) => {
    return axios.put(`poll/${uuid}`, { choices })
}