<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { systemAPI } from '../api/info'
import pkg from '../../package.json'

const osInfo = ref('')

const fetchOSInfo = async () => {
    try {
        const res = await systemAPI.getOSInfo()
        osInfo.value = res.data.os_info
    } catch (error) {
        console.error('获取操作系统信息失败:', error)
    }
}

onMounted(() => {
    fetchOSInfo()
    setInterval(fetchOSInfo, 50000)
})
</script>

<template>
    <div class="flex items-center gap-2">
        <i class="ri-server-line text-2xl text-primary"></i>
        <div class="flex flex-col">
            <span class="text-lg font-bold">{{ pkg.name.charAt(0).toUpperCase() + pkg.name.slice(1) }}</span>
            <span class="text-xs text-base-content/60">{{ osInfo }}</span>
        </div>
    </div>
</template>