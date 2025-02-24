
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

export interface IPInfo {
    local_ips: LocalIP[]
    public_ip: string
    public_ipv6: string
    dns_servers: string[]
    network_cards: NetworkCard[]
}