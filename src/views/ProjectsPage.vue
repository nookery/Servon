<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { RiNodeTree, RiFolderOpenLine, RiFileListLine } from '@remixicon/vue'
import PageContainer from '../layouts/PageContainer.vue'
import ProjectTopology from '../components/projects/ProjectTopology.vue'
import ProjectFiles from '../components/projects/ProjectFiles.vue'
import ProjectLogs from '../components/projects/ProjectLogs.vue'
import type { SortBy, SortOrder } from '../types/FileInfo'

const route = useRoute()
const router = useRouter()

// 本地存储键名
const STORAGE_KEY = 'projectsLastPath'

// 默认项目目录
const defaultProjectPath = '/data/projects'
const currentPath = ref(route.query.path as string || '')

// 当前激活的 Tab
const activeTab = ref(route.query.tab as string || 'topology')

// 定义标签页
const tabs = [
    { key: 'topology', title: '项目拓扑', icon: RiNodeTree },
    { key: 'files', title: '文件管理', icon: RiFolderOpenLine },
    { key: 'logs', title: '项目日志', icon: RiFileListLine }
]

// 监听路径变化，更新 URL 和本地存储
watch(() => currentPath.value, (newPath) => {
    // 更新 URL
    router.replace({
        query: { ...route.query, path: newPath }
    })

    // 保存到本地存储
    localStorage.setItem(STORAGE_KEY, newPath)
})

// 监听 Tab 变化，更新 URL
watch(() => activeTab.value, (newTab) => {
    router.replace({
        query: { ...route.query, tab: newTab }
    })
})

// 初始化时加载路径和 Tab
onMounted(() => {
    // 加载 Tab
    if (route.query.tab) {
        activeTab.value = route.query.tab as string
    }

    // 优先级：URL 参数 > 本地存储 > 默认路径
    if (route.query.path) {
        // 如果 URL 中有路径参数，使用它
        currentPath.value = route.query.path as string
    } else {
        // 尝试从本地存储获取上次访问的路径
        const savedPath = localStorage.getItem(STORAGE_KEY)
        if (savedPath) {
            currentPath.value = savedPath
        } else {
            // 如果没有保存的路径，使用默认路径
            currentPath.value = defaultProjectPath
        }
    }
})
</script>

<template>
    <PageContainer title="项目管理" :tabs="tabs" v-model="activeTab">
        <!-- 项目拓扑 Tab -->
        <template #topology>
            <ProjectTopology :project-path="currentPath" />
        </template>

        <!-- 文件管理 Tab -->
        <template #files>
            <ProjectFiles />
        </template>

        <!-- 日志 Tab -->
        <template #logs>
            <ProjectLogs />
        </template>
    </PageContainer>
</template>