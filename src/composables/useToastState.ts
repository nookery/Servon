import { ref } from 'vue'

interface ToastItem {
    id: number
    message: string
    type: 'success' | 'error' | 'info' | 'warning'
}

const toasts = ref<ToastItem[]>([])
let nextId = 0

export const toastState = {
    toasts,
    add(message: string, type: ToastItem['type'] = 'info') {
        const id = nextId++
        toasts.value.push({ id, message, type })
        setTimeout(() => {
            toasts.value = toasts.value.filter(t => t.id !== id)
        }, 3000)
    }
} 