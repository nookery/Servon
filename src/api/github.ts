import axios from 'axios'

interface GitHubSetupParams {
    name: string
    description?: string
}

export interface WebhookData {
    id: string
    type: string
    timestamp: number
    payload: any
}

export const githubAPI = {
    setup: (params: GitHubSetupParams) =>
        axios.post('/web_api/github/setup', params, {
            responseType: 'text'
        }),
    getWebhooks: () =>
        axios.get<WebhookData[]>('/web_api/github/webhooks'),
} 