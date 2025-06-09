<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { systemAPI } from '../api/info'

const cpuUsage = ref(0)
const memoryUsage = ref(0)
const diskUsage = ref(0)

const fetchSystemResources = async () => {
    try {
        const res = await systemAPI.getResources()
        cpuUsage.value = res.data.cpu_usage
        memoryUsage.value = res.data.memory_usage
        diskUsage.value = res.data.disk_usage
    } catch (error) {
        console.error('获取系统资源信息失败:', error)
    }
}

onMounted(() => {
    fetchSystemResources()
    setInterval(fetchSystemResources, 50000)
})
</script>

<template>
    <div class="flex gap-6">
        <!-- CPU Usage -->
        <div class="w-36">
            <div class="text-xs text-base-content/70 mb-1">
                CPU: {{ cpuUsage.toFixed(1) }}%
            </div>
            <progress class="progress progress-primary h-2" :value="cpuUsage" max="100"></progress>
        </div>

        <!-- Memory Usage -->
        <div class="w-36">
            <div class="text-xs text-base-content/70 mb-1">
                内存: {{ memoryUsage.toFixed(1) }}%
            </div>
            <progress class="progress progress-primary h-2" :value="memoryUsage" max="100"></progress>
        </div>

        <!-- Disk Usage -->
        <div class="w-36">
            <div class="text-xs text-base-content/70 mb-1">
                磁盘: {{ diskUsage.toFixed(1) }}%
            </div>
            <progress class="progress progress-primary h-2" :value="diskUsage" max="100"></progress>
        </div>
    </div>
</template>
