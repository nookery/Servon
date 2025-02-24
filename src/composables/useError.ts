import { inject } from 'vue'

export function useError() {
    const showError = inject<(message: string, showInHeader?: boolean) => void>('showError')

    if (!showError) {
        throw new Error('useError must be used within App.vue')
    }

    return {
        // 默认在头部显示
        error: (message: string, showInHeader = true) => showError(message, showInHeader),

        // 在右下角显示的便捷方法
        notifyError: (message: string) => showError(message, false)
    }
} 