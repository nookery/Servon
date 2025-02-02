<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { NCard, NDataTable, useMessage } from 'naive-ui'
import axios from 'axios'

const processes = ref<any[]>([])
const message = useMessage()

const columns = [
    { title: 'PID', key: 'pid' },
    { title: '名称', key: 'name' },
    { title: 'CPU', key: 'cpu' },
    { title: '内存', key: 'memory' }
]

async function loadProcesses() {
    try {
        const res = await axios.get('http://localhost:8080/api/system/processes')
        processes.value = res.data
    } catch (error) {
        message.error('获取进程列表失败')
    }
}

onMounted(() => {
    loadProcesses()
})
</script>

<template>
    <n-card title="进程列表">
        <n-data-table :columns="columns" :data="processes" />
    </n-card>
</template>