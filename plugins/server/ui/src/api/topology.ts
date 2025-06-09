import axios from 'axios'
import type { Project } from '../types/Project'

export const topologyAPI = {
    // 获取所有网关
    getGateways: () =>
        axios.get<string[]>('/web_api/topology/gateways'),

    // 获取指定网关的项目列表
    getProjects: (gateway: string) =>
        axios.get<Project[]>(`/web_api/topology/gateways/${gateway}/projects`),

    // 添加项目到网关
    addProject: (gateway: string, project: Project) =>
        axios.post(`/web_api/topology/gateways/${gateway}/projects`, project),

    // 从网关移除项目
    removeProject: (gateway: string, projectName: string) =>
        axios.delete(`/web_api/topology/gateways/${gateway}/projects/${projectName}`),

    // 获取网关配置
    getGatewayConfig: (gateway: string) =>
        axios.get<{ config: string }>(`/web_api/topology/gateways/${gateway}/config`),

    // 设置网关配置
    setGatewayConfig: (gateway: string, config: string) =>
        axios.put<void>(`/web_api/topology/gateways/${gateway}/config`, { config }),
} 