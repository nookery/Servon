<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { NCard, NDataTable, NBreadcrumb, NBreadcrumbItem, useMessage } from 'naive-ui'
import axios from 'axios'

const files = ref<any[]>([])
const currentPath = ref('/')
const message = useMessage()

const columns = [
    { title: '名称', key: 'name' },
    { title: '类型', key: 'type' },
    { title: '大小', key: 'size' },
    { title: '修改时间', key: 'modTime' }
]

async function loadFiles(path: string) {
    try {
        const res = await axios.get(`/web_api/system/files?path=${path}`)
        files.value = res.data
    } catch (error) {
        message.error('获取文件列表失败')
    }
}

onMounted(() => {
    loadFiles(currentPath.value)
})
</script>

<template>
    <n-card title="文件管理">
        <n-breadcrumb>
            <n-breadcrumb-item>根目录</n-breadcrumb-item>
            <n-breadcrumb-item v-for="(part, index) in currentPath.split('/').filter(Boolean)" :key="index">
                {{ part }}
            </n-breadcrumb-item>
        </n-breadcrumb>
        <n-data-table :columns="columns" :data="files" />
    </n-card>
</template>