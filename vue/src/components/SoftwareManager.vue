<script setup lang="ts">
import { ref, onMounted, h } from 'vue'
import { NCard, NDataTable, NButton, NSpace, useMessage, NModal, NLog, NSpin } from 'naive-ui'
import axios from 'axios'

const message = useMessage()
const software = ref<any[]>([])
const loading = ref(false)
const showLogModal = ref(false)
const currentLogs = ref<string[]>([])
const installing = ref(false)
const currentSoftware = ref<string>('')

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
        `http://localhost:8080/api/system/software/${software.name}/${software.status === 'not_installed' ? 'install' : 'uninstall'}`
    )

    eventSource.onmessage = (event) => {
        currentLogs.value.push(event.data)
        if (event.data.includes('完成')) {
            eventSource.close()
            installing.value = false
            loadSoftwareList()
            message.success(`${software.status === 'not_installed' ? '安装' : '卸载'}完成`)
        }
    }

    eventSource.onerror = () => {
        eventSource.close()
        installing.value = false
        if (!currentLogs.value[currentLogs.value.length - 1]?.includes('完成')) {
            currentLogs.value.push('操作异常终止')
            message.error('操作失败')
        }
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

async function loadSoftwareList() {
    try {
        loading.value = true
        const res = await axios.get('http://localhost:8080/api/system/software')
        software.value = res.data.map((name: string) => ({ name }))

        // 获取每个软件的状态
        for (const item of software.value) {
            try {
                const statusRes = await axios.get(`http://localhost:8080/api/system/software/${item.name}/status`)
                item.status = statusRes.data.status
            } catch (error) {
                item.status = 'not_installed'
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

onMounted(() => {
    loadSoftwareList()
})
</script>

<template>
    <n-card title="软件管理">
        <n-data-table :columns="columns" :data="software" :loading="loading" />

        <n-modal v-model:show="showLogModal" style="width: 600px" :mask-closable="false" preset="card"
            :title="`${currentSoftware} ${installing ? '操作执行中' : '操作日志'}`" :bordered="false">
            <n-spin :show="installing">
                <n-log :lines="currentLogs" :rows="15" :loading="installing" />
            </n-spin>
            <template #footer>
                <n-button @click="closeLogModal" :disabled="installing">关闭</n-button>
            </template>
        </n-modal>
    </n-card>
</template>