<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { NCard, NDataTable, useMessage, NButton, NSpace } from 'naive-ui'
import axios from 'axios'

const processes = ref<any[]>([])
const message = useMessage()
const showRawData = ref(false)
const rawProcessesData = ref<any>(null)

const columns = [
    { title: 'PID', key: 'pid' },
    { title: '名称', key: 'name' },
    { title: 'CPU', key: 'cpu' },
    { title: '内存', key: 'memory' }
]

async function loadProcesses() {
    try {
        const res = await axios.get('/web_api/system/processes')
        processes.value = res.data
        rawProcessesData.value = res.data
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
        <n-space vertical>
            <n-space justify="end">
                <n-button @click="showRawData = !showRawData">
                    {{ showRawData ? '显示进程列表' : '显示原始数据' }}
                </n-button>
            </n-space>

            <n-data-table v-if="!showRawData" :columns="columns" :data="processes" />

            <pre v-else
                style="width: 100%; padding: 16px; background: #f9f9f9; border-radius: 6px; overflow: auto; font-family: monospace; text-align: left; white-space: pre">{{ JSON.stringify(rawProcessesData, null, 2) }}</pre>
        </n-space>
    </n-card>
</template>