<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { NCard, NDataTable, useMessage } from 'naive-ui'
import axios from 'axios'

const ports = ref<any[]>([])
const message = useMessage()

const columns = [
    { title: '端口', key: 'port' },
    { title: '协议', key: 'protocol' },
    { title: '进程', key: 'process' },
    { title: '状态', key: 'state' }
]

async function loadPorts() {
    try {
        const res = await axios.get('http://localhost:8080/api/system/ports')
        ports.value = res.data
    } catch (error) {
        message.error('获取端口列表失败')
    }
}

onMounted(() => {
    loadPorts()
})
</script>

<template>
    <n-card title="端口列表">
        <n-data-table :columns="columns" :data="ports" />
    </n-card>
</template>