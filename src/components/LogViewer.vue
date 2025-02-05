<template>
    <div>
        <div class="mockup-code p-4 rounded-none font-mono text-sm h-[350px] overflow-auto">
            <pre v-for="(log, index) in formattedLogs" :key="index" v-html="log"></pre>
        </div>
        <div class="flex justify-end mt-2">
            <button class="btn btn-sm" @click="copyLogs">复制日志</button>
        </div>
    </div>
</template>

<script setup lang="ts">
import { defineProps, computed } from 'vue'
import { useToast } from '../composables/useToast'

const props = defineProps<{
    logs: string[]
}>()

const toast = useToast()

const formattedLogs = computed(() => {
    return props.logs.map(log => log.replace(/\r/g, '').replace(/\n/g, '<br>'))
})

const copyLogs = () => {
    try {
        const plainTextLogs = props.logs.join('\n')
        navigator.clipboard.writeText(plainTextLogs)
        toast.success('日志已复制到剪贴板')
    } catch (error) {
        toast.error('复制失败')
    }
}
</script>