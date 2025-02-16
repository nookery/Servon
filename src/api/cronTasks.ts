import axios from 'axios'

export interface CronTask {
    id: number
    name: string
    command: string
    schedule: string
    description: string
    enabled: boolean
    last_run?: string
    next_run?: string
}

export async function getTasks() {
    const res = await axios.get('/web_api/cron/tasks')
    return res.data
}

export async function createTask(task: CronTask) {
    const res = await axios.post('/web_api/cron/tasks', task)
    return res.data
}

export async function updateTask(task: CronTask) {
    const res = await axios.put(`/web_api/cron/tasks/${task.id}`, task)
    return res.data
}

export async function deleteTask(id: number) {
    await axios.delete(`/web_api/cron/tasks/${id}`)
}

export async function toggleTask(id: number) {
    await axios.post(`/web_api/cron/tasks/${id}/toggle`)
} 