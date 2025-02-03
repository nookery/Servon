<script setup lang="ts">
import { ref, onMounted } from 'vue'
import axios from 'axios'
import { useToast } from '../composables/useToast'

const software = ref<any[]>([])
const loading = ref(false)
const showLogModal = ref(false)
const currentLogs = ref<string[]>([])
const installing = ref(false)
const currentSoftware = ref<string>('')
const showRawData = ref(false)
const rawSoftwareData = ref<any>(null)

const toast = useToast()

async function handleAction(software: any) {
    currentSoftware.value = software.name
    currentLogs.value = []
    showLogModal.value = true
    installing.value = true

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
            toast.error('操作失败')
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

function closeLogModal() {
    showLogModal.value = false
    currentLogs.value = []
}

function copyLogs() {
    try {
        navigator.clipboard.writeText(currentLogs.value.join('\n'))
        toast.success('日志已复制到剪贴板')
    } catch (error) {
        toast.error('复制失败')
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

        <!-- Modal -->
        <dialog :open="showLogModal" class="modal">
            <div class="modal-box">
                <h3 class="font-bold text-lg">
                    {{ currentSoftware }} {{ installing ? '操作执行中' : '操作日志' }}
                </h3>

                <div class="py-4">
                    <div v-if="installing" class="loading loading-spinner loading-lg"></div>
                    <div class="bg-base-200 p-4 rounded-lg font-mono text-sm h-[300px] overflow-auto">
                        <div v-for="(log, index) in currentLogs" :key="index" v-html="log.replace(/\n/g, '<br>')">
                        </div>
                    </div>
                </div>

                <div class="modal-action">
                    <button class="btn" @click="copyLogs" :disabled="installing">复制日志</button>
                    <button class="btn" @click="closeLogModal" :disabled="installing">关闭</button>
                </div>
            </div>
            <form method="dialog" class="modal-backdrop">
                <button :disabled="installing">关闭</button>
            </form>
        </dialog>
    </div>
</template>

<style scoped>
.toast {
    position: fixed;
    bottom: 1rem;
    right: 1rem;
    padding: 1rem;
    border-radius: 0.5rem;
    z-index: 1000;
}
</style>