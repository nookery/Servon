<script setup lang="ts">
import { ref, onMounted } from 'vue'
import Alert from '../components/Alert.vue'
import { systemAPI } from '../api/system'

interface SystemBasicInfo {
    hostname: string
    os: string
    platform: string
}

const systemInfo = ref<SystemBasicInfo | null>(null)
const currentUser = ref<string>('')
const error = ref<string>('')

onMounted(async () => {
    try {
        const [infoRes, userRes] = await Promise.all([
            systemAPI.getBasicInfo(),
            systemAPI.getCurrentUser()
        ])
        systemInfo.value = infoRes.data
        currentUser.value = userRes.data.username
        error.value = ''
    } catch (err) {
        error.value = '获取系统信息失败'
    }
})
</script>

<template>
    <div class="card bg-base-100 shadow-xl">
        <div class="card-body">
            <h2 class="card-title">系统信息</h2>

            <Alert v-if="error" type="error" :message="error" />

            <div v-if="systemInfo" class="overflow-x-auto">
                <table class="table table-zebra">
                    <tbody>
                        <tr>
                            <td class="font-bold">主机名</td>
                            <td>{{ systemInfo.hostname }}</td>
                        </tr>
                        <tr>
                            <td class="font-bold">操作系统</td>
                            <td>{{ systemInfo.os }}</td>
                        </tr>
                        <tr>
                            <td class="font-bold">平台</td>
                            <td>{{ systemInfo.platform }}</td>
                        </tr>
                        <tr>
                            <td class="font-bold">当前用户</td>
                            <td>{{ currentUser }}</td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>
    </div>
</template>