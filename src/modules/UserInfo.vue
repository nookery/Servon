<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { systemAPI } from '../api/info'

const currentUser = ref('')

const fetchUserInfo = async () => {
    try {
        const res = await systemAPI.getCurrentUser()
        currentUser.value = res.data.username
    } catch (error) {
        console.error('获取用户信息失败:', error)
    }
}

onMounted(() => {
    fetchUserInfo()
})
</script>

<template>
    <div class="flex items-center gap-2">
        <div class="avatar placeholder">
            <div class="bg-primary text-primary-content rounded-full w-8">
                <span class="text-xs">{{ currentUser.charAt(0).toUpperCase() }}</span>
            </div>
        </div>
        <span class="text-sm">{{ currentUser }}</span>
    </div>
</template>