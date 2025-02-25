<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { RiPlayCircleLine, RiStopCircleLine, RiRefreshLine } from '@remixicon/vue'
import IconButton from '../IconButton.vue'
import { getServiceList, startService, stopService, restartService } from '../../api/services'

const props = defineProps<{
    currentService: string
}>()

const emit = defineEmits<{
    (e: 'selectService', serviceName: string): void
}>()

// 服务列表
const services = ref<Array<{
    name: string
    status: string
    pid?: string
    uptime?: string
}>>([])

// 加载状态
const loading = ref(false)
const actionInProgress = ref<Record<string, boolean>>({})

// 获取服务列表
const fetchServices = async () => {
    loading.value = true
    try {
        const response = await getServiceList()
        services.value = parseServiceList(response)
    } catch (error) {
        console.error('获取服务列表失败:', error)
    } finally {
        loading.value = false
    }
}

// 解析服务列表响应
const parseServiceList = (response: string): Array<any> => {
    const lines = response.trim().split('\n')
    return lines.map(line => {
        const parts = line.trim().split(/\s+/)
        const result: any = {
            name: parts[0],
            status: parts[1]
        }

        if (parts[1] === 'RUNNING') {
            result.pid = parts[2]
            result.uptime = parts.slice(3).join(' ')
        }

        return result
    })
}

// 选择服务
const selectService = (serviceName: string) => {
    emit('selectService', serviceName)
}

// 启动服务
const handleStartService = async (serviceName: string) => {
    actionInProgress.value[serviceName] = true
    try {
        await startService(serviceName)
        await fetchServices()
    } catch (error) {
        console.error(`启动服务 ${serviceName} 失败:`, error)
    } finally {
        actionInProgress.value[serviceName] = false
    }
}

// 停止服务
const handleStopService = async (serviceName: string) => {
    actionInProgress.value[serviceName] = true
    try {
        await stopService(serviceName)
        await fetchServices()
    } catch (error) {
        console.error(`停止服务 ${serviceName} 失败:`, error)
    } finally {
        actionInProgress.value[serviceName] = false
    }
}

// 重启服务
const handleRestartService = async (serviceName: string) => {
    actionInProgress.value[serviceName] = true
    try {
        await restartService(serviceName)
        await fetchServices()
    } catch (error) {
        console.error(`重启服务 ${serviceName} 失败:`, error)
    } finally {
        actionInProgress.value[serviceName] = false
    }
}

// 获取服务状态样式
const getStatusClass = (status: string) => {
    switch (status) {
        case 'RUNNING':
            return 'badge-success'
        case 'STOPPED':
            return 'badge-error'
        case 'STARTING':
            return 'badge-warning'
        default:
            return 'badge-ghost'
    }
}

// 初始加载和定时刷新
onMounted(() => {
    fetchServices()
    // 每30秒刷新一次
    const interval = setInterval(fetchServices, 30000)

    // 组件卸载时清除定时器
    onUnmounted(() => {
        clearInterval(interval)
    })
})
</script>

<template>
    <div class="service-list-container">
        <div class="flex justify-between items-center mb-4">
            <h2 class="text-xl font-bold">系统服务</h2>
            <IconButton icon="refresh" :loading="loading" @click="fetchServices" tooltip="刷新服务列表">
                <RiRefreshLine />
            </IconButton>
        </div>

        <div v-if="loading && services.length === 0" class="loading-container">
            <span class="loading loading-spinner loading-lg"></span>
            <p>加载服务列表中...</p>
        </div>

        <div v-else-if="services.length === 0" class="empty-state">
            <p>没有找到任何服务</p>
        </div>

        <div v-else class="overflow-x-auto">
            <table class="table table-zebra w-full">
                <thead>
                    <tr>
                        <th>服务名称</th>
                        <th>状态</th>
                        <th>PID</th>
                        <th>运行时间</th>
                        <th>操作</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="service in services" :key="service.name"
                        :class="{ 'bg-base-200': service.name === props.currentService }"
                        @click="selectService(service.name)" class="cursor-pointer hover:bg-base-200">
                        <td>{{ service.name }}</td>
                        <td>
                            <span :class="['badge', getStatusClass(service.status)]">
                                {{ service.status }}
                            </span>
                        </td>
                        <td>{{ service.pid || '-' }}</td>
                        <td>{{ service.uptime || '-' }}</td>
                        <td class="flex gap-2">
                            <IconButton v-if="service.status !== 'RUNNING'" icon="play" color="success"
                                :loading="actionInProgress[service.name]" @click.stop="handleStartService(service.name)"
                                tooltip="启动服务">
                                <RiPlayCircleLine />
                            </IconButton>

                            <IconButton v-if="service.status === 'RUNNING'" icon="stop" color="error"
                                :loading="actionInProgress[service.name]" @click.stop="handleStopService(service.name)"
                                tooltip="停止服务">
                                <RiStopCircleLine />
                            </IconButton>

                            <IconButton icon="refresh" color="info" :loading="actionInProgress[service.name]"
                                @click.stop="handleRestartService(service.name)" tooltip="重启服务">
                                <RiRefreshLine />
                            </IconButton>
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>
    </div>
</template>

<style scoped>
.service-list-container {
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