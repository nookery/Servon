<script setup lang="ts">
import { ref, onMounted } from 'vue'
import FileManager from '../files/FileManager.vue'
import type { SortBy, SortOrder } from '../../types/FileInfo'

// 本地存储键名
const STORAGE_KEY = 'projectFilesLastPath'
// 默认项目目录
const defaultProjectPath = '/data/projects'

// 移除 props 依赖
const currentPath = ref(defaultProjectPath)
const currentSort = ref<SortBy>('name')
const sortOrder = ref<SortOrder>('asc')

onMounted(() => {
    // 从本地存储加载上次访问的路径，如果没有则使用默认路径
    const savedPath = localStorage.getItem(STORAGE_KEY)
    currentPath.value = savedPath || defaultProjectPath
})

// 监听路径变化，保存到本地存储
function handlePathChange(newPath: string) {
    currentPath.value = newPath
    localStorage.setItem(STORAGE_KEY, newPath)
}
</script>

<template>
    <FileManager v-model:path="currentPath" :initial-path="currentPath" :show-breadcrumbs="true" :show-toolbar="true"
        :show-pagination="true" :show-shortcuts="false" :read-only="false" v-model:sort-by="currentSort"
        v-model:sort-order="sortOrder" @update:path="handlePathChange" />
</template>