import axios from 'axios'

export interface UserInfo {
    username: string
}

export interface User {
    username: string
    groups: string[]
    shell: string
    home_dir: string
    create_time: string
    last_login: string
    sudo: boolean
}

export interface NewUser {
    username: string
    password: string
}

export async function getUsers() {
    const res = await axios.get('/web_api/users')
    return res.data
}

export async function createUser(user: NewUser) {
    await axios.post('/web_api/users', user)
}

export async function deleteUser(username: string) {
    await axios.delete(`/web_api/users/${username}`)
} 