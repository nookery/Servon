import axios from 'axios'
import type { FileInfo } from '../models/FileInfo'
import type { GitHubRepo, GitHubSetupParams, WebhookData } from '../models/GitHubTypes'

export async function getAuthorizedRepos(): Promise<GitHubRepo[]> {
    const res = await axios.get('/web_api/integrations/github/repos')
    return res.data
}

export async function getGitHubLogs(): Promise<FileInfo[]> {
    const res = await axios.get('/web_api/integrations/github/logs')
    return res.data
}

export async function getFileContent(path: string): Promise<string> {
    const res = await axios.get(`/web_api/files/content`, {
        params: { path }
    })
    return res.data
}

export const githubAPI = {
    setup: (params: GitHubSetupParams) =>
        axios.post<string>('/web_api/github/setup', params, {
            responseType: 'text',
            headers: {
                'Content-Type': 'application/json'
            }
        }),
    getWebhooks: () =>
        axios.get<WebhookData[]>('/web_api/github/webhooks'),
} 