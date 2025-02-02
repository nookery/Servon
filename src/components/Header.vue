<script setup lang="ts">
import { NLayoutHeader, NSpace, NAvatar, NIcon, NProgress } from 'naive-ui'
import { ServerOutline } from '@vicons/ionicons5'
import { ref, onMounted } from 'vue'
import axios from 'axios'

const currentUser = ref('')
const cpuUsage = ref(0)
const memoryUsage = ref(0)
const diskUsage = ref(0)

// 获取系统资源使用情况
const fetchSystemResources = async () => {
    try {
        const res = await axios.get('http://localhost:8080/api/system/resources')
        cpuUsage.value = res.data.cpu_usage
        memoryUsage.value = res.data.memory_usage
        diskUsage.value = res.data.disk_usage
    } catch (error) {
        console.error('获取系统资源信息失败:', error)
    }
}

onMounted(async () => {
    try {
        const res = await axios.get('http://localhost:8080/api/system/user')
        currentUser.value = res.data.username
    } catch (error) {
        console.error('获取用户信息失败:', error)
    }

    // 初始获取系统资源信息
    fetchSystemResources()

    // 每5秒更新一次系统资源信息
    setInterval(fetchSystemResources, 5000)
})
</script>

<template>
    <n-layout-header bordered style="height: 64px; padding: 0 24px">
        <div class="header-content">
            <div class="header-title">
                <n-icon size="24" :component="ServerOutline" class="header-icon" />
                服务器管理面板
            </div>
            <n-space align="center" :size="24">
                <div class="resource-info">
                    <span>CPU: {{ cpuUsage.toFixed(1) }}%</span>
                    <n-progress type="line" :percentage="cpuUsage" :height="8" :border-radius="4"
                        :show-indicator="false" />
                </div>
                <div class="resource-info">
                    <span>内存: {{ memoryUsage.toFixed(1) }}%</span>
                    <n-progress type="line" :percentage="memoryUsage" :height="8" :border-radius="4"
                        :show-indicator="false" />
                </div>
                <div class="resource-info">
                    <span>磁盘: {{ diskUsage.toFixed(1) }}%</span>
                    <n-progress type="line" :percentage="diskUsage" :height="8" :border-radius="4"
                        :show-indicator="false" />
                </div>
                <n-space align="center">
                    <n-avatar round size="small" :style="{ background: '#2080f0' }">
                        {{ currentUser.charAt(0).toUpperCase() }}
                    </n-avatar>
                    <span>{{ currentUser }}</span>
                </n-space>
            </n-space>
        </div>
    </n-layout-header>
</template>

<style scoped>
.header-content {
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: space-between;
}

.header-title {
    font-size: 18px;
    font-weight: bold;
    display: flex;
    align-items: center;
    gap: 8px;
}

.header-icon {
    color: #2080f0;
}

.resource-info {
    width: 150px;
}

.resource-info span {
    font-size: 12px;
    color: #666;
}
</style>