import { reactive } from 'vue'

export type ConfirmType = 'warning' | 'error' | 'info'

interface ConfirmOptions {
    title: string
    message: string
    type?: ConfirmType
    confirmText?: string
    cancelText?: string
}

interface ConfirmState {
    show: boolean
    options: ConfirmOptions
    resolve: ((value: boolean) => void) | null
}

export const confirmState = reactive<ConfirmState>({
    show: false,
    options: {
        title: '',
        message: '',
        type: 'warning',
        confirmText: '确认',
        cancelText: '取消'
    },
    resolve: null
}) 