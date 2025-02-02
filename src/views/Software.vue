<script setup lang="ts">
import { ref, onMounted, h } from 'vue'
import { NCard, NDataTable, NButton, NSpace, useMessage, NModal, NLog, NSpin, NCode } from 'naive-ui'
import axios from 'axios'

const message = useMessage()
const software = ref<any[]>([])
const loading = ref(false)
const showLogModal = ref(false)
const currentLogs = ref<string[]>([])
const installing = ref(false)
const currentSoftware = ref<string>('')
const showRawData = ref(false)
const rawSoftwareData = ref<any>(null)

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
                            type: row.status === 'not_installed' ? 'primary' : 'error',
                            disabled: installing.value && currentSoftware.value === row.name,
                            onClick: () => handleAction(row)
                        },
                        { default: () => row.status === 'not_installed' ? '安装' : '卸载' }
                    ),
                    row.status === 'running' && h(
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

async function handleAction(software: any) {
    currentSoftware.value = software.name
    currentLogs.value = []
    showLogModal.value = true
    installing.value = true

    const eventSource = new EventSource(
        `/web_api/system/software/${software.name}/${software.status === 'not_installed' ? 'install' : 'uninstall'}`
    )

    eventSource.onmessage = (event) => {
        // 收到第一条消息时隐藏 loading
        installing.value = false
        currentLogs.value.push(event.data)
        currentLogs.value = [...currentLogs.value]
        if (event.data.includes('完成')) {
            eventSource.close()
            loadSoftwareList()
            message.success(`${software.status === 'not_installed' ? '安装' : '卸载'}完成`)
        }
    }

    eventSource.onerror = () => {
        eventSource.close()
        installing.value = false
        if (!currentLogs.value[currentLogs.value.length - 1]?.includes('完成')) {
            currentLogs.value.push('操作异常终止')
            currentLogs.value = [...currentLogs.value]
            message.error('操作失败')
        }
        loadSoftwareList()
    }
}

async function handleStop(name: string) {
    try {
        await axios.post(`/web_api/system/software/${name}/stop`)
        message.success('服务已停止')
        loadSoftwareList()
    } catch (error) {
        message.error('停止服务失败')
    }
}

async function loadSoftwareList() {
    try {
        loading.value = true
        const res = await axios.get('/web_api/system/software')
        software.value = res.data.map((name: string) => ({ name }))
        rawSoftwareData.value = res.data

        // 获取每个软件的状态
        for (const item of software.value) {
            try {
                const statusRes = await axios.get(`/web_api/system/software/${item.name}/status`)
                item.status = statusRes.data.status
            } catch (error: any) {
                item.status = 'error'
                message.error(`获取 ${item.name} 状态失败: ${error.response?.data?.message || error.message}`)
            }
        }
    } catch (error) {
        message.error('获取软件列表失败')
    } finally {
        loading.value = false
    }
}

function closeLogModal() {
    showLogModal.value = false
    currentLogs.value = []
}

function copyLogs() {
    try {
        navigator.clipboard.writeText(currentLogs.value.join('\n'))
        message.success('日志已复制到剪贴板')
    } catch (error) {
        message.error('复制失败')
    }
}

async function handleUninstall(name: string) {
    try {
        const res = await axios.post(`/web_api/system/software/${name}/uninstall`)
        const logChan = new WebSocket(`ws://${window.location.host}/web_api/system/software/${res.data.id}/log`)

        showLogModal.value = true
        logChan.onmessage = (event) => {
            currentLogs.value.push(event.data)
        }

        logChan.onclose = async () => {
            // 卸载完成后重新加载软件状态
            await loadSoftwareList()
        }
    } catch (error) {
        message.error('卸载失败')
    }
}

onMounted(() => {
    loadSoftwareList()
})
</script>

<template>
    <n-card title="软件管理">
        <n-space vertical>
            <n-space justify="end">
                <n-button @click="showRawData = !showRawData" :disabled="loading">
                    {{ showRawData ? '显示软件列表' : '显示原始数据' }}
                </n-button>
            </n-space>

            <n-data-table v-if="!showRawData" :columns="columns" :data="software" :loading="loading" />

            <pre v-else
                style="width: 100%; padding: 16px; background: #f9f9f9; border-radius: 6px; overflow: auto; font-family: monospace; text-align: left; white-space: pre">{{ JSON.stringify(rawSoftwareData, null, 2) }}</pre>
        </n-space>

        <n-modal v-model:show="showLogModal" style="width: 600px" :mask-closable="false" preset="card"
            :title="`${currentSoftware} ${installing ? '操作执行中' : '操作日志'}`" :bordered="false">
            <n-spin :show="installing">
                <n-log :lines="currentLogs" :rows="15" :font-family="'JetBrains Mono, Menlo, Consolas, monospace'"
                    :line-height="1.25" trim style="white-space: pre" />
            </n-spin>
            <template #footer>
                <n-space>
                    <n-button @click="copyLogs" :disabled="installing">复制日志</n-button>
                    <n-button @click="closeLogModal" :disabled="installing">关闭</n-button>
                </n-space>
            </template>
        </n-modal>
    </n-card>
</template>