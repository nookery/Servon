import axios from 'axios'

interface GitHubSetupParams {
    name: string
    description?: string
}

export const githubAPI = {
    setup: (params: GitHubSetupParams) =>
        axios.post('/web_api/github/setup', params, {
            responseType: 'text'
        })
} 