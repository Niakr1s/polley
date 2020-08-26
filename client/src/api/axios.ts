import Axios from 'axios'
import baseURL from './baseUrl'

const axios = Axios.create({
    baseURL: baseURL,
    withCredentials: true,
})

export default axios