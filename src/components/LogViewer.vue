<template>
    <div :class="[
        'fixed right-0 top-16 h-[calc(100vh-4rem)] transition-all duration-300 z-40 terminal-container',
        visible ? 'w-96' : 'w-0'
    ]">
        <div class="h-full flex flex-col" v-show="visible">
            <div class="terminal-header px-4 py-2 flex items-center justify-between">
                <h3 class="text-terminal-text font-medium">系统日志</h3>
            </div>
            <div id="log-container" class="flex-1 overflow-y-auto terminal-body p-4 font-mono text-sm">
                <div v-for="(log, index) in logs" :key="index" class="py-0.5 terminal-line">
                    <span class="text-terminal-prompt mr-2">></span>
                    <span class="text-terminal-text">{{ log }}</span>
                </div>
                <div v-if="logs.length === 0" class="text-terminal-text/50 text-center py-4">
                    等待日志输出...
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
        // 自动滚动到底部
        setTimeout(() => {
            const container = document.getElementById('log-container')
            if (container) {
                container.scrollTop = container.scrollHeight
            }
        }, 0)
    })
})

onUnmounted(() => {
    if (eventSource.value) {
        eventSource.value.close()
    }
})
</script>

<style scoped>
.terminal-container {
    background-color: #1e1e1e;
    border-left: 1px solid #333;
}

.terminal-header {
    background-color: #2d2d2d;
    border-bottom: 1px solid #333;
}

.terminal-body {
    background-color: #1e1e1e;
}

.terminal-line {
    line-height: 1.4;
    white-space: pre-wrap;
    word-break: break-all;
}

.text-terminal-text {
    color: #e0e0e0;
}

.text-terminal-prompt {
    color: #4ec9b0;
}

/* 自定义滚动条样式 */
.terminal-body::-webkit-scrollbar {
    width: 8px;
}

.terminal-body::-webkit-scrollbar-track {
    background: #1e1e1e;
}

.terminal-body::-webkit-scrollbar-thumb {
    background: #424242;
    border-radius: 4px;
}

.terminal-body::-webkit-scrollbar-thumb:hover {
    background: #4f4f4f;
}
</style>