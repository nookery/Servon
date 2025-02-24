export interface GitHubRepo {
    id: number
    name: string
    full_name: string
    description: string
    private: boolean
}

export interface GitHubSetupParams {
    base_url: string
    name: string
    description?: string
}

export interface WebhookData {
    id: string
    type: string
    payload: any
    created_at: string
} 