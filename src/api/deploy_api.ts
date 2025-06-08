import axios from 'axios'
import type { DeployLog } from '../types/Deploy'

// 获取所有部署日志
export async function getDeployLogs() {
    const res = await axios.get('/web_api/deploy/logs')
    return res.data as DeployLog[]
}

// 获取单个部署日志
export async function getDeployLog(id: string) {
    const res = await axios.get('/web_api/deploy/log', {
        params: { id }
    })
    return res.data as DeployLog
}

// 删除部署日志
export async function deleteDeployLog(id: string) {
    await axios.delete('/web_api/deploy/log', {
        params: { id }
    })
}

// 部署仓库
export async function deployRepository(id: string) {
    const res = await axios.post('/web_api/deploy/repository', null, {
        params: { id }
    })
    return res.data
}
