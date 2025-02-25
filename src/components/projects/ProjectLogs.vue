<script setup lang="ts">
import { ref, onMounted } from 'vue'
import FileManager from '../files/FileManager.vue'
import LogView from '../logs/LogView.vue'
import type { SortBy, SortOrder } from '../../types/FileInfo'

// 默认日志路径
const logsPath = ref('/data/logs')
const logsSort = ref<SortBy>('modTime')
const logsSortOrder = ref<SortOrder>('desc')

onMounted(() => {
    // 可以在这里添加初始化逻辑
})
</script>

<template>
    <LogView current-dir="" />
</template>