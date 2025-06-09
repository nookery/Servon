export interface Project {
    name: string
    domain: string
    upstream_url: string
    enabled: boolean
    config: Record<string, any>
}