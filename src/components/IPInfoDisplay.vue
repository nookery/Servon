<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { systemAPI } from '../api/info'
import type { IPInfo } from '../models/IpInfo'

const ipInfo = ref<IPInfo | null>(null)

const fetchIPInfo = async () => {
    try {
        const res = await systemAPI.getIPInfo()
        ipInfo.value = res.data
    } catch (error) {
        console.error('获取IP信息失败:', error)
    }
}

onMounted(() => {
    fetchIPInfo()
    // IP信息不需要频繁更新，设置较长的间隔
    setInterval(fetchIPInfo, 300000) // 5分钟更新一次
})
</script>

<template>
    <div class="flex items-center gap-2">
        <i class="ri-global-line text-info"></i>
        <div class="text-xs" v-if="ipInfo">
            <span>{{ ipInfo.local_ips[0].ip }}</span>
            <span class="mx-1 text-base-content/40">|</span>
            <span class="text-base-content/70">{{ ipInfo.public_ip }}</span>
        </div>
    </div>
</template>