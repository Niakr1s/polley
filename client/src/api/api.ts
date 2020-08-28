import { IPollToSend } from '../models/poll'
import axios from './axios'
import qs from 'querystring'

export const ApiCreatePoll = (poll: IPollToSend) => {
    return axios.post("/createPoll", poll)
}

export const ApiGetPoll = (uuid: string) => {
    return axios.get(`/poll/${uuid}`)
}
export const ApiGetUUIDs = (pageSize: number = 10, page: number = 0) => {
    return axios.get(`/getUUIDs?${qs.stringify({pageSize, page})}`)
}

export const ApiPutPollChoices = (uuid: string, choices: string[]) => {
    return axios.put(`/poll/${uuid}`, { choices })
}
