import axios from 'axios'

export interface SystemResources {
    cpu_usage: number
    memory_usage: number
    disk_usage: number
}

export interface NetworkResources {
    download_speed: number
    upload_speed: number
}

export interface OSInfo {
    os_info: string
}

export interface SystemBasicInfo {
    hostname: string
    os: string
    platform: string
}

export interface Software {
    name: string
    status: string
}

export interface Process {
    pid: number
    command: string
    cpu: number
    memory: number
}

export interface IPInfo {
    local_ips: LocalIP[]
    public_ip: string
    public_ipv6: string
    dns_servers: string[]
    network_cards: NetworkCard[]
}

interface LocalIP {
    ip: string
    interface: string
    is_ipv6: boolean
    netmask: string
}

interface NetworkCard {
    name: string
    mac_address: string
    is_up: boolean
    mtu: number
    ips: string[]
}

export const systemAPI = {
    getResources: () =>
        axios.get<SystemResources>('/web_api/info/resources'),

    getOSInfo: () =>
        axios.get<OSInfo>('/web_api/info/os'),

    getNetworkResources: () =>
        axios.get<NetworkResources>('/web_api/info/network'),

    getCurrentUser: () =>
        axios.get<string>('/web_api/info/user'),

    getBasicInfo: () =>
        axios.get<SystemBasicInfo>('/web_api/info/basic'),

    // 获取软件列表
    getSoftwareList: () =>
        axios.get<string[]>('/web_api/soft'),

    // 获取软件状态
    getSoftwareStatus: (name: string) =>
        axios.get<{ status: string }>(`/web_api/soft/${name}/status`),

    // 安装/卸载软件
    manageSoftware: (name: string, action: 'install' | 'uninstall') =>
        axios.post(`/web_api/soft/${name}/${action}`),

    // 启动软件
    startSoftware: (name: string) =>
        axios.post(`/web_api/soft/${name}/start`),

    // 停止软件
    stopSoftware: (name: string) =>
        axios.post(`/web_api/soft/${name}/stop`),

    // 获取进程列表
    getProcesses: () =>
        axios.get<Process[]>('/web_api/processes'),

    // 结束进程
    killProcess: (pid: number) =>
        axios.post(`/web_api/processes/${pid}/kill`),

    getIPInfo: () =>
        axios.get<IPInfo>('/web_api/info/ip'),
} 