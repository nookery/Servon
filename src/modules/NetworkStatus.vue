<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { systemAPI } from '../api/info'

const downloadSpeed = ref(0)
const uploadSpeed = ref(0)

// 获取网络资源使用情况
const fetchNetworkResources = async () => {
    try {
        const res = await systemAPI.getNetworkResources()
        downloadSpeed.value = res.data.download_speed
        uploadSpeed.value = res.data.upload_speed
    } catch (error) {
        console.error('获取网络资源信息失败:', error)
    }
}

onMounted(() => {
    fetchNetworkResources()
    setInterval(fetchNetworkResources, 50000)
})
</script>

<template>
    <div class="flex flex-col">
        <div class="flex items-center gap-1">
            <i class="ri-download-line text-xs text-base-content/70"></i>
            <span class="text-xs text-base-content/70">
                {{ downloadSpeed ? (downloadSpeed / 1024 / 1024).toFixed(1) : '0.0' }} MB/s
            </span>
        </div>
        <div class="flex items-center gap-1">
            <i class="ri-upload-line text-xs text-base-content/70"></i>
            <span class="text-xs text-base-content/70">
                {{ uploadSpeed ? (uploadSpeed / 1024 / 1024).toFixed(1) : '0.0' }} MB/s
            </span>
        </div>
    </div>
</template>
