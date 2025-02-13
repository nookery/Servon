<script setup lang="ts">
import { toastState } from '../composables/useToastState'

const { toasts } = toastState
</script>

<template>
    <div class="fixed bottom-4 right-4 z-50 flex flex-col gap-2 toast toast-end">
        <div v-for="toast in toasts" :key="toast.id" :class="[
            'alert',
            'animate-fade-in-up',
            {
                'alert-success': toast.type === 'success',
                'alert-error': toast.type === 'error',
                'alert-info': toast.type === 'info'
            }
        ]">
            <span>{{ toast.message }}</span>
        </div>
    </div>
</template>

<style scoped>
@keyframes fadeInUp {
    from {
        opacity: 0;
        transform: translateY(1rem);
    }

    to {
        opacity: 1;
        transform: translateY(0);
    }
}

.animate-fade-in-up {
    animation: fadeInUp 0.3s ease-out;
}
</style>