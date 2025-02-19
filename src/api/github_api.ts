import axios from 'axios'
import type { FileInfo } from '../models/FileInfo'

export interface GitHubRepo {
    id: number
    name: string
    fullName: string
    description: string
    private: boolean
    htmlUrl: string
}

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