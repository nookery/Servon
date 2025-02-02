<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { NCard, NDescriptions, NDescriptionsItem, useMessage } from 'naive-ui'
import axios from 'axios'

const systemInfo = ref<any>(null)
const currentUser = ref<string>('')
const message = useMessage()

onMounted(async () => {
    try {
        const [infoRes, userRes] = await Promise.all([
            axios.get('/web_api/system/basic'),
            axios.get('/web_api/system/user')
        ])
        systemInfo.value = infoRes.data
        currentUser.value = userRes.data.username
    } catch (error) {
        message.error('获取系统信息失败')
    }
})
</script>

<template>
    <n-card title="系统信息">
        <n-descriptions v-if="systemInfo" bordered>
            <n-descriptions-item label="主机名">
                {{ systemInfo.hostname }}
            </n-descriptions-item>
            <n-descriptions-item label="操作系统">
                {{ systemInfo.os }}
            </n-descriptions-item>
            <n-descriptions-item label="平台">
                {{ systemInfo.platform }}
            </n-descriptions-item>
            <n-descriptions-item label="当前用户">
                {{ currentUser }}
            </n-descriptions-item>
        </n-descriptions>
    </n-card>
</template>