<script setup lang="ts">
import { ref } from 'vue'
import { systemAPI } from '../api/info'
import type { IPInfo } from '../types/IpInfo'
import IPInfoDisplay from '../modules/IPInfoDisplay.vue'
import IconButton from '../components/IconButton.vue'

const ipInfo = ref<IPInfo | null>(null)
const isIPInfoVisible = ref(false)

const toggleIPInfo = () => {
    isIPInfoVisible.value = !isIPInfoVisible.value
}

const fetchIPInfo = async () => {
    try {
        const res = await systemAPI.getIPInfo()
        ipInfo.value = res.data
    } catch (error) {
        console.error('获取IP信息失败:', error)
    }
}

// 初始化和定时更新IP信息
fetchIPInfo()
setInterval(fetchIPInfo, 50000)
</script>

<template>
    <div class="relative">
        <IconButton @click="toggleIPInfo" icon="ri-global-line" variant="ghost" circle title="IP信息" />
        <IPInfoDisplay :ip-info="ipInfo" :is-visible="isIPInfoVisible" />
    </div>
</template>