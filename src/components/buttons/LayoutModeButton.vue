<script setup lang="ts">
import { useLayoutModeStore } from '../../stores/layoutModeStore'
import IconButton from '../IconButton.vue'

const layoutModeStore = useLayoutModeStore()

const toggleLayoutMode = () => {
    layoutModeStore.toggleMode()
}
</script>

<template>
    <div class="tooltip tooltip-bottom" :data-tip="layoutModeStore.mode === 'web' ? '切换到OS布局' : '切换到Web布局'">
        <div v-if="layoutModeStore.mode === 'web'">
            <IconButton icon="ri-window-line" variant="ghost" circle @click="toggleLayoutMode" />
        </div>
        <div v-else>
            <IconButton icon="ri-layout-grid-line" variant="ghost" circle @click="toggleLayoutMode" />
        </div>
    </div>
</template>