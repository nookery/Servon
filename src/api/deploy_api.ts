import axios from 'axios'

// 部署日志的数据接口定义
export interface DeployLog {
    id: string
    timestamp: string
    status: string
    message: string
}

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