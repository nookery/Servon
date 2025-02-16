<script setup lang="ts">
import { useLogViewerStore } from '../stores/logViewer'
import IconButton from '../components/IconButton.vue'

const logViewerStore = useLogViewerStore()
</script>

<template>
    <IconButton @click="logViewerStore.toggleVisibility" icon="ri-file-list-line" variant="ghost" circle
        title="日志查看器" />
</template>