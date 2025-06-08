<script setup lang="ts">
import { onMounted } from 'vue'
import LogView from '../logs/LogView.vue'

// 默认日志路径

onMounted(() => {
    // 可以在这里添加初始化逻辑
})
</script>

<template>
    <LogView current-dir="" />
</template>