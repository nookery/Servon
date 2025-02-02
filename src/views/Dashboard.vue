<script setup lang="ts">
import { NSpace } from 'naive-ui'
import SystemInfo from '../components/SystemInfo.vue'
</script>

<template>
    <n-space vertical :size="24">
        <SystemInfo />
        <!-- 后续可以添加更多组件，如：
        <ResourceUsage />
        <ServiceStatus />
        等 -->
    </n-space>
</template>