<script setup lang="ts">
import SystemInfo from '../components/SystemInfo.vue'
</script>

<template>
    <div class="flex flex-col gap-6">
        <SystemInfo />
        <!-- 后续可以添加更多组件，如：
        <ResourceUsage />
        <ServiceStatus />
        等 -->
    </div>
</template>