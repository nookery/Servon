<template>
    <div :class="[
        'fixed right-0 top-16 h-[calc(100vh-4rem)] bg-base-100 border-l border-base-300 transition-all duration-300 z-40',
        visible ? 'w-96' : 'w-0'
    ]">
        <div class="p-4 h-full flex flex-col" v-show="visible">
            <h3 class="font-bold mb-2 flex items-center justify-between">
                系统日志
            </h3>
            <div class="flex-1 overflow-y-auto text-sm font-mono">
                <div v-for="(log, index) in logs" :key="index" class="py-1 border-b border-base-200 last:border-0">
                    {{ log }}
                </div>
                <div v-if="logs.length === 0" class="text-base-content/50 text-center py-4">
                    暂无日志
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'

const logs = ref<string[]>([])
const eventSource = ref<EventSource | null>(null)
const visible = ref(false)

defineExpose({ visible })

onMounted(() => {
    eventSource.value = new EventSource('/web_api/logs/default')

    eventSource.value.addEventListener('log', (event) => {
        logs.value.push(event.data)
        // 保持最新的100条日志
        if (logs.value.length > 100) {
            logs.value.shift()
        }
    })
})

onUnmounted(() => {
    if (eventSource.value) {
        eventSource.value.close()
    }
})
</script>