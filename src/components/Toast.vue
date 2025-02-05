<script setup lang="ts">
import { toastState } from '../composables/useToastState'

const { toasts } = toastState
</script>

<template>
    <div class="toast toast-end">
        <div v-for="toast in toasts" :key="toast.id" :class="{
            'alert': true,
            'alert-success': toast.type === 'success',
            'alert-error': toast.type === 'error',
            'alert-info': toast.type === 'info'
        }">
            <span>{{ toast.message }}</span>
        </div>
    </div>
</template>

<style scoped>
.toast {
    position: fixed;
    bottom: 1rem;
    right: 1rem;
    z-index: 1000;
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
}
</style>