<script setup lang="ts">
import { ref, onMounted } from 'vue'
import axios from 'axios'
import { useToast } from '../composables/useToast'
import LogViewer from '../components/LogViewer.vue'
import Modal from '../components/Modal.vue'

const software = ref<any[]>([])
const loading = ref(false)
const showLogModal = ref(false)
const currentLogs = ref<string[]>([])
const installing = ref(false)
const currentSoftware = ref<string>('')
const showRawData = ref(false)
const rawSoftwareData = ref<any>(null)
const operationFailed = ref(false)

const toast = useToast()

async function handleAction(software: any) {
    currentSoftware.value = software.name
    currentLogs.value = []
    showLogModal.value = true
    installing.value = true
    operationFailed.value = false

    const eventSource = new EventSource(
        `/web_api/system/software/${software.name}/${software.status === 'not_installed' ? 'install' : 'uninstall'}`
    )

    eventSource.onmessage = (event) => {
        installing.value = false
        currentLogs.value.push(event.data)
        currentLogs.value = [...currentLogs.value]
        if (event.data.includes('完成')) {
            eventSource.close()
            loadSoftwareList()
            toast.success(`${software.status === 'not_installed' ? '安装' : '卸载'}完成`)
        }
    }

    eventSource.onerror = () => {
        eventSource.close()
        installing.value = false
        if (!currentLogs.value[currentLogs.value.length - 1]?.includes('完成')) {
            currentLogs.value.push('操作异常终止')
            currentLogs.value = [...currentLogs.value]
            operationFailed.value = true
        }
        loadSoftwareList()
    }
}

async function handleStop(name: string) {
    try {
        await axios.post(`/web_api/system/software/${name}/stop`)
        toast.success('服务已停止')
        loadSoftwareList()
    } catch (error) {
        toast.error('停止服务失败')
    }
}

async function loadSoftwareList() {
    try {
        loading.value = true
        const res = await axios.get('/web_api/system/software')
        software.value = res.data.map((name: string) => ({ name }))
        rawSoftwareData.value = res.data

        for (const item of software.value) {
            try {
                const statusRes = await axios.get(`/web_api/system/software/${item.name}/status`)
                item.status = statusRes.data.status
            } catch (error: any) {
                item.status = 'error'
                toast.error(`获取 ${item.name} 状态失败: ${error.response?.data?.message || error.message}`)
            }
        }
    } catch (error) {
        toast.error('获取软件列表失败')
    } finally {
        loading.value = false
    }
}

onMounted(() => {
    loadSoftwareList()
})
</script>

<template>
    <div class="card bg-base-100 shadow-xl">
        <div class="card-body">
            <h2 class="card-title">软件管理</h2>

            <div class="flex justify-end mb-4">
                <button class="btn btn-primary" @click="showRawData = !showRawData" :disabled="loading">
                    {{ showRawData ? '显示软件列表' : '显示原始数据' }}
                </button>
            </div>

            <div v-if="!showRawData" class="overflow-x-auto">
                <table class="table table-zebra w-full">
                    <thead>
                        <tr>
                            <th>软件名称</th>
                            <th>状态</th>
                            <th>操作</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="item in software" :key="item.name">
                            <td>{{ item.name }}</td>
                            <td>{{ item.status }}</td>
                            <td>
                                <div class="flex gap-2">
                                    <button class="btn btn-sm"
                                        :class="item.status === 'not_installed' ? 'btn-primary' : 'btn-error'"
                                        :disabled="installing && currentSoftware === item.name"
                                        @click="handleAction(item)">
                                        {{ item.status === 'not_installed' ? '安装' : '卸载' }}
                                    </button>
                                    <button v-if="item.status === 'running'" class="btn btn-sm"
                                        @click="handleStop(item.name)">
                                        停止
                                    </button>
                                </div>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>

            <pre v-else
                class="bg-base-200 p-4 rounded-lg overflow-auto font-mono text-left whitespace-pre">{{ JSON.stringify(rawSoftwareData, null, 2) }}</pre>
        </div>

        <!-- 使用 Modal 组件 -->
        <Modal v-model:show="showLogModal"
            :title="currentSoftware + (installing ? '操作执行中' : (operationFailed ? '操作失败' : '操作日志'))"
            :loading="installing" :error="operationFailed">
            <template #default>
                <div v-if="installing" class="loading loading-spinner loading-lg"></div>
                <LogViewer :logs="currentLogs" />
            </template>
        </Modal>
    </div>
</template>