<template>
    <div class="bg-white shadow rounded-lg p-6" :class="{
        'border-2 border-red-300': installFailed,
        'border-2 border-blue-300': installing,
        'border-2 border-green-300': installSuccess
    }">
        <div class="flex justify-between items-start">
            <div>
                <h3 class="text-lg font-semibold text-gray-900">
                    {{ software.name }}
                </h3>
                <p class="text-sm text-gray-500 mt-1">
                    {{ software.version || '未知版本' }}
                </p>
                <p v-if="software.description" class="text-gray-700 mt-2">
                    {{ software.description }}
                </p>
            </div>
            <div class="flex flex-col items-end gap-2">
                <span class="px-2 py-1 text-sm rounded-full" :class="getStatusClass(software.status)">
                    {{ getStatusText(software.status) }}
                </span>
                <button class="px-3 py-1 text-sm rounded-md text-white" :class="[
                    software.status === 'not_installed' ? 'bg-blue-500 hover:bg-blue-600' : 'bg-red-500 hover:bg-red-600',
                    installing ? 'opacity-50 cursor-not-allowed' : ''
                ]" :disabled="installing" @click="handleAction">
                    {{
                        installing ? '安装中...' :
                            software.status === 'not_installed' ? '安装' : '卸载'
                    }}
                </button>
            </div>
        </div>
        <p v-if="software.path" class="text-sm text-gray-500 mt-2">
            路径: {{ software.path }}
        </p>
        <!-- 安装状态提示 -->
        <p v-if="installFailed" class="text-sm text-red-600 mt-2">
            安装过程中断
        </p>
        <!-- 安装日志 -->
        <div v-if="logs.length > 0" class="mt-4 p-4 bg-gray-50 rounded-md">
            <div class="flex justify-between items-center mb-2">
                <h3 class="text-sm font-medium text-gray-900">安装日志</h3>
                <div class="flex gap-2">
                    <button @click="copyLogs" class="text-sm text-gray-500 hover:text-gray-700">
                        复制日志
                    </button>
                    <button @click="clearLogs" class="text-sm text-gray-500 hover:text-gray-700">
                        关闭日志
                    </button>
                </div>
            </div>
            <div class="space-y-1">
                <div v-for="(log, index) in logs" :key="index" class="text-sm text-gray-600">
                    {{ log }}
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, onUnmounted } from 'vue'

const props = defineProps({
    software: {
        type: Object,
        required: true
    }
})

const emit = defineEmits(['refresh'])

const installing = ref(false)
const installFailed = ref(false)
const installSuccess = ref(false)
const logs = ref([])

// 组件卸载时清理 EventSource
onUnmounted(() => {
    if (eventSource) {
        eventSource.close()
        eventSource = null
    }
})

function getStatusClass(status) {
    switch (status) {
        case 'running':
            return 'bg-green-100 text-green-800'
        case 'stopped':
            return 'bg-red-100 text-red-800'
        default:
            return 'bg-gray-100 text-gray-800'
    }
}

function getStatusText(status) {
    switch (status) {
        case 'running':
            return '运行中'
        case 'stopped':
            return '已停止'
        default:
            return '未安装'
    }
}

function install() {
    installing.value = true
    const eventSource = new EventSource(`/api/system/software/${props.software.name}/install`)

    eventSource.onmessage = (event) => {
        // 直接更新组件内的日志
        logs.value.push(event.data)

        if (event.data === '安装完成') {
            eventSource.close()
            installing.value = false
            installSuccess.value = true
            // 3秒后清除成功状态
            setTimeout(() => {
                installSuccess.value = false
            }, 3000)
        }
    }

    eventSource.onerror = () => {
        eventSource.close()
        installing.value = false
        // 检查最后一条日志
        const lastLog = logs.value[logs.value.length - 1]
        if (lastLog === '安装完成') {
            installSuccess.value = true
            setTimeout(() => {
                installSuccess.value = false
            }, 3000)
        } else {
            // 添加连接关闭的提示到日志
            logs.value.push('连接已关闭，安装可能未完成')
        }
    }
}

function clearLogs() {
    logs.value = []
}

async function uninstall() {
    if (!confirm(`确定要卸载 ${props.software.name} 吗？这可能会删除相关的数据和配置。`)) return

    try {
        const response = await fetch(`/api/system/software/${props.software.name}/uninstall`, {
            method: 'POST'
        })
        if (!response.ok) {
            const data = await response.json()
            throw new Error(data.error || '卸载失败')
        }
        alert('卸载成功！')
    } catch (err) {
        alert(err.message)
    }
}

function handleAction() {
    if (props.software.status === 'not_installed') {
        install()
    } else {
        uninstall()
    }
}

// 添加复制日志功能
function copyLogs() {
    const logText = logs.value.join('\n')
    navigator.clipboard.writeText(logText).then(() => {
        alert('日志已复制到剪贴板')
    }).catch(err => {
        console.error('复制失败:', err)
        alert('复制失败，请手动复制')
    })
}
</script>