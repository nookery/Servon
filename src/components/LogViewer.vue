<template>
    <div class="mockup-code p-4 rounded-lg font-mono text-sm h-[300px] overflow-auto">
        <pre v-for="(log, index) in formattedLogs" :key="index" v-html="log"></pre>
    </div>
</template>

<script setup lang="ts">
import { defineProps, computed } from 'vue'

const props = defineProps<{
    logs: string[]
}>()

const formattedLogs = computed(() => {
    return props.logs.map(log => log.replace(/\r/g, '').replace(/\n/g, '<br>'))
})
</script>