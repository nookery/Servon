<script setup lang="ts">
import { useLockStore } from '../../stores/lockStore'
import IconButton from '../IconButton.vue'

const lockStore = useLockStore()

const handleLock = () => {
    lockStore.lock()
}
</script>

<template>
    <IconButton @click="handleLock">
        <template #icon>
            <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect>
                <path d="M7 11V7a5 5 0 0 1 10 0v4"></path>
            </svg>
        </template>
        锁定
    </IconButton>
</template>