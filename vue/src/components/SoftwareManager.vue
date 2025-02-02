<script setup lang="ts">
import { ref, onMounted, h } from 'vue'
import { NCard, NDataTable, NButton, NSpace, useMessage } from 'naive-ui'
import axios from 'axios'

const message = useMessage()
const software = ref<any[]>([])
const loading = ref(false)

const columns = [
    { title: '软件名称', key: 'name' },
    { title: '状态', key: 'status' },
    {
        title: '操作',
        key: 'actions',
        render: (row: any) => h(
            NSpace,
            null,
            {
                default: () => [
                    h(
                        NButton,
                        {
                            size: 'small',
                            onClick: () => handleInstall(row.name)
                        },
                        { default: () => '安装' }
                    ),
                    h(
                        NButton,
                        {
                            size: 'small',
                            onClick: () => handleUninstall(row.name)
                        },
                        { default: () => '卸载' }
                    ),
                    h(
                        NButton,
                        {
                            size: 'small',
                            onClick: () => handleStop(row.name)
                        },
                        { default: () => '停止' }
                    )
                ]
            }
        )
    }
]

async function loadSoftwareList() {
    try {
        const res = await axios.get('http://localhost:8080/api/system/software')
        software.value = res.data.map((name: string) => ({ name }))

        // 获取每个软件的状态
        for (const item of software.value) {
            try {
                const statusRes = await axios.get(`http://localhost:8080/api/system/software/${item.name}/status`)
                item.status = statusRes.data.status
            } catch (error) {
                item.status = '未知'
            }
        }
    } catch (error) {
        message.error('获取软件列表失败')
    }
}

async function handleInstall(name: string) {
    const eventSource = new EventSource(`http://localhost:8080/api/system/software/${name}/install`)

    eventSource.onmessage = (event) => {
        message.info(event.data)
    }

    eventSource.onerror = () => {
        eventSource.close()
        loadSoftwareList()
    }
}

async function handleUninstall(name: string) {
    const eventSource = new EventSource(`http://localhost:8080/api/system/software/${name}/uninstall`)

    eventSource.onmessage = (event) => {
        message.info(event.data)
    }

    eventSource.onerror = () => {
        eventSource.close()
        loadSoftwareList()
    }
}

async function handleStop(name: string) {
    try {
        await axios.post(`http://localhost:8080/api/system/software/${name}/stop`)
        message.success('服务已停止')
        loadSoftwareList()
    } catch (error) {
        message.error('停止服务失败')
    }
}

onMounted(() => {
    loadSoftwareList()
})
</script>

<template>
    <n-card title="软件管理">
        <n-data-table :columns="columns" :data="software" :loading="loading" />
    </n-card>
</template>