import { request } from './request'

// 获取服务列表
export const getServiceList = async () => {
    const response = await request.get('/api/services')
    return response.data
}

// 启动服务
export const startService = async (serviceName: string) => {
    const response = await request.post(`/api/services/${serviceName}/start`)
    return response.data
}

// 停止服务
export const stopService = async (serviceName: string) => {
    const response = await request.post(`/api/services/${serviceName}/stop`)
    return response.data
}

// 重启服务
export const restartService = async (serviceName: string) => {
    const response = await request.post(`/api/services/${serviceName}/restart`)
    return response.data
}

// 获取服务配置
export const getServiceConfig = async (serviceName: string) => {
    const response = await request.get(`/api/services/${serviceName}/config`)
    return response.data
}

// 更新服务配置
export const updateServiceConfig = async (serviceName: string, config: any) => {
    const response = await request.put(`/api/services/${serviceName}/config`, config)
    return response.data
}

// 获取服务日志
export const getServiceLogs = async (serviceName: string, lines: number = 100) => {
    const response = await request.get(`/api/services/${serviceName}/logs`, {
        params: { lines }
    })
    return response.data
}

// 获取服务详情
export const getServiceDetails = async (serviceName: string) => {
    const response = await request.get(`/api/services/${serviceName}/details`)
    return response.data
} 