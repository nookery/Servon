<script setup lang="ts">
defineProps<{
    type?: 'info' | 'success' | 'warning' | 'error'
    message: string
}>()
</script>

<template>
    <div :class="[
        'alert',
        {
            'alert-info': type === 'info',
            'alert-success': type === 'success',
            'alert-warning': type === 'warning',
            'alert-error': type === 'error'
        }
    ]">
        {{ message }}
    </div>
</template>