<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import PageContainer from '../layouts/PageContainer.vue'
import FileManager from '../components/files/FileManager.vue'
import type { SortBy, SortOrder } from '../models/FileInfo'

const route = useRoute()
const router = useRouter()

const currentPath = ref(route.query.path as string || '/')
const currentSort = ref<SortBy>('name')
const sortOrder = ref<SortOrder>('asc')

// 监听路径变化，更新 URL
watch(() => currentPath.value, (newPath) => {
    router.replace({
        query: { ...route.query, path: newPath }
    })
})

// 初始化时从 URL 读取路径
onMounted(() => {
    if (route.query.path) {
        currentPath.value = route.query.path as string
    }
})
</script>

<template>
    <PageContainer title="文件管理">
        <FileManager v-model:path="currentPath" :initial-path="currentPath" :show-breadcrumbs="true"
            :show-toolbar="true" :show-pagination="true" :show-shortcuts="true" :read-only="false"
            v-model:sort-by="currentSort" v-model:sort-order="sortOrder" />
    </PageContainer>
</template>