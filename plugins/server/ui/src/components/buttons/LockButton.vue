<script setup lang="ts">
import { useLockStore } from '../../stores/lockStore'
import IconButton from '../IconButton.vue'
import { RiLockLine } from '@remixicon/vue'

const lockStore = useLockStore()

const handleLock = () => {
    lockStore.lock()
}
</script>

<template>
    <IconButton @click="handleLock" variant="ghost" circle>
        <RiLockLine class="w-4 h-4" />
    </IconButton>
</template>