import axios from 'axios'

export async function getPorts() {
    const res = await axios.get('/web_api/ports')
    return res.data
} 