<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { systemAPI } from '../api/info'

const username = ref('')

onMounted(async () => {
    try {
        const response = await systemAPI.getCurrentUser()
        username.value = response.data.username
    } catch (error) {
        console.error('获取用户信息失败:', error)
    }
})
</script>

<template>
    <div>{{ username }}</div>
</template>