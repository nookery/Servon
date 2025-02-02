<script setup lang="ts">
import { NLayoutHeader, NSpace, NAvatar, NIcon } from 'naive-ui'
import { ServerOutline } from '@vicons/ionicons5'
import { ref, onMounted } from 'vue'
import axios from 'axios'

const currentUser = ref('')

onMounted(async () => {
    try {
        const res = await axios.get('http://localhost:8080/api/system/user')
        currentUser.value = res.data.username
    } catch (error) {
        console.error('获取用户信息失败:', error)
    }
})
</script>

<template>
    <n-layout-header bordered style="height: 64px; padding: 0 24px">
        <div class="header-content">
            <div class="header-title">
                <n-icon size="24" :component="ServerOutline" class="header-icon" />
                服务器管理面板
            </div>
            <n-space align="center">
                <n-avatar round size="small" :style="{ background: '#2080f0' }">{{ currentUser.charAt(0).toUpperCase()
                    }}</n-avatar>
                <span>{{ currentUser }}</span>
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
</style>