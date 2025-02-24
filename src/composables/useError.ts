import { inject, type Ref } from 'vue'

export function useError() {
    const showError = inject<(message: string) => void>('showError')

    if (!showError) {
        throw new Error('useError must be used within App.vue')
    }

    return {
        error: showError
    }
} 