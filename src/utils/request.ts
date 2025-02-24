import axios from 'axios'

const axiosInstance = axios.create({
    baseURL: import.meta.env.VITE_API_BASE_URL || '/web_api',
    timeout: 10000
})

axiosInstance.interceptors.response.use(
    response => response.data,
    error => {
        const message = error.response?.data?.error || error.message
        return Promise.reject(new Error(message))
    }
)

export const request = {
    async get<T>(url: string, config?: any): Promise<T> {
        return axiosInstance.get(url, config)
    },
    async post<T>(url: string, config?: any): Promise<T> {
        return axiosInstance.post(url, config)
    }
} 