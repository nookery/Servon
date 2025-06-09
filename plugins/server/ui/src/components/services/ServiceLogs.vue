<script setup lang="ts">
import { ref, onMounted, watch, onUnmounted } from 'vue'
import { RiRefreshLine, RiDownloadLine } from '@remixicon/vue'
import IconButton from '../IconButton.vue'
import { getServiceLogs } from '../../api/services'

const props = defineProps<{
    serviceName: string
}>()

// 日志内容和控制
const logs = ref('')
const loading = ref(false)
const autoRefresh = ref(false)
const lineCount = ref(100)
const refreshInterval = ref<number | null>(null)

// 获取服务日志
const fetchLogs = async () => {
    if (!props.serviceName) return

    loading.value = true
    try {
        const response = await getServiceLogs(props.serviceName, lineCount.value)
        logs.value = response
    } catch (error) {
        console.error(`获取服务 ${props.serviceName} 日志失败:`, error)
    } finally {
        loading.value = false
    }
}

// 下载日志
const downloadLogs = () => {
    if (!logs.value) return

    const blob = new Blob([logs.value], { type: 'text/plain' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `${props.serviceName}-logs-${new Date().toISOString().slice(0, 10)}.log`
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    URL.revokeObjectURL(url)
}

// 切换自动刷新
const toggleAutoRefresh = () => {
    autoRefresh.value = !autoRefresh.value

    if (autoRefresh.value) {
        // 启动自动刷新 (每5秒)
        refreshInterval.value = window.setInterval(fetchLogs, 5000)
    } else if (refreshInterval.value !== null) {
        // 停止自动刷新
        clearInterval(refreshInterval.value)
        refreshInterval.value = null
    }
}

// 当服务名称变化时重新加载日志
watch(() => props.serviceName, (newServiceName) => {
    if (newServiceName) {
        fetchLogs()
    }
})

// 当行数变化时重新加载日志
watch(() => lineCount.value, () => {
    fetchLogs()
})

onMounted(() => {
    if (props.serviceName) {
        fetchLogs()
    }
})

onUnmounted(() => {
    // 清除自动刷新定时器
    if (refreshInterval.value !== null) {
        clearInterval(refreshInterval.value)
    }
})
</script>

<template>
    <div class="service-logs-container">
        <div class="flex justify-between items-center mb-4">
            <h2 class="text-xl font-bold">{{ serviceName }} 日志</h2>
            <div class="flex gap-2 items-center">
                <div class="form-control">
                    <select v-model="lineCount" class="select select-bordered select-sm">
                        <option :value="50">50行</option>
                        <option :value="100">100行</option>
                        <option :value="500">500行</option>
                        <option :value="1000">1000行</option>
                    </select>
                </div>

                <div class="form-control">
                    <label class="label cursor-pointer">
                        <span class="label-text mr-2">自动刷新</span>
                        <input type="checkbox" class="toggle toggle-sm toggle-primary" v-model="autoRefresh"
                            @change="toggleAutoRefresh" />
                    </label>
                </div>

                <IconButton icon="refresh" :loading="loading" @click="fetchLogs" tooltip="刷新日志">
                    <RiRefreshLine />
                </IconButton>

                <IconButton icon="download" :disabled="!logs" @click="downloadLogs" tooltip="下载日志">
                    <RiDownloadLine />
                </IconButton>
            </div>
        </div>

        <div v-if="loading && !logs" class="loading-container">
            <span class="loading loading-spinner loading-lg"></span>
            <p>加载日志中...</p>
        </div>

        <div v-else-if="!logs" class="empty-state">
            <p>没有可用的日志</p>
        </div>

        <div v-else class="logs-viewer">
            <pre
                class="bg-base-200 p-4 rounded-lg text-sm font-mono overflow-auto h-[calc(100vh-300px)]">{{ logs }}</pre>
        </div>
    </div>
</template>

<style scoped>
.service-logs-container {
    padding: 1rem;
}

.loading-container,
.empty-state {
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    height: 200px;
    color: #666;
    background-color: #f9f9f9;
    border-radius: 8px;
    margin: 20px 0;
}
</style>