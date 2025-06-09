<script setup lang="ts">
import { useWindowStore } from '../stores/windowStore'
import Window from './Window.vue'

const windowStore = useWindowStore()
</script>

<template>
    <div class="window-manager">
        <Window v-for="window in windowStore.windows" :key="window.id" :window="window" />
    </div>
</template>

<style scoped>
.window-manager {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    pointer-events: none;
}

.window-manager :deep(.window) {
    pointer-events: auto;
}
</style>